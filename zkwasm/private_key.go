package zkwasm

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

// PrivateKey represents a private key on the elliptic curve
type PrivateKey struct {
	key  *CurveField
	pubk *PublicKey
}

// NewPrivateKey creates a new private key from a given CurveField value
func NewPrivateKey(key *CurveField) *PrivateKey {
	return &PrivateKey{key: key}
}

// RandomPrivateKey generates a random private key
func RandomPrivateKey() *PrivateKey {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("random byte generation failed")
	}
	return NewPrivateKey(NewCurveField(new(big.Int).SetBytes(bytes)))
}

// PrivateKeyFromString creates a private key from a hex string
func PrivateKeyFromString(s string) *PrivateKey {
	value, _ := new(big.Int).SetString(s, 16)
	return NewPrivateKey(NewCurveField(value))
}

// ToString converts the private key to a hexadecimal string
func (pk *PrivateKey) ToString() string {
	return fmt.Sprintf("%x", pk.key.v.Bytes())
}

// R generates a random scalar value r
func (pk *PrivateKey) R() *CurveField {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("random byte generation failed")
	}
	return NewCurveField(new(big.Int).SetBytes(bytes))
}

// PublicKey returns the public key corresponding to this private key
func (pk *PrivateKey) PublicKey() *PublicKey {
	if pk.pubk == nil {
		pk.pubk = PublicKeyFromPrivateKey(pk)
	}
	return pk.pubk
}

// Sign signs a message using the private key
func (pk *PrivateKey) Sign(message []byte) [2][][]byte {
	// Ax = public key's x coordinate
	Ax := pk.PublicKey().key.x
	r := pk.R()                       // Random value r
	R := PointBase().Mul((*Field)(r)) // R = r * G
	Rx := R.x

	// Create the content for hashing: Rx || Ax || message
	var content []byte
	content = append(content, Rx.v.Bytes()...)
	content = append(content, Ax.v.Bytes()...)
	content = append(content, message...)

	// Hash the content using SHA-256
	hash := sha256.Sum256(content)
	H := new(big.Int).SetBytes(hash[:])

	// Calculate S = r + H * privateKey
	S := r.Add(pk.key.Mul(NewCurveField(H)))

	return [2][][]byte{{Rx.v.Bytes(), Rx.v.Bytes()}, {S.v.Bytes()}}
}
