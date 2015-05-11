package main

import(
  "net"
  "fmt"
  "os"
  "crypto/rand"
  "golang.org/x/crypto/nacl/box"
)

type Server struct {
  port int
}

func (s *Server) Run() {
  address := fmt.Sprintf("127.0.0.1:%d", s.port)
  networkAddress, _ := net.ResolveTCPAddr("tcp", address)

  listener, err := net.ListenTCP("tcp", networkAddress)
  if err != nil {
    fmt.Print(err)
    os.Exit(2)
  }

  for {
    conn, err := listener.AcceptTCP()

    if err != nil {
      fmt.Print(err)
    }

    defer conn.Close()

    go handleConnection(conn)
  }
}

func handleConnection(conn *net.TCPConn) {
  sharedKey := serverHandshake(conn)
  secureConnection := SecureConnection{conn: conn, sharedKey: sharedKey}

  for {
    msg := make([]byte, 1024)
    _, err := secureConnection.Read(msg)

    if err != nil {
      break
    }

    secureConnection.Write(msg)
  }
}

func serverHandshake(conn *net.TCPConn) *[32] byte {
  var peerKey, sharedKey [32]byte

  publicKey, privateKey, _ := box.GenerateKey(rand.Reader)

  peerKeyArray := make([]byte, 32)
  conn.Read(peerKeyArray)
  copy(peerKey[:], peerKeyArray)

  conn.Write(publicKey[:])

  box.Precompute(&sharedKey, &peerKey, privateKey)

  return &sharedKey
}

