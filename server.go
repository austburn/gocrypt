package main

import(
  "net"
  "flag"
  "fmt"
  "os"
)

func main() {
  port := flag.Int("p", 9000, "Port to host server on.")
  flag.Parse()

  address := fmt.Sprintf("127.0.0.1:%d", *port)
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
  for {
    msg := make([]byte, 1024)
    _, err := conn.Read(msg)

    if err != nil {
      break
    }

    sm := ConstructSecureMessage(msg)
    fmt.Printf("%s %s\n", sm.nonce, sm.msg)

    conn.Write(msg)
  }
}
