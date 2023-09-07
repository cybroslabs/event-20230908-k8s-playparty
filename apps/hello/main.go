package main

import (
  "fmt"
  "net/http"
  "github.com/google/uuid"
)

var instanceID string

func main() {
  // Generate a UUID for the application on startup
  instanceID = uuid.New().String()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    html := fmt.Sprintf("<html><body><h1 style=\"font-family: Sans-Serif; font-size: 4rem; margin: 0 auto; margin-top: 8rem;\">Hello from instance: %s</h1></body></html>", instanceID)
    fmt.Fprintln(w, html)
  })

  port := "8080"
  fmt.Printf("Listening on port %s...\n", port)
  err := http.ListenAndServe(":"+port, nil)
  if err != nil {
    panic(err)
  }
}