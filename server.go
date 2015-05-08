package main

import(
  "net"
  "errors"
  "flag"
  "fmt"
)

func main() {
  port := flag.Int("p", 9000, "Port to host server on.")
  flag.Parse()

  address := fmt.Sprintf("127.0.0.1:%d", *port)
  networkAddress, _ := net.ResolveTCPAddr("tcp", address)

  listener, err := net.ListenTCP("tcp", networkAddress)
  if err != nil {
    errMsg := fmt.Sprintf("Problem connecting to port :%d on localhost\n", *port)
    errors.New(errMsg)
  }

  for {
    conn, err := listener.AcceptTCP()

    if err != nil {
      errors.New("Problem accepting connection.")
    }

    defer conn.Close()
    go handleConnection(conn)
  }
}

func handleConnection(conn *net.TCPConn) {
  for {
    msg := make([]byte, 1024)
    conn.Read(msg)
    conn.Write(msg)
  }
}
