package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/parser"
	"github.com/tsthght/slowparser/dao"
)

type Event struct {
	IndexNames map[string]interface{}
	SqlText    string
	Count      int
	IndexName  string
	PlanSql    string
}

var (
	input  string
	output string
	index  int
	h      bool
)

var list = make(map[string]interface{})
var state int = 1 //1:begin 2: end

func main() {
	flag.StringVar(&input, "i", "", "input name")
	flag.StringVar(&output, "o", "", "output name")
	flag.IntVar(&index, "c", 0, "only show index > 1")
	flag.BoolVar(&h, "h", false, "this help")
	flag.Parse()
	if h {
		flag.Usage()
	}

	if len(input) == 0 || len(output) == 0 {
		fmt.Printf("input or output file is empty\n")
		return
	}

	fi, err := os.Open(input)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	reg_index := regexp.MustCompile(`^# Index_names: .*$`)
	reg_plan := regexp.MustCompile(`^# Plan: .*$`)
	reg_nosql := regexp.MustCompile(`^# .*$`)

	br := bufio.NewReader(fi)

	evt := Event{make(map[string]interface{}), "", 0, "", ""}
	db, err := dao.NewDatabase(
		"root",
		"",
		"tcp",
		"127.0.0.1",
		"4000",
		"")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		if p := reg_index.FindAllString(string(a), -1); len(p) > 0 {
			s := strings.Split(p[0], "[")
			s1 := strings.Split(s[1], "]")
			evt.IndexName = s1[0]
		} else if p := reg_plan.FindAllString(string(a), -1); len(p) > 0 {
			s := strings.Split(p[0], "# Plan:")
			evt.PlanSql = "select"
			evt.PlanSql += s[1]
		} else {
			if q := reg_nosql.FindAllString(string(a), -1); len(q) == 0 {
				evt.SqlText = parser.Normalize(string(a))
				state = 2
			}
		}

		if state == 2 {
			fp := parser.DigestNormalized(evt.SqlText)
			if list[fp] == nil {
				evt.Count++
				evt.IndexNames[evt.IndexName] = 1
				if evt.IndexName != "" {
					idx := strings.Split(evt.IndexName, ":")
					filter := fmt.Sprintf("index:%s", idx[0])
					res, e := dao.QueryIndexResult(db, evt.PlanSql, filter)
					if e != nil {
						fmt.Printf("%s", e.Error())
						return
					}
					r := strings.Split(res, "index:")
					r1 := strings.Split(r[0], "range:")
					filter += r1[1]
					evt.IndexName = filter
				}
				list[fp] = evt
			} else {
				e := list[fp].(Event)
				e.Count++
				if e.IndexNames[evt.IndexName] == nil {
					e.IndexNames[evt.IndexName] = 1
				} else {
					c := e.IndexNames[evt.IndexName].(int)
					e.IndexNames[evt.IndexName] = c + 1
				}
				list[fp] = e
			}
			state = 1
			evt = Event{make(map[string]interface{}), "", 0, "", ""}
		}
	}

	fo, err := os.Create(output)
	defer fo.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for i, v := range list {
			if index > 0 && len(v.(Event).IndexNames) <= 1 {
				continue
			}
			var str string
			str += fmt.Sprintf("# FingerPrint: %v \n", i)
			str += fmt.Sprintf("# Count: %d\n", v.(Event).Count)
			str += fmt.Sprintf("# Indexs: %v\n", v.(Event).IndexNames)
			str += fmt.Sprintf("# Plan: %v\n", v.(Event).PlanSql)
			str += fmt.Sprintf("%v\n", v.(Event).SqlText)
			_, err = fo.Write([]byte(str))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

}
