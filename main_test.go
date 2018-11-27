package main

import (
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

func Test_getErrorJSONSource(t *testing.T) {
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
			if gotOut := getErrorJSONSource(tt.args.source, tt.args.offset); gotOut != tt.wantOut {
				t.Errorf("getErrorJSONSource() = %q, want %q", gotOut, tt.wantOut)
			}
		})
	}
}
