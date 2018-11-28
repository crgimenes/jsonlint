package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/crgimenes/goconfig"
	"github.com/gosidekick/jl"
)

func printError(a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, "\x1b[91m%v\033[0;00m\n", a...)
	if err != nil {
		fmt.Println(err)
	}
}

func printIndicator(a ...interface{}) {
	_, err := fmt.Fprintf(os.Stderr, "\x1b[96m%v\033[0;00m\n", a...)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	type configFlags struct {
		Input  string `json:"i" cfg:"i" cfgDefault:"stdin" cfgHelper:"input from"`
		Output string `json:"o" cfg:"o" cfgDefault:"stdout" cfgHelper:"output to"`
	}

	cfg := configFlags{}
	goconfig.PrefixEnv = "JSON_LINT"
	err := goconfig.Parse(&cfg)
	if err != nil {
		printError(err)
		os.Exit(-1)
	}
	var j []byte

	if cfg.Input == "stdin" {
		j, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			printError(err)
			os.Exit(-1)
		}
	} else {
		j, err = ioutil.ReadFile(cfg.Input)
		if err != nil {
			printError(err)
			os.Exit(-1)
		}
	}
	var m interface{}
	err = json.Unmarshal(j, &m)
	if err != nil {
		out, offset := jl.ParseJSONError(j, err)
		printError(out)
		if offset > 0 {
			out = jl.GetErrorJSONSource(j, offset)
			printIndicator(out)
		}
		os.Exit(-1)
	}
	j, err = json.MarshalIndent(m, "", "\t")
	if err != nil {
		printError(err)
		os.Exit(-1)
	}
	if cfg.Output == "stdout" {
		fmt.Println(string(j))
		return
	}
	err = ioutil.WriteFile(cfg.Output, j, 0644)
	if err != nil {
		printError(err)
	}
}
