# GoCrypto

This is a *basic* SSL/TLS implementation using [box](https://godoc.org/golang.org/x/crypto/nacl/box).

Basic meaning that we have a concept of a client and server that have a secure session. A handshake is performed in which they
trade public keys and begin a session. This session consists of public key encryption and sending the message along with a nonce. The receiver
decrpyts the message with the corresponding private key.

To use this, we will use the prebuilt ```play-crypto``` binary.

First, start the server:

```./play-crypto -s [-p <port>]```

Next, connect as a client:

```./play-crypto [-p <port>]```
