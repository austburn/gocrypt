package main

import(
  "fmt"
  "golang.org/x/crypto/nacl/box"
  "crypto/rand"
)

func main() {
  var clientNonce [24]byte
  rand.Read(clientNonce[:])

  clientPublicKey, clientPrivateKey, _ := box.GenerateKey(rand.Reader)
  serverPublicKey, serverPrivateKey, _ := box.GenerateKey(rand.Reader)

  var sharedClientKey, sharedServerKey [32]byte
  box.Precompute(&sharedClientKey, serverPublicKey, clientPrivateKey)
  box.Precompute(&sharedServerKey, clientPublicKey, serverPrivateKey)

  message := []byte{'g', 'o', 'l', 'a', 'n', 'g'}

  // Seal with server's public key, aka encrypt it
  // In this scenario, we are about to send a request to the server
  // Encrypt with public, decrypt with private
  encryptedOnClient := box.SealAfterPrecomputation(nil, message, &clientNonce, &sharedClientKey)
  fmt.Printf("The client has encrypted the message: %s as %s\n", message, encryptedOnClient)

  // Open with the server's private key, aka decrypt
  decryptedOnServer, _ := box.OpenAfterPrecomputation(nil, encryptedOnClient, &clientNonce, &sharedServerKey)
  fmt.Printf("The server has decrypted the message: %s as %s\n", encryptedOnClient, decryptedOnServer)

  // Do whatever with the message we decrypted...

  // Encrypt our response, in this case, the message again
  encryptedOnServer := box.SealAfterPrecomputation(nil, decryptedOnServer, &clientNonce, &sharedServerKey)
  fmt.Printf("The server encrypted the message: %s as %s\n", decryptedOnServer, encryptedOnServer)

  // Receive on the client
  decryptedOnClient, _ := box.OpenAfterPrecomputation(nil, encryptedOnServer, &clientNonce, &sharedClientKey)
  fmt.Printf("The client decrypted the message: %s as %s\n", encryptedOnServer, decryptedOnClient)
}
