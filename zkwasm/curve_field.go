package zkwasm

import (
	"math/big"
)

type CurveField struct {
	v       *big.Int
	modulus *big.Int
}

func NewCurveField(v interface{}) *CurveField {
	modulus := new(big.Int)
	modulus.SetString("2736030358979909402780800718157159386076813972158567259200215660948447373041", 10)

	var value *big.Int
	switch v := v.(type) {
	case *Field:
		value = v.v
	case string:
		value = new(big.Int)
		value.SetString(v, 10)
	case int:
		value = big.NewInt(int64(v))
	case uint64:
		value = big.NewInt(int64(v))
	case int64:
		value = big.NewInt(v)
	case *big.Int:
		value = v
	default:
		panic("v must be an int, string, uint64,int64, or Field")
	}

	value.Mod(value, modulus)
	return &CurveField{v: value, modulus: modulus}
}

func (cf *CurveField) Add(f *CurveField) *CurveField {
	result := new(big.Int).Add(cf.v, f.v)
	result.Mod(result, cf.modulus)
	return &CurveField{v: result, modulus: cf.modulus}
}

func (cf *CurveField) Mul(f *CurveField) *CurveField {
	result := new(big.Int).Mul(cf.v, f.v)
	result.Mod(result, cf.modulus)
	return &CurveField{v: result, modulus: cf.modulus}
}

func (cf *CurveField) Sub(f *CurveField) *CurveField {
	result := new(big.Int).Sub(cf.v, f.v)
	result.Mod(result, cf.modulus)
	return &CurveField{v: result, modulus: cf.modulus}
}

func (cf *CurveField) Neg() *CurveField {
	result := new(big.Int).Neg(cf.v)
	result.Mod(result, cf.modulus)
	return &CurveField{v: result, modulus: cf.modulus}
}

func (cf *CurveField) Div(f *CurveField) *CurveField {
	inv := f.Inv()
	result := new(big.Int).Mul(cf.v, inv.v)
	result.Mod(result, cf.modulus)
	return &CurveField{v: result, modulus: cf.modulus}
}

func (cf *CurveField) Inv() *CurveField {
	if cf.v.Cmp(big.NewInt(0)) == 0 {
		panic("Cannot calculate the inverse of zero")
	}

	newt, t := big.NewInt(1), big.NewInt(0)
	newr, r := new(big.Int).Set(cf.v), new(big.Int).Set(cf.modulus)

	for newr.Cmp(big.NewInt(0)) != 0 {
		quotient := new(big.Int).Div(r, newr)
		t, newt = newt, new(big.Int).Sub(t, new(big.Int).Mul(quotient, newt))
		r, newr = newr, new(big.Int).Sub(r, new(big.Int).Mul(quotient, newr))
	}

	if t.Cmp(big.NewInt(0)) < 0 {
		t.Add(t, cf.modulus)
	}

	return &CurveField{v: t, modulus: cf.modulus}
}

func (cf *CurveField) String() string {
	return cf.v.String()
}
