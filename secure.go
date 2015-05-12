package main

import(
  "bytes"
  "net"
  "golang.org/x/crypto/nacl/box"
  "crypto/rand"
  "errors"
)

type SecureMessage struct {
  msg   []byte
  nonce [24]byte
}

func (s *SecureMessage) toByteArray() []byte {
  return append(s.nonce[:], s.msg[:]...)
}

func ConstructSecureMessage(sm []byte) SecureMessage {
  var nonce [24]byte
  nonceArray := sm[:24]
  copy(nonce[:], nonceArray)

  // Trim out all unnecessary bytes
  msg := bytes.Trim(sm[24:], "\x00")

  return SecureMessage{msg: msg, nonce: nonce}
}

type SecureConnection struct {
  conn *net.TCPConn
  sharedKey *[32]byte

}

func (s *SecureConnection) Read(p []byte) (int, error) {
  message := make([]byte, 2048)
  // Read the message from the buffer
  n, err := s.conn.Read(message)

  secureMessage := ConstructSecureMessage(message)
  decryptedMessage, ok := box.OpenAfterPrecomputation(nil, secureMessage.msg, &secureMessage.nonce, s.sharedKey)

  if !ok {
    return 0, errors.New("Problem decrypting the message.\n")
  }

  // Actually copy it to the destination byte array
  n = copy(p, decryptedMessage)

  return n, err
}

func (s *SecureConnection) Write(p []byte) (int, error) {
  var nonce [24]byte

  // Create a new nonce for each message sent
  rand.Read(nonce[:])

  encryptedMessage := box.SealAfterPrecomputation(nil, p, &nonce, s.sharedKey)
  sm := SecureMessage{msg: encryptedMessage, nonce: nonce}

  // Write it to the connection
  return s.conn.Write(sm.toByteArray())
}

func Handshake(conn *net.TCPConn) *[32] byte {
  var peerKey, sharedKey [32]byte

  publicKey, privateKey, _ := box.GenerateKey(rand.Reader)

  conn.Write(publicKey[:])

  peerKeyArray := make([]byte, 32)
  conn.Read(peerKeyArray)
  copy(peerKey[:], peerKeyArray)

  box.Precompute(&sharedKey, &peerKey, privateKey)

  return &sharedKey
}

