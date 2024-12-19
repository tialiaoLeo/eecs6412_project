package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateRSAKeys: Generate a pair of public and private keys for each node
func GenerateRSAKeys() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) //setting the length of the integer to be large enough
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// EncryptWithPublicKey: use the public key from the target vertex to encrypt
func EncryptWithPublicKey(message string, publicKey *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(message), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// DecryptWithPrivateKey: use its own private key to decrypt
func DecryptWithPrivateKey(encryptedMessage string, privateKey *rsa.PrivateKey) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedMessage)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedBytes, nil)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}

/*
func main() {
	//
	senderPrivateKey, _ := GenerateRSAKeys()
	receiverPrivateKey, _ := GenerateRSAKeys()

	//
	senderPublicKey := &senderPrivateKey.PublicKey
	receiverPublicKey := &receiverPrivateKey.PublicKey

	//
	message := "k-core: 5"
	encryptedMessage, err := EncryptWithPublicKey(message, receiverPublicKey)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Encrypted Message:", encryptedMessage)

	//
	decryptedMessage, err := DecryptWithPrivateKey(encryptedMessage, receiverPrivateKey)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Println("Decrypted Message:", decryptedMessage)
}
*/
