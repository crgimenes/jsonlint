package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	j, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	var m interface{}
	err = json.Unmarshal(j, &m)
	var offset int64
	if err != nil {
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			log.Printf("UnmarshalTypeError: Value[%s] Type[%v]\n", e.Value, e.Type)
			fmt.Println("1>", e.Offset, e.Struct)
			offset = e.Offset
		case *json.InvalidUnmarshalError:
			log.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)
		case *json.SyntaxError:
			fmt.Println("2>", e.Offset, e.Error())
			offset = e.Offset
		default:
			log.Printf("3> %T %v", err, err)
		}

		printErrorSource([]byte(j), offset)
		lin, col := getErrorLineCol([]byte(j), offset)
		fmt.Println("lin:", lin, "col:", col)
		return
	}
}

func getErrorLineCol(source []byte, offset int64) (lin, col int) {
	for i := int64(0); i < offset; i++ {
		v := source[i]
		if v == '\r' {
			continue
		}
		if v == '\n' {
			col = 0
			lin++
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
	fmt.Printf("%s\n", source[start:end])
	fmt.Printf("%v↑\n", space)
}
