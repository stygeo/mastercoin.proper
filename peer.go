package main

import (
  "fmt"
  "github.com/conformal/btcwire"
  "net"
  "errors"
)

type OutMsg struct {
  msg       btcwire.Message
  //doneChan  chan bool
}

type Peer struct {
  server          *Server
  conn            net.Conn
  outputQueue     chan OutMsg
  version         uint32
  btcnet          btcwire.BitcoinNet
}

func (peer *Peer) InHandler() {
  for {
    msg, _, err := btcwire.ReadMessage(peer.conn, peer.version, peer.btcnet)
    if err != nil {
      fmt.Println("Error %v", err)
      continue
    }

    switch msg := msg.(type) {
    case *btcwire.MsgVerAck:
      // Do nothing
      fmt.Println("Version ackknowledged")
    default:
      fmt.Printf("Received unhandled message of type %v\n", msg.Command())
    }
  }
}

func (peer *Peer) OutHandler() {
  for {
    select {
    case msg := <- peer.outputQueue:
      switch msg.msg.(type) {
      case *btcwire.MsgVersion:
        // Version is ok
      default:
        fmt.Printf("Unknown message %v\n", msg)
        return
      }

      fmt.Println("Sending message to peer")
      peer.WriteMessage(msg.msg)
    }
  }
}

func (peer *Peer) WriteMessage(msg btcwire.Message) {
  err := btcwire.WriteMessage(peer.conn, msg, peer.version, peer.btcnet)
  if err != nil {
    fmt.Println("Error writing message %v", err)

    return
  }
}

func (peer *Peer) Start() error {
  err := peer.SendVersionMsg()
  if err != nil {
    return err
  }

  go peer.OutHandler()
  go peer.InHandler()

  return nil
}

func (peer *Peer) Stop() {
}

func (peer *Peer) QueueMsg(msg btcwire.Message) {
  peer.outputQueue <- OutMsg{msg: msg}
}

func (peer *Peer) SendVersionMsg() error {
  me, err := NewNetAddress(peer.conn.LocalAddr(), 0)
  if err != nil {
    return err
  }

  you, err := NewNetAddress(peer.conn.RemoteAddr(), 0)
  if err != nil {
    return err
  }

  msg := btcwire.NewMsgVersion(me, you, 1, "mscd", int32(peer.server.db.LastBlock()))

  msg.AddrYou.Services = btcwire.SFNodeNetwork
  msg.Services = btcwire.SFNodeNetwork

  peer.QueueMsg(msg)

  return nil
}

func NewNetAddress(addr net.Addr, services btcwire.ServiceFlag) (*btcwire.NetAddress, error) {
  // addr will be addr net.TCPAddr when not using a proxy.
  if tcpAddr, ok := addr.(*net.TCPAddr); ok {
    ip := tcpAddr.IP
    port := uint16(tcpAddr.Port)
    na := btcwire.NewNetAddressIPPort(ip, port, services)
    return na, nil
  }

  return nil, errors.New("Couldn't create address")
}

func NewPeer(server *Server) (*Peer, error) {
  conn, err := net.Dial("tcp", "62.194.122.22:8333")
  if err != nil {
    return nil, err
  }

  p := Peer{
    server:     server,
    conn:       conn,
    version:    btcwire.ProtocolVersion,
    btcnet:     btcwire.MainNet,
    outputQueue: make(chan OutMsg, 20),
  }

  return &p, nil
}
