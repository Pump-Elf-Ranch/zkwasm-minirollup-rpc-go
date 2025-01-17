package zkwasm

type PublicKey struct {
	key *Point
}

func NewPublicKey(key *Point) *PublicKey {
	return &PublicKey{key: key}
}

func PublicKeyFromPrivateKey(pk *PrivateKey) *PublicKey {
	return NewPublicKey(PointBase().Mul((*Field)(pk.key)))
}
