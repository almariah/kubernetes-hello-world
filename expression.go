package main

import (
  "github.com/corywalker/expreduce/expreduce"
  "github.com/corywalker/expreduce/expreduce/parser"
	"strconv"
	"bytes"
	"github.com/golang/glog"
  "fmt"
  "log"
  "github.com/corywalker/expreduce/pkg/expreduceapi"
)

type Expression struct {
  es *expreduce.EvalState
}

func (e Expression) EvaluateExpression(v string, val int64) float64 {

  f := strconv.FormatInt(val, 10)
  buf := bytes.NewBufferString("n=" + fmt.Sprintf("%s", f) + ";" + v + "//N")
  res, err := parser.InterpBuf(buf, "nofile", e.es)
  res = e.es.Eval(res)
  if err != nil {
    glog.Fatal(err)
  }

  i, err := strconv.ParseFloat(res.StringForm(expreduceapi.ToStringParams{}), 64)
  if err != nil {
    log.Printf("Could't evaluate the expression: %v", err)
    return 0
  }
  return i
}
