package jsonlint

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func Test_getErrorRowCol(t *testing.T) {
	type args struct {
		source []byte
		offset int64
	}
	tests := []struct {
		name    string
		args    args
		wantRow int
		wantCol int
	}{
		{
			name: "multiple lines",
			args: args{
				source: []byte("xxxx\nxxxx\nxxXxx\nxxxx\nxxxx\n"),
				offset: 13,
			},
			wantRow: 2,
			wantCol: 3,
		},
		{
			name: "multiple lines on Windows",
			args: args{
				source: []byte("xxxx\r\nxxxx\r\nxxXxx\r\nxxxx\r\nxxxx\n"),
				offset: 15,
			},
			wantRow: 2,
			wantCol: 3,
		},
		{
			name: "empty string",
			args: args{
				source: []byte(""),
				offset: 0,
			},
			wantRow: 0,
			wantCol: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRow, gotCol := getErrorRowCol(tt.args.source, tt.args.offset)
			if gotRow != tt.wantRow {
				t.Errorf("getErrorRowCol() gotLin = %v, want %v", gotRow, tt.wantRow)
			}
			if gotCol != tt.wantCol {
				t.Errorf("getErrorRowCol() gotCol = %v, want %v", gotCol, tt.wantCol)
			}
		})
	}
}

func Test_GetErrorJSONSource(t *testing.T) {
	type args struct {
		source []byte
		offset int64
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "simple offset",
			args: args{
				source: []byte("xxxx\nXXXXX\nxxxx"),
				offset: 8,
			},
			wantOut: "XXXXX\n  ↑",
		},
		{
			name: "with tab",
			args: args{
				source: []byte("xxxx\nX\tXXXX\nxxxx"),
				offset: 8,
			},
			wantOut: "X\tXXXX\n \t↑",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := GetErrorJSONSource(tt.args.source, tt.args.offset); gotOut != tt.wantOut {
				t.Errorf("GetErrorJSONSource() = %q, want %q", gotOut, tt.wantOut)
			}
		})
	}
}

func TestParseJSONError(t *testing.T) {
	type args struct {
		source []byte
		err    error
	}
	tests := []struct {
		name       string
		args       args
		wantOut    string
		wantOffset int64
	}{
		{
			name: "simple UnmarshalTypeError",
			args: args{
				source: []byte("xxxx\nXXXXX\nxxxx"),
				err: &json.UnmarshalTypeError{
					Value:  "array",
					Type:   reflect.TypeOf(""),
					Offset: 7,
					Struct: "T",
					Field:  "X"},
			},
			wantOut:    "UnmarshalTypeError: json: cannot unmarshal array into Go struct field T.X of type string, Value[array], Type[string], offset: 7, row: 1, col: 2",
			wantOffset: 7,
		},
		{
			name: "simple SyntaxError",
			args: args{
				source: []byte("xxxx\nXXXXX\nxxxx"),
				err: &json.SyntaxError{
					Offset: 7,
				},
			},
			wantOut:    "SyntaxError: , offset: 7, row: 1, col: 2",
			wantOffset: 7,
		},
		{
			name: "simple InvalidUnmarshalError",
			args: args{
				source: []byte(""),
				err: &json.InvalidUnmarshalError{
					Type: reflect.TypeOf(""),
				},
			},
			wantOut:    "InvalidUnmarshalError: json: Unmarshal(non-pointer string), Type[string]",
			wantOffset: -1,
		},
		{
			name: "default error",
			args: args{
				source: []byte(""),
				err:    errors.New("test"),
			},
			wantOut:    "error: test",
			wantOffset: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, gotOffset := ParseJSONError(tt.args.source, tt.args.err)
			if gotOut != tt.wantOut {
				t.Errorf("ParseJSONError() gotOut = %q, want %q", gotOut, tt.wantOut)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("ParseJSONError() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}
