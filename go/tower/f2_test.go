package tower

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	f00 = &F2{F1Zero(), F1Zero()}
	f01 = &F2{F1Zero(), F1One()}
	f10 = &F2{F1One(), F1Zero()}
	f11 = &F2{F1One(), F1One()}
)

func TestF2Add(t *testing.T) {
	assert.True(t, NewF2().Set(f00).Add(f00).Equal(f00))
	assert.True(t, NewF2().Set(f01).Add(f01).Equal(f00))
	assert.True(t, NewF2().Set(f10).Add(f01).Equal(f11))
	assert.True(t, NewF2().Set(f11).Add(f01).Equal(f10))
	assert.True(t, NewF2().Set(f11).Add(f10).Equal(f01))
	assert.True(t, NewF2().Set(f11).Add(f11).Equal(f00))
}

func TestF2Mul(t *testing.T) {
	assert.True(t, NewF2().Set(f00).Mul(f00).Equal(f00))
	assert.True(t, NewF2().Set(f01).Mul(f01).Equal(f01))
	assert.True(t, NewF2().Set(f10).Mul(f01).Equal(f10))
	assert.True(t, NewF2().Set(f11).Mul(f01).Equal(f11))
	assert.True(t, NewF2().Set(f11).Mul(f10).Equal(f01))
	assert.True(t, NewF2().Set(f11).Mul(f11).Equal(f10))
}

func TestF2Neg(t *testing.T) {
	assert.True(t, NewF2().Set(f00).Neg().Add(f00).Equal(F2Zero()))
	assert.True(t, NewF2().Set(f01).Neg().Add(f01).Equal(F2Zero()))
	assert.True(t, NewF2().Set(f10).Neg().Add(f10).Equal(F2Zero()))
	assert.True(t, NewF2().Set(f11).Neg().Add(f11).Equal(F2Zero()))
}

func TestF2Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF2().Set(F2Zero()).Inv()
	})

	assert.True(t, NewF2().Set(f01).Inv().Mul(f01).Equal(F2One()))
	assert.True(t, NewF2().Set(f10).Inv().Mul(f10).Equal(F2One()))
	assert.True(t, NewF2().Set(f11).Inv().Mul(f11).Equal(F2One()))
}
