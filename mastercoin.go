package main

import (
  "fmt"
  "os"
  "runtime"
  "container/list"
)

const (
  // First block that includes transactions to the exodus address
  GenesisBlock         uint32 = 249498
  // Server version
  Version              string = "0.1"
)

// Server interface. Handles communications to and between peers
// and provides a database interface
type Server struct {
  // List of connected peers
  peers         *list.List
  shutdownChan  chan bool
  db            *Database
}

func NewServer() (*Server, error) {
  db, err := NewDatabase()
  if err != nil {
    return nil, err
  }

  server := &Server{
    shutdownChan:   make(chan bool),
    peers:          list.New(),
    db:             db,
  }

  return server, nil
}

// Main server loop
func (s *Server) Start() {
  for {

  }

  s.shutdownChan <- true
}

func (s *Server) WaitForShutdown() {
  <- s.shutdownChan
}

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  server, err := NewServer()
  if err != nil {
    fmt.Printf("Unable to start server due to error(s) :: %v", err)
    os.Exit(1)
  }

  server.Start()

  shutdownChan := make(chan bool)
  go func() {
    // Main loop. Wait for server shutdown
    server.WaitForShutdown()

    shutdownChan <- true
  }()

  <- shutdownChan
}
