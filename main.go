package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/crgimenes/goconfig"
)

func main() {
	type configFlags struct {
		Input  string `json:"i" cfg:"i" cfgDefault:"stdin" cfgHelper:"input from"`
		Output string `json:"o" cfg:"o" cfgDefault:"stdout" cfgHelper:"output to"`
	}

	cfg := configFlags{}
	goconfig.PrefixEnv = "JSON_LINT"
	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var j []byte

	if cfg.Input == "stdin" {
		j, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	} else {
		j, err = ioutil.ReadFile(cfg.Input)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	var m interface{}
	err = json.Unmarshal(j, &m)
	var offset int64
	if err != nil {
		out := ""
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			out = fmt.Sprintf("UnmarshalTypeError: %v, Value[%s], Type[%v], offset: %v",
				e.Error(),
				e.Value,
				e.Type,
				e.Offset)
			offset = e.Offset
		case *json.SyntaxError:
			out = fmt.Sprintf("SyntaxError: %v, offset: %v",
				e.Error(),
				e.Offset)
			offset = e.Offset
		case *json.InvalidUnmarshalError:
			fmt.Fprintf(os.Stderr, "InvalidUnmarshalError: %v, Type[%v]\n",
				e.Error(),
				e.Type)
			os.Exit(-1)
		default:
			fmt.Fprintf(os.Stderr, "error: %v\n", e.Error())
			os.Exit(-1)
		}

		row, col := getErrorRowCol(j, offset)
		fmt.Fprintf(os.Stderr, "%v, row: %v, col: %v\n", out, row, col)
		printErrorSource(j, offset)
		os.Exit(-1)
	}
	j, err = json.MarshalIndent(m, "", "\t")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if cfg.Output == "stdout" {
		fmt.Println(string(j))
	} else {
		err = ioutil.WriteFile(cfg.Output, j, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getErrorRowCol(source []byte, offset int64) (row, col int) {
	for i := int64(0); i < offset; i++ {
		v := source[i]
		if v == '\r' {
			continue
		}
		if v == '\n' {
			col = 0
			row++
			continue
		}
		col++
	}
	return
}

func printErrorSource(source []byte, offset int64) {
	start := offset - 1
	limit := 0
	for ; start > 0; start-- {
		if source[start] == '\r' ||
			source[start] == '\n' ||
			limit > 38 {
			break
		}
		limit++
	}
	start++
	end := offset
	limit = 0
	for ; int64(len(source)) > end; end++ {
		if source[end] == '\r' ||
			source[end] == '\n' ||
			limit > 38 {
			break
		}
		limit++
	}
	space := ""
	for i := start; i < offset-1; i++ {
		if source[i] == '\t' {
			space += "\t"
			continue
		}
		space += " "
	}
	fmt.Fprintf(os.Stderr, "%s\n", source[start:end])
	fmt.Fprintf(os.Stderr, "%v↑\n", space)
}
