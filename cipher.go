package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func Decrypt(priv *rsa.PrivateKey, msg []byte) (string, error) {
	rng := rand.Reader
	bytes, err := rsa.DecryptPKCS1v15(rng, priv, msg)
	return string(bytes), err
}

func ReadPrivateKey(path string) (*rsa.PrivateKey, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(file)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}
