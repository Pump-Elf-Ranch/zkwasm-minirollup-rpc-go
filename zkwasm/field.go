package zkwasm

import (
	"math/big"
)

type Field struct {
	v       *big.Int
	modulus *big.Int
}

func NewField(v *big.Int) *Field {
	modulus := new(big.Int)
	modulus.SetString("21888242871839275222246405745257275088548364400416034343698204186575808495617", 10)
	value := new(big.Int).Mod(v, modulus)
	return &Field{v: value, modulus: modulus}
}

func (f *Field) String() string {
	return f.v.String()
}

func (f *Field) Add(other *Field) *Field {
	result := new(big.Int).Add(f.v, other.v)
	result.Mod(result, f.modulus)
	return &Field{v: result, modulus: f.modulus}
}

func (f *Field) Mul(other *Field) *Field {
	result := new(big.Int).Mul(f.v, other.v)
	result.Mod(result, f.modulus)
	return &Field{v: result, modulus: f.modulus}
}

func (f *Field) Sub(other *Field) *Field {
	result := new(big.Int).Sub(f.v, other.v)
	result.Mod(result, f.modulus)
	return &Field{v: result, modulus: f.modulus}
}

func (f *Field) Neg() *Field {
	result := new(big.Int).Neg(f.v)
	result.Mod(result, f.modulus)
	return &Field{v: result, modulus: f.modulus}
}

func (f *Field) Div(other *Field) *Field {
	inv := other.Inv()
	result := new(big.Int).Mul(f.v, inv.v)
	result.Mod(result, f.modulus)
	return &Field{v: result, modulus: f.modulus}
}

func (f *Field) Inv() *Field {
	if f.v.Cmp(big.NewInt(0)) == 0 {
		return f
	}

	newt, t := big.NewInt(1), big.NewInt(0)
	newr, r := new(big.Int).Set(f.v), new(big.Int).Set(f.modulus)

	for newr.Cmp(big.NewInt(0)) != 0 {
		quotient := new(big.Int).Div(r, newr)
		t, newt = newt, new(big.Int).Sub(t, new(big.Int).Mul(quotient, newt))
		r, newr = newr, new(big.Int).Sub(r, new(big.Int).Mul(quotient, newr))
	}

	if t.Cmp(big.NewInt(0)) < 0 {
		t.Add(t, f.modulus)
	}

	return &Field{v: t, modulus: f.modulus}
}
