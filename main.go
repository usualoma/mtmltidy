package main

import (
    "flag"
    "os"
    "io"
    "bufio"
    "fmt"
	"regexp"
	"encoding/json"
    "strings"
)

type Tags struct {
	Functions []string `json:"functions"`
	Blocks    []string `json:"blocks"`
}

var version = "v0.0.1"

func replace(f *os.File) {
    var err error

    data, _ := Asset("data/tags.json")

    tags := Tags{}
	err = json.Unmarshal(data, &tags)
	if err != nil {
		return
	}

    attrs := `(?:|\s(?:<[^>]+?>|"(?:<[^>]+?>|.)*?"|'(?:<[^>]+?>|.)*?'|.)*?)`
    attrs_re, _ := regexp.Compile(`"(?:<[^>]+?>|.)*?"|'(?:<[^>]+?>|.)*?'`)

    exp := `(?i)<\$?(?:MT:?)(?:` + strings.Join(tags.Functions, "|") + `)` + attrs + `>`
	func_start_re, _ := regexp.Compile(`(?i)^<\$?(?:MT:?)`)
	func_start_re2, _ := regexp.Compile(`(?i)^<\$?(?:MT:?)[a-zA-Z:]+`)
	func_start_re3, _ := regexp.Compile(`(?i)^</(?:MT:?)[a-zA-Z:]+`)
	func_end_re, _ := regexp.Compile(`\s*[\$/]?>$`)
	func_re, _ := regexp.Compile(exp)

    exp = `(?i)</?(?:MT:?)(?:` + strings.Join(tags.Blocks, "|") + `)` + attrs + `>`
	block_close_start_re, _ := regexp.Compile(`(?i)^</(?:MT:?)`)
	block_re, _ := regexp.Compile(exp)

    tag_map := make(map[string]string)
    for _,t := range tags.Functions {
        tag_map[strings.ToLower(t)] = t
    }
    for _,t := range tags.Blocks {
        tag_map[strings.ToLower(t)] = t
    }


    var repl func(string) string
    repl = func(str string) string {

        str = func_re.ReplaceAllStringFunc(str, func(tag string) string {
            tag = func_start_re.ReplaceAllStringFunc(tag, func(_ string) string {
                return "<$mt:"
            })
            tag = func_start_re2.ReplaceAllStringFunc(tag, func(t string) string {
                return "<$mt:" + tag_map[strings.ToLower(t[5:])]
            })
            tag = func_end_re.ReplaceAllStringFunc(tag, func(_ string) string {
                return "$>"
            })
            tag = attrs_re.ReplaceAllStringFunc(tag, func(attr string) string {
                return repl(attr)
            })
            return tag
	    })

        str = block_re.ReplaceAllStringFunc(str, func(tag string) string {
            tag = func_start_re.ReplaceAllStringFunc(tag, func(_ string) string {
                return "<mt:"
            })
            tag = func_start_re2.ReplaceAllStringFunc(tag, func(t string) string {
                return "<mt:" + tag_map[strings.ToLower(t[4:])]
            })
            tag = block_close_start_re.ReplaceAllStringFunc(tag, func(_ string) string {
                return "</mt:"
            })
            tag = func_start_re3.ReplaceAllStringFunc(tag, func(t string) string {
                return "</mt:" + tag_map[strings.ToLower(t[5:])]
            })
            tag = attrs_re.ReplaceAllStringFunc(tag, func(attr string) string {
                return repl(attr)
            })
            return tag
	    })

        return str
    }

    reader := bufio.NewReaderSize(f, 4096)
    var line string
    for line = ""; err == nil; line, err = reader.ReadString('\n') {
	    fmt.Print(repl(line))
    }
    fmt.Print(repl(line))
    if err != io.EOF {
        panic(err)
    }
}

func main() {
    flag.Parse()
    if flag.NArg() == 0 {
        replace(os.Stdin)
    }
    for i := 0; i < flag.NArg(); i++ {
        f, err := os.Open(flag.Arg(i))
        if f == nil {
            fmt.Fprintf(os.Stderr, "can't open %s: error %s\n", flag.Arg(i), err)
            os.Exit(1)
        }
        replace(f)
        f.Close()
    }
}
