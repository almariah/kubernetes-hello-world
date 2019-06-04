package main

import (
  "github.com/corywalker/expreduce/expreduce"
  "time"
  "log"
  "math"
  //"runtime"
)

func LoadCPU() {

  // workers represrnts cores
  workers := int(math.Ceil(*cores)) - 1

  // loadPercentage will estimate CPU percentage if cores is non-integer
  loadPercentage := (*cores)/(math.Ceil(*cores))

  // initial load
  load := 1.0

  // bias will be used to fix delay that happens due to any processing time
  bias := 0.0;

  for {
    select {
    case cpu := <-currentCPU:
        log.Printf("utilize %.2f%% of CPU", cpu)
        load = cpu * loadPercentage
    default:
      t1 := time.Now()
      sum := 0
      // runtime.NumCPU()-1
      done := make(chan int)
      for i := 0; i < workers; i++ {
        go func() {
          for {
            select {
            case <-done:
                return
            default:
            }
          }
        }()
      }

      // utilize CPU usage for small slice
      for i := 0; i < 16384; i++ {
        sum += i
        time.Sleep(20 * time.Nanosecond)
      }

      // stop all utilization workers
      close(done)

      t2 := time.Now()
      d := t2.Sub(t1)

      // delay CPU usage to emulate percentage usage
      delay := float64(d.Nanoseconds()) * ( (100.0 / load) - 1.0)
      // subtract bias from delay
      delay -= bias
      time.Sleep(time.Duration(delay) * time.Nanosecond)
      t3 := time.Now()
      db := t3.Sub(t2)
      bias = float64(db.Nanoseconds()) - delay
    }
  }

}

func RefreshCurrentCPU(cpuFunc string, res int) {

  var val float64

  es := expreduce.NewEvalState()
  e := Expression{
    es: es,
  }

  i := int64(0)

  for {
    val = e.EvaluateExpression(cpuFunc, i)
    if val <= 0 {
      log.Print("current cpu expression is less than or equal 0")
      val = 1.0
    } else if val >= 100 {
      log.Print("current cpu expression is bigger than 100%")
      val = 99.0
    }
    currentCPU <- val
    time.Sleep(time.Duration(res) * time.Second)
    i++
  }

}
