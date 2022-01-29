package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	crypt "github.com/amoghe/go-crypt"
)

var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ=_")

// Salt generates a random salt according to given size
func Salt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("error generating salt: %s", err)
	}

	arc := uint8(0)
	for i, x := range salt {
		arc = x & 63
		salt[i] = letters[arc]
	}
	return salt, nil
}

// PasswordEncrypt encrypts given password with sha512
func PasswordEncrypt(password string) (string, error) {
	saltBytes, err := Salt(16)
	if err != nil {
		return "", err
	}
	salt := "$6$" + string(saltBytes) + "$"

	passwordEncrypted, err := crypt.Crypt(password, salt)
	if err != nil {
		return "", fmt.Errorf("error encrypting the password: %s", err)
	}
	return base64.StdEncoding.EncodeToString([]byte(passwordEncrypted)), nil
}

// TryPasswordEncrypt tries to encrypt given password if it's not encrypted
func TryPasswordEncrypt(password string) (string, error) {
	if _, err := base64.StdEncoding.DecodeString(password); err != nil {
		return PasswordEncrypt(password)
	}
	return password, nil
}
