package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
)

type encryptedPayload struct {
	Version int `json:"version"`
	Nonce string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func Encrypt(plain_text []byte, key []byte) ([]byte, error) {
	block,err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm,err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _,err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plain_text, nil)
	
	payload := encryptedPayload{
		Version: 1,
		Nonce: base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}
	return json.Marshal(payload)
}

func Decrypt (data []byte, key []byte) ([]byte, error) {
	var payload encryptedPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	nonce, err := base64.StdEncoding.DecodeString(payload.Nonce)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(payload.Ciphertext)
	if err != nil {
		return nil, err
	}

	block,err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm,err := cipher.NewGCM(block)
	if err != nil{
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}