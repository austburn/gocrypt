package main

import(
  "net"
  "fmt"
  "flag"
  "os"
  "bufio"
)

func main() {
  port := flag.Int("p", 9000, "Port to connect to.")
  flag.Parse()

  address := fmt.Sprintf("127.0.0.1:%d", *port)
  fmt.Printf("Connection on %s\n", address)
  serverAddress, _ := net.ResolveTCPAddr("tcp", address)

  conn, err := net.DialTCP("tcp", nil, serverAddress)
  defer conn.Close()

  if err != nil {
    fmt.Print(err)
    os.Exit(2)
  }
  reader := bufio.NewReader(os.Stdin)
  nonce := [24]byte{'l', 'e', 't', 's', 'p', 'r', 'e', 't', 'e', 'n', 'd', 't', 'h', 'i', 's', 'i', 's', 'm', 'y', 'n', 'o', 'n', 'c', 'e'}
  for {
    fmt.Print("> ")
    msg, _ := reader.ReadBytes(0xA)

    // Kill the newline char
    msg = msg[:len(msg) - 1]
    sm := SecureMessage{msg: msg, nonce: nonce}
    _, err := conn.Write(sm.toByteArray())

    response := make([]byte, 1024)

    _, err = conn.Read(response)
    if err != nil {
      fmt.Print("Connection to the server was closed.\n")
      break
    }

    fmt.Printf("%s", response)
  }
}
