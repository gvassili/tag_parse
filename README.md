# tag_parser [![GoDoc](https://godoc.org/github.com/gvassili/tag_parser?status.svg)](https://godoc.org/github.com/gvassili/tag_parser)

A golang tag key value parser

## Installation
``` sh
go get github.com/gvassili/tag_parser
```

## Example
``` go
package main

import (
  "fmt"
  "github.com/gvassili/tag_parser"
  "reflect"
)

type Test struct {
  Field string `tag:"key1,key2=value1,key3='value 1' 'value 2'"`
}

func main() {
  tag := reflect.TypeOf(Test{}).Field(0).Tag.Get("tag")
  params, err := tag_parser.Parse(tag)
  if err != nil {
    panic(fmt.Errorf("parse struct tag: %w", err))
  }
  for i, param := range params {
    fmt.Printf("param %d key=%s, values=%+v\n", i, param.Key, param.Values)
  }
}
```
output:
```
param 0 key=key1, values=[]
param 1 key=key2, values=[value1]
param 2 key=key3, values=[value 1 value 2]
```