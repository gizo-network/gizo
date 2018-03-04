package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"github.com/kpango/glg"
)

//GenKeys returns private and public keys
func GenKeys() (private_key_bytes, public_key_bytes []byte) {
	privKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		glg.Fatal(err)
	}
	privKeyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		glg.Fatal(err)
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		glg.Fatal(err)
	}
	return privKeyBytes, pubKeyBytes
}
