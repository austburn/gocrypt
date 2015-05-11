# GoCrypto

This is a *basic* SSL/TLS implementation using [box](https://godoc.org/golang.org/x/crypto/nacl/box).

Basic meaning that we have a concept of a client and server that have a secure session. A handshake is performed in which they
trade public keys and begin a session in which a nonce is generated by the sender along with a message that the client/server
encrypt with their public key and it is decrypted on the other end with their private key. 