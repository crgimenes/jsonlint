package jsonlint

import (
	"encoding/json"
	"fmt"
)

// ParseJSONError return better JSON error message
func ParseJSONError(source []byte, err error) (out string, offset int64) {
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
		out = fmt.Sprintf("InvalidUnmarshalError: %v, Type[%v]",
			e.Error(),
			e.Type)
	default:
		out = fmt.Sprintf("error: %v", e.Error())
	}
	return
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

/*
GetErrorJSONSource returns the error in the JSON string with
an arrow pointing exactly to the error, I hope.
*/
func GetErrorJSONSource(source []byte, offset int64) (out string) {
	start := getStart(source, offset)
	end := getEnd(source, offset)
	spaces := getSpaces(source, start, offset)
	out = fmt.Sprintf("%s\n%v↑", source[start:end], spaces)
	return
}
