package field

import (
	"log"

	"github.com/cznic/mathutil"
)

// UInt64 - GO basic object notation
// to perform operations on basic types
type UInt64 uint64

// P - size of field
// The field size -- 2305843009213693951
const P UInt64 = (1 << 61) - 1

// Field -- It is consistent iff 0 <= src.0 <= P.
type Field struct {
	fp UInt64
}

func (src UInt64) reduceOnce() UInt64 {
	var value = (src & P) + (src >> 61)
	if value == P {
		return 0
	}
	return value
}

func (src UInt64) reduceOnceAssert() UInt64 {
	var res = src.reduceOnce()
	if res >= P {
		log.Fatal("Error: Expected result should be less than field size")
	}
	return res
}

func (src UInt64) reduceOnceMul(op2 UInt64) UInt64 {
	var value = (src << 3) | (op2 >> 61)
	value = (op2 & P) + value
	if value == P {
		return 0
	}
	return value
}

func asLimbs(x UInt64) (uint32, uint32) {
	return uint32(x >> 32), uint32(x)
}

// NewField reduces initial value in the field
func NewField(value UInt64) Field {
	return Field{value.reduceOnce().reduceOnceAssert()}
}

// Neg negates number in the field
func (src Field) Neg() Field {
	return Field{(P - src.fp).reduceOnce().reduceOnceAssert()}
}

// Add adds two field elements and reduces the resulting number in the field
func (src Field) Add(op2 Field) Field {
	return Field{(src.fp + op2.fp).reduceOnce().reduceOnceAssert()}
}

// AddAssign works same as Add, assigns final value to src
func (src *Field) AddAssign(op2 Field) {
	*src = src.Add(op2)
}

// Sub subtracts two field elements and reduces the resulting number in the field
func (src Field) Sub(op2 Field) Field {
	if op2.fp > src.fp {
		return Field{(P - op2.fp + src.fp).reduceOnce().reduceOnceAssert()}
	}
	return Field{(src.fp - op2.fp).reduceOnce().reduceOnceAssert()}
}

// SubAssign works same as Sub, assigns final value to src
func (src *Field) SubAssign(op2 Field) {
	*src = src.Sub(op2)
}

// Mul muliplies two field elements and reduces the resulting number in the field
func (src Field) Mul(op2 Field) Field {
	var high, low uint64 = mathutil.MulUint128_64(uint64(src.fp), uint64(op2.fp))
	var rh, rl = UInt64(high), UInt64(low)
	var res = rh.reduceOnceMul(rl).reduceOnceAssert()
	return Field{res}
}

// MulAssign works same as Mul, assigns final value to src
func (src *Field) MulAssign(op2 Field) {
	*src = src.Mul(op2)
}
