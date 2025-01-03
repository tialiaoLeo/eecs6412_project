package main

import (
	"github.com/tuneinsight/lattigo/v4/bfv"
	"github.com/tuneinsight/lattigo/v4/rlwe"
)

// Encrypt function: encrypts a plain text string using AES and a private key.
func Encrypt(plainText int, msg *secure_msg) rlwe.Ciphertext {
	value1 := uint64(plainText) // First encrypted value

	plaintext1 := bfv.NewPlaintext(msg.params, msg.params.MaxLevel())
	msg.encoder.Encode([]uint64{value1}, plaintext1)

	ciphertext1 := msg.encryptor.EncryptNew(plaintext1)
	return *ciphertext1
}
