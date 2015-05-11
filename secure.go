package main

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
  msg := sm[24:]
  return SecureMessage{msg: msg, nonce: nonce}
}
