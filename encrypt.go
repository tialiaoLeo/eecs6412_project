package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt function: encrypts a plain text string using AES and a private key.
func Encrypt(plainText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(createKey(key)))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))

	// Encode to Base64 for safe transport as a string.
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt function: decrypts an encrypted string using AES and a private key.
func Decrypt(cipherText, key string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(createKey(key)))
	if err != nil {
		return "", err
	}

	if len(cipherBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := cipherBytes[:aes.BlockSize]
	cipherBytes = cipherBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherBytes, cipherBytes)

	// Convert bytes back to a string.
	return string(cipherBytes), nil
}

// createKey: ensures the key is exactly 32 bytes long by padding or trimming.
func createKey(key string) string {
	const keyLength = 32 // AES-256 key length
	if len(key) > keyLength {
		return key[:keyLength]
	}
	for len(key) < keyLength {
		key += "0" // Pad with zeros
	}
	return key
}

/*
	func main() {
		privateKey := "mysecretprivatekey123"
		plainText := "abc"

		// Encrypt the string
		encrypted, err := Encrypt(plainText, privateKey)
		if err != nil {
			fmt.Println("Error encrypting:", err)
			return
		}
		fmt.Println("Encrypted:", encrypted)

		// Decrypt the string
		decrypted, err := Decrypt(encrypted, privateKey)
		if err != nil {
			fmt.Println("Error decrypting:", err)
			return
		}
		fmt.Println("Decrypted:", decrypted)
	}
*/
