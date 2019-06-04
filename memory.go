package main

// #cgo CFLAGS: -g -Wall
// #include <stdlib.h>
// #include <sys/mman.h>
import "C"

import (
  "github.com/corywalker/expreduce/expreduce"
  "time"
  "log"
  "unsafe"
)

func LoadMemory() {

  ptr := C.malloc(C.sizeof_char * 1)

  for {
    select {
    case mem := <-currentMemory:
      log.Printf("utilize %d bytes of memory", mem)
      C.free(unsafe.Pointer(ptr))
      ptr = C.malloc(C.sizeof_char * C.ulong(mem))
      x := C.mlock(unsafe.Pointer(ptr), C.sizeof_char * C.ulong(mem))
      if x != 0 {
        log.Fatal("Error: couldn't use mlock of specified memory")
      }
    }
  }

}

func RefreshCurrentMemory(memFunc string, res int) {

  var val float64

  es := expreduce.NewEvalState()
  e := Expression{
    es: es,
  }

  i := int64(0)

  for {
    val = e.EvaluateExpression(memFunc, i)
    if val < 0 {
      log.Print("current memory expression is less than 0")
      val = 0
    } else if val > float64(*maxMem) {
      log.Print("current memory expression is bigger than max-memory")
      val = float64(*maxMem)
    }
    currentMemory <- int64(val) * 1024 * 1024
    time.Sleep(time.Duration(res) * time.Second)
    i++
  }

}
