// Copyright 2025 Alexey Yanchenko <mail@yanchenko.me>
// SPDX-License-Identifier: Apache-2.0

package gufodao

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// GetAesKey loads or generates AES key for config encryption.
// Priority: ENV -> /etc/gufo/secret.key -> generate new.
func GetAesKey() []byte {
	if key := os.Getenv("GUFO_AES_KEY"); key != "" {
		return []byte(key)
	}

	const keyPath = "/etc/gufo/secret.key"

	if data, err := os.ReadFile(keyPath); err == nil {
		return bytes.TrimSpace(data)
	}

	newKey := make([]byte, 16)
	_, _ = rand.Read(newKey)
	_ = os.WriteFile(keyPath, newKey, 0600)
	return newKey
}

// EncryptAES encrypts plain text with AES-GCM.
func EncryptAES(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts AES-GCM ciphertext.
func DecryptAES(key []byte, encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid ciphertext size")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
