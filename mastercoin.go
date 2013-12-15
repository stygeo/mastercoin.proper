package main

import (
  "fmt"
  "os"
  "os/signal"
  "runtime"
  "container/list"
  "time"
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
  // Channel for shutting down the server
  shutdownChan  chan bool
  // Database interface
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
  /*
   * TEMP
   */

  peer, _ := NewPeer(s)
  peer.Start()

  for {
    // Temp
    time.Sleep( time.Second )
  }
}

func (s *Server) Stop() {
  for e := s.peers.Front(); e != nil; e = e.Next() {
    if peer, ok := e.Value.(Peer); ok {
      peer.Stop()
    }
  }

  // Close database
  s.db.Close()

  s.shutdownChan <- true
}

func (s *Server) WaitForShutdown() {
  <- s.shutdownChan
}

func RegisterInterrupts(server *Server) {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func() {
    for sig := range c {
      fmt.Printf("Shutting down (%v) ...\n", sig)

      server.Stop()
    }
  }()
}

func main() {
  runtime.GOMAXPROCS(runtime.NumCPU())

  server, err := NewServer()
  if err != nil {
    fmt.Printf("Unable to start server due to error(s) :: %v", err)
    os.Exit(1)
  }

  // Register interrupt handlers for graceful shutdowns
  RegisterInterrupts(server)

  go server.Start()

  shutdownChan := make(chan bool)
  go func() {
    // Main loop. Wait for server shutdown
    server.WaitForShutdown()

    shutdownChan <- true
  }()

  <- shutdownChan
}
