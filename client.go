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

  if err != nil {
    fmt.Print(err)
  }
  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Print("> ")
    msg, _ := reader.ReadBytes(0xA)

    conn.Write(msg)

    response := make([]byte, 1024)
    conn.Read(response)
    fmt.Printf("%s", response)
  }
}
