package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

const (
	//32 bytes
	aesKey = "a5b8e760f241d59c3ae78c66020e0d24"
)

func padPlaintext(plaintext []byte) []byte {
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padText...)
}

func encryptPassword(plaintextPassword string) (string, error) {
	plaintext := []byte(plaintextPassword)
	paddedPlaintext := padPlaintext(plaintext)

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(paddedPlaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// comparePasswordsAES compares a plaintext password with an encrypted password using AES.
// It returns true if the passwords match, and false otherwise.

func comparePasswordsAES(plaintextPassword, encryptedPassword string) (bool, error) {
	decodedCipher, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return false, err
	}

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return false, err
	}

	if len(decodedCipher) < aes.BlockSize {
		return false, fmt.Errorf("ciphertext too short")
	}
	iv := decodedCipher[:aes.BlockSize]
	decodedCipher = decodedCipher[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedCipher, decodedCipher)

	// Remove padding from plaintext password
	padding := int(decodedCipher[len(decodedCipher)-1])
	plaintext := decodedCipher[:len(decodedCipher)-padding]

	// Compare plaintext password with provided password
	match := subtle.ConstantTimeCompare([]byte(plaintextPassword), plaintext) == 1

	// Log the result
	if match {
		log.Println("Passwords match!")
	} else {
		log.Println("Passwords do not match!")
	}

	return match, nil
}
