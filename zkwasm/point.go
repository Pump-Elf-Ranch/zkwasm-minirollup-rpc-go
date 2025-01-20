package zkwasm

import (
	"fmt"
	"math/big"
)

// Point struct represents a point on the elliptic curve
type Point struct {
	x, y *Field
}

// NewPoint creates a new Point with the given x and y coordinates
func NewPoint(x, y *Field) *Point {
	return &Point{x: x, y: y}
}

// Zero returns the point at infinity (the identity element)
func (p *Point) Zero() *Point {
	return NewPoint(NewField(toBigInt("0")), NewField(toBigInt("1"))) // point at infinity
}

// IsZero checks if the point is the point at infinity
func (p *Point) IsZero() bool {
	return p.x.v.Cmp(big.NewInt(0)) == 0 && p.y.v.Cmp(big.NewInt(1)) == 0
}

// Base returns the base point of the elliptic curve
func (p *Point) Base() *Point {
	gX := Constants["gX"]
	gY := Constants["gY"]
	return NewPoint(gX, gY)
}

// Add adds two points on the elliptic curve
func (p *Point) Add(other *Point) *Point {
	// Curve equation: y^2 = x^3 + ax + b (here we use a simplified version)
	u1 := p.x
	v1 := p.y
	u2 := other.x
	v2 := other.y

	// Compute the numerator and denominator for x-coordinate
	u3_m := u1.Mul(v2).Add(v1.Mul(u2))                                                  // u3_m = u1*v2 + v1*u2
	u3_d := Constants["d"].Mul(u1).Mul(u2).Mul(v1).Mul(v2).Add(NewField(toBigInt("1"))) // u3_d = 1 + d*u1*u2*v1*v2
	u3 := u3_m.Div(u3_d)                                                                // u3 = u3_m / u3_d

	// Compute the numerator and denominator for y-coordinate
	v3_m := v1.Mul(v2).Sub(Constants["a"].Mul(u1).Mul(u2))                              // v3_m = v1*v2 - a*u1*u2
	v3_d := NewField(toBigInt("1")).Sub(Constants["d"].Mul(u1).Mul(u2).Mul(v1).Mul(v2)) // v3_d = 1 - d*u1*u2*v1*v2
	v3 := v3_m.Div(v3_d)                                                                // v3 = v3_m / v3_d

	return NewPoint(u3, v3)
}

// Mul performs scalar multiplication (k * P)
func (p *Point) Mul(k *Field) *Point {
	// Scalar multiplication: repeated doubling and adding (binary method)
	result := p.Zero() // start with point at infinity
	acc := p

	for i := 0; i < k.v.BitLen(); i++ {
		if k.v.Bit(i) == 1 {
			result = result.Add(acc)
		}
		acc = acc.Add(acc) // double the point
	}

	return result
}

// String returns a string representation of the point (x, y)
func (p *Point) String() string {
	return fmt.Sprintf("Point(x: %s, y: %s)", p.x.String(), p.y.String())
}

func PointBase() *Point {
	gX := Constants["gX"]
	gY := Constants["gY"]
	return NewPoint(gX, gY)
}
