package enc

import (
	"core-ledger/keys"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func LoadPublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode(keys.RsaPublicKey())
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func LoadPKCS1PrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func LoadPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

func LoadGTelPublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode(keys.GTelProdRsaPublicKey())
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func LoadGTelPrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keys.GTelProdRsaPrivateKey())
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), nil
}

func EncryptRSA(plainText string) (string, error) {
	pub, err := LoadPublicKey()
	if err != nil {
		return "", err
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptRSA(prvKey *rsa.PrivateKey, cipherText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, data)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
