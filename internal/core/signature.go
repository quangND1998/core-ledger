package core

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"strings"
)

func SignEd25519(privHex string, msg []byte) (sigHex string, err error) {
	raw, err := hex.DecodeString(strings.TrimSpace(privHex))
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}
	if len(raw) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid private key length: got %d bytes, want %d", len(raw), ed25519.PrivateKeySize)
	}

	sig := ed25519.Sign(raw, msg)

	return hex.EncodeToString(sig), nil
}

// VerifyEd25519Hex verifies msg using hex-encoded public key and signature.
// pubHex must be 32 bytes (64 hex chars). sigHex must be 64 bytes (128 hex chars).
func VerifyEd25519Hex(pubHex string, msg []byte, sigHex string) (bool, error) {
	pub, err := hex.DecodeString(strings.TrimSpace(pubHex))
	if err != nil {
		return false, fmt.Errorf("invalid public key hex: %w", err)
	}
	if len(pub) != ed25519.PublicKeySize {
		return false, fmt.Errorf("bad public key length: got %d, want %d", len(pub), ed25519.PublicKeySize)
	}

	sig, err := hex.DecodeString(strings.TrimSpace(sigHex))
	if err != nil {
		return false, fmt.Errorf("invalid signature hex: %w", err)
	}
	if len(sig) != ed25519.SignatureSize {
		return false, fmt.Errorf("bad signature length: got %d, want %d", len(sig), ed25519.SignatureSize)
	}

	ok := ed25519.Verify(ed25519.PublicKey(pub), msg, sig)
	return ok, nil
}

func GenerateEd25519Hex() (priv string, pub string, err error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}
	priv = hex.EncodeToString(privKey) // 64 bytes
	pub = hex.EncodeToString(pubKey)   // 32 bytes
	return
}
