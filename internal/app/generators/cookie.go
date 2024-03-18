package generators

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strconv"
)

func GenerateCookie(value int) (string, error) {
	key, err := GenerateCrypto(aes.BlockSize)
	if err != nil {
		return "", err
	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce, err := GenerateCrypto(aesgcm.NonceSize())
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(
		aesgcm.Seal(nil, nonce, []byte(strconv.Itoa(value)), nil)), nil
}
