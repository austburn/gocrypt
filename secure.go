package main
import "bytes"

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
