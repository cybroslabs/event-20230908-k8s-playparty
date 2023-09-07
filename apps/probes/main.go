package main

import (
  "flag"
  "fmt"
  "math/rand"
  "net/http"
  "time"
)

var (
  startupOK = false

  // Flags
  startupProbeDelaySeconds   = 30
  notReadyTrasholdPercentage = 30
  randomizeReadiness         = new(bool)
)

func init() {
  flag.IntVar(&startupProbeDelaySeconds, "startup-delay", 30, "Initial startup probe delay in seconds")
  flag.IntVar(&notReadyTrasholdPercentage, "ready-trashold", 30, "% trashold (int)")
  randomizeReadiness = flag.Bool("ready-random", false, "Randomize readiness probe response. Use the flag to use the feature")
}

func main() {
  flag.Parse()
  rand.Seed(time.Now().UnixNano())

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "OK")
  })

  http.HandleFunc("/-/startup", func(w http.ResponseWriter, r *http.Request) {
    if !startupOK {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "NOT OK")
      return
    }

    fmt.Fprintln(w, "OK")
  })

  http.HandleFunc("/-/livez", func(w http.ResponseWriter, r *http.Request) {
    if !startupOK {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "NOT OK")
      return
    }

    fmt.Fprintln(w, "OK")
  })

  http.HandleFunc("/-/readyz", func(w http.ResponseWriter, r *http.Request) {
    if !startupOK {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "NOT READY")
      return
    }

    if !*randomizeReadiness {
      fmt.Fprintln(w, "OK")
      return
    }

    randomNumber := rand.Intn(100) // aka %

    // If randomNumber is smaller than the trashold, we respond with not ok
    if randomNumber < notReadyTrasholdPercentage {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintln(w, "NOT READY")
      return
    }

    fmt.Fprintln(w, "OK")
  })

  go waitForStartup()

  port := "8080"
  fmt.Printf("Listening on port %s...\n", port)
  err := http.ListenAndServe(":"+port, nil)
  if err != nil {
    panic(err)
  }
}

func waitForStartup() {
  time.Sleep(time.Duration(startupProbeDelaySeconds) * time.Second)
  startupOK = true
}
