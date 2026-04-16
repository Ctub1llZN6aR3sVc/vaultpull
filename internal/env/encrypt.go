package env

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts plaintext using AES-GCM with the provided key.
// key must be 16, 24, or 32 bytes.
func Encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded AES-GCM ciphertext using the provided key.
func Decrypt(encoded string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncryptSecrets encrypts all values in the map, returning a new map.
func EncryptSecrets(secrets map[string]string, key []byte) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := Encrypt([]byte(v), key)
		if err != nil {
			return nil, err
		}
		out[k] = enc
	}
	return out, nil
}

// DecryptSecrets decrypts all values in the map, returning a new map.
func DecryptSecrets(secrets map[string]string, key []byte) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		dec, err := Decrypt(v, key)
		if err != nil {
			return nil, err
		}
		out[k] = string(dec)
	}
	return out, nil
}
