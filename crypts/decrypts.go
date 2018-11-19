package crypts

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	privatePem = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`)
	privateKey []byte
)

func DecryptPrivateKey(contentBytes []byte) error {
	key, err := RSADecrypt(contentBytes)
	if err != nil {
		return err
	}
	privateKey = key
	return nil
}

// RSADecrypt decrypts a content by a private key
func RSADecrypt(contentBytes []byte) ([]byte, error) {
	prvKey, err := getPrvKey(privatePem)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, prvKey, contentBytes)
}

func getPrvKey(prvBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(prvBytes)
	if block == nil {
		return nil, errors.New("Fail to decode private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(contentBytes []byte) ([]byte, error) {
	block, err := aes.NewCipher(privateKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	contentBytes = PKCS5Padding(contentBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, privateKey[:blockSize])
	crypted := make([]byte, len(contentBytes))
	blockMode.CryptBlocks(crypted, contentBytes)
	return crypted, nil
}

func AesDecrypt(contentBytes []byte) ([]byte, error) {
	block, err := aes.NewCipher(privateKey)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, privateKey[:blockSize])
	origData := make([]byte, len(contentBytes))
	blockMode.CryptBlocks(origData, contentBytes)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
