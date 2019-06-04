package main

import (
  "fmt"
  "net/http"
  "log"
  "runtime"
  "time"
  "os"
  "encoding/json"
  "flag"
)

type HealthCheck struct {
  Status string `json:"status"`
}

var (

  currentMemory chan int64
  currentCPU chan float64

	color = flag.String("color", "white", "The color of the hello world web page.")
  initDuration = flag.Int("init-duration", 0, "The delay duration to simulate initialization.")

  maxMem = flag.Int("max-memory", 0, "The maximum memory to be allocated.")
  memoryExpression = flag.String("memory-expression", "0", "The mathematical expression of memory in MB as a function of n.")

  cores = flag.Float64("cores", 0, "The value of cores that is allocated to the application. It could be non-integer value for cores in mili")
  cpuExpression = flag.String("cpu-expression", "0", "The mathematical expression of CPU percentage as a function of n.")

  resolution = flag.Int("resolution", 10, "The duration in seconds at which the expressions will be evaluated.")
)

func main() {

  flag.Parse()

  if (*maxMem != 0 && *memoryExpression == "0")  || (*maxMem == 0 && *memoryExpression != "0") {
    log.Fatal("Error: you have to set both max-memory and memory-expression!")
  }

  if (*cores != 0 && *cpuExpression == "0") || (*cores == 0 && *cpuExpression != "0") {
    log.Fatal("Error: you have to set both cores and cpu-expression!")
  }

  if *resolution < 10 {
    log.Fatal("Error: resolution is less than 10 seconds!")
  }

  // to simulate init
  log.Printf("Simulate delay of %d seconds!", *initDuration)
  time.Sleep(time.Duration(*initDuration) * time.Second)

  if *memoryExpression != "0" {
    currentMemory = make(chan int64)
    go RefreshCurrentMemory(*memoryExpression, *resolution)
    go LoadMemory()
  }

  if *cpuExpression != "0" {
    currentCPU = make(chan float64)
    go RefreshCurrentCPU(*cpuExpression, *resolution)
    go LoadCPU()
  }

  router := http.NewServeMux()
	router.HandleFunc("/", hello)
	router.HandleFunc("/health", health)

  log.Print("Starting Kubernetes Hello World at port 8080")
  log.Fatal(http.ListenAndServe(":8080", router))
}

func hello(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
  fmt.Fprint(w, render(*color))
  return
}

func health(w http.ResponseWriter, r *http.Request) {
  hc := HealthCheck{
    Status: "up",
  }
  b, err := json.Marshal(hc)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(b)
}

func getHostname() (string, error) {
  name, err := os.Hostname()
  if err != nil {
    return "", err
  }
  return name, nil
}

func render(color string) string {

  hostname, _ := getHostname()

  alloc, totalAlloc, sys, numGC := getMemUsage()

  html := fmt.Sprintf(`
  <html>
    <body style="background-color:%s;">
      <h1>Kubernetes Hello World</h1>
      <p>
        Hostname = %s<br>
        Memory Alloc = %v MiB"<br>
        Memory Total Alloc = %v MiB<br>
        Memory Sys = %v MiB<br>
        NumGC = %v
      </p>
    </body>
  </html>
  `, color, hostname, alloc,totalAlloc,sys,numGC)
  return html
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

func getMemUsage() (uint64, uint64, uint64, uint32) {
  var m runtime.MemStats
  runtime.ReadMemStats(&m)
  return bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC
}
