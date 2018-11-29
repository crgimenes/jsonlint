# JSON Lint
[![Build Status](https://travis-ci.org/gosidekick/jsonlint.svg?branch=master)](https://travis-ci.org/gosidekick/jsonlint)
[![Go Report Card](https://goreportcard.com/badge/github.com/gosidekick/jsonlint)](https://goreportcard.com/report/github.com/gosidekick/jsonlint)
[![GoDoc](https://godoc.org/github.com/gosidekick/jsonlint?status.png)](https://godoc.org/github.com/gosidekick/jsonlint)
[![Go project version](https://badge.fury.io/go/github.com%2Fgosidekick%2Fjsonlint.svg)](https://badge.fury.io/go/github.com/gosidekick/jsonlint)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

A small utility to validate, format JSON and when it has an error, if possible indicate with an arrow where the error is.

It can also be used as a Go package.

## Example

```console
cat << EOF | jl                                                                       {
    "name": "John",
    "age": 30,
    "cars": {
        "car1": "Ford",
        "car2":: "BMW",
        "car3": "Fiat"
    }
}
EOF
SyntaxError: invalid character ':' looking for beginning of value, offset: 91, row: 5, col: 16
        "car2":: "BMW",
               ↑
```

## As a package

Use as a Go package to get a better error message.

```console
go get -u github.com/gosidekick/jsonlint
```

```go
j := `{
    "name": "John",
    "age": 30,
    "cars": {
        "car1": "Ford",
        "car2":: "BMW",
        "car3": "Fiat"
    }
}`
err := json.Unmarshal(j, &m)
if err != nil {
	out, offset := jsonlint.ParseJSONError(j, err)
	fmt.Println(out) // print the error message
	if offset > 0 {
		out = jsonlint.GetErrorJSONSource(j, offset)
		fmt.Println(out) // print the arrow
	}
	return
}
```

## Contributing

- Fork the repo on GitHub
- Clone the project to your own machine
- Create a *branch* with your modifications `git checkout -b fantastic-feature`.
- Then _commit_ your changes `git commit -m 'Implementation of new fantastic feature'`
- Make a _push_ to your _branch_ `git push origin fantastic-feature`.
- Submit a **Pull Request** so that we can review your changes
