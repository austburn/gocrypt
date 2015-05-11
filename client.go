package main

import(
  "net"
  "fmt"
  "flag"
  "os"
  "bufio"
  "crypto/rand"
  "golang.org/x/crypto/nacl/box"
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

  var nonce [24]byte
  sharedKey := handshake(conn)

  reader := bufio.NewReader(os.Stdin)

  for {
    fmt.Print("> ")
    msg, _ := reader.ReadBytes(0xA)

    // Kill the newline char
    msg = msg[:len(msg) - 1]

    rand.Read(nonce[:])
    encryptedMessage := box.SealAfterPrecomputation(nil, msg, &nonce, sharedKey)

    sm := SecureMessage{msg: encryptedMessage, nonce: nonce}

    _, err := conn.Write(sm.toByteArray())

    response := make([]byte, 1024)

    _, err = conn.Read(response)
    if err != nil {
      fmt.Print("Connection to the server was closed.\n")
      break
    }

    serverMessage := ConstructSecureMessage(response)
    decryptedServerMessage, ok := box.OpenAfterPrecomputation(nil, serverMessage.msg, &serverMessage.nonce, sharedKey)

    if ok == false {
      fmt.Printf("Problem decrypting message.\n")
    }

    fmt.Printf("%s\n", decryptedServerMessage)
  }
}

func handshake(conn *net.TCPConn) *[32] byte {
  var peerKey, sharedKey [32]byte

  publicKey, privateKey, _ := box.GenerateKey(rand.Reader)

  conn.Write(publicKey[:])

  peerKeyArray := make([]byte, 32)
  conn.Read(peerKeyArray)
  copy(peerKey[:], peerKeyArray)

  box.Precompute(&sharedKey, &peerKey, privateKey)

  return &sharedKey
}
