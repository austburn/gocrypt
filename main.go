package main

import(
  "fmt"
  "golang.org/x/crypto/nacl/box"
  "crypto/rand"
)

func main() {
  publicKey, privateKey, _ := box.GenerateKey(rand.Reader)
  var nonce [24]byte

  for i := 0; i < 4; i++ {
    test := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
    crypt := box.Seal(nil, test, &nonce, publicKey, privateKey)
    uncrypt, _ := box.Open(nil, crypt, &nonce, publicKey, privateKey)
    fmt.Printf("%s, %s, %s\n", test, crypt, uncrypt)
  }
}
