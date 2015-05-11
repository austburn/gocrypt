package main

import(
  "net"
  "flag"
  "fmt"
  "os"
  "crypto/rand"
  "golang.org/x/crypto/nacl/box"
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
  var nonce [24]byte
  sharedKey := handshake(conn)

  for {
    msg := make([]byte, 1024)
    _, err := conn.Read(msg)

    if err != nil {
      break
    }

    sm := ConstructSecureMessage(msg)
    decryptedClientMessage, ok := box.OpenAfterPrecomputation(nil, sm.msg, &sm.nonce, sharedKey)

    if ok == false {
      fmt.Print("Problem decrypting the message.\n")
    }

    rand.Read(nonce[:])
    encryptedServerMessage := box.SealAfterPrecomputation(nil, decryptedClientMessage, &nonce, sharedKey)
    response := SecureMessage{msg: encryptedServerMessage, nonce: nonce}

    conn.Write(response.toByteArray())
  }
}

func handshake(conn *net.TCPConn) *[32] byte {
  var peerKey, sharedKey [32]byte

  publicKey, privateKey, _ := box.GenerateKey(rand.Reader)

  peerKeyArray := make([]byte, 32)
  conn.Read(peerKeyArray)
  copy(peerKey[:], peerKeyArray)

  conn.Write(publicKey[:])

  box.Precompute(&sharedKey, &peerKey, privateKey)

  return &sharedKey
}

