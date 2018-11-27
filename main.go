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
		printError("%v\n", err)
		os.Exit(-1)
	}
	var j []byte

	if cfg.Input == "stdin" {
		j, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			printError("%v\n", err)
			os.Exit(-1)
		}
	} else {
		j, err = ioutil.ReadFile(cfg.Input)
		if err != nil {
			printError("%v\n", err)
			os.Exit(-1)
		}
	}
	var m interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		out, offset := parseJSONError(j, err)
		printError("%v\n", out)
		if offset > 0 {
			out = getErrorJSONSource(j, offset)
			printError("%v\n", out)
		}
		os.Exit(-1)
	}
	j, err = json.MarshalIndent(m, "", "\t")
	if err != nil {
		printError("%v\n", err)
		os.Exit(-1)
	}
	if cfg.Output == "stdout" {
		fmt.Println(string(j))
		return
	}
	err = ioutil.WriteFile(cfg.Output, j, 0644)
	if err != nil {
		printError("%v\n", err)
	}
}

func parseJSONError(source []byte, err error) (out string, offset int64) {
	offset = -1
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		row, col := getErrorRowCol(source, e.Offset)
		out = fmt.Sprintf("UnmarshalTypeError: %v, Value[%s], Type[%v], offset: %v, row: %v, col: %v",
			e.Error(),
			e.Value,
			e.Type,
			e.Offset,
			row,
			col)
		offset = e.Offset
	case *json.SyntaxError:
		row, col := getErrorRowCol(source, e.Offset)
		out = fmt.Sprintf("SyntaxError: %v, offset: %v, row: %v, col: %v",
			e.Error(),
			e.Offset,
			row,
			col)
		offset = e.Offset
	case *json.InvalidUnmarshalError:
		out = fmt.Sprintf("InvalidUnmarshalError: %v, Type[%v]\n",
			e.Error(),
			e.Type)
	default:
		out = fmt.Sprintf("error: %v\n", e.Error())
	}
	return
}

func printError(format string, a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, format, a...)
	if err != nil {
		fmt.Println(err)
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

func getStart(source []byte, offset int64) (start int64) {
	start = offset - 1
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
	return
}

func getEnd(source []byte, offset int64) (end int64) {
	end = offset
	limit := 0
	for ; int64(len(source)) > end; end++ {
		if source[end] == '\r' ||
			source[end] == '\n' ||
			limit > 38 {
			break
		}
		limit++
	}
	return
}

func getSpaces(source []byte, start, offset int64) (spaces string) {
	for i := start; i < offset-1; i++ {
		if source[i] == '\t' {
			spaces += "\t"
			continue
		}
		spaces += " "
	}
	return
}

func getErrorJSONSource(source []byte, offset int64) (out string) {
	start := getStart(source, offset)
	end := getEnd(source, offset)
	spaces := getSpaces(source, start, offset)
	out = fmt.Sprintf("%s\n%v↑", source[start:end], spaces)
	return
}
