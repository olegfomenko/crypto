package tower

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF32Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF32()
		assert.True(t, NewF32().Set(x).Neg().Add(x).Equal(F32Zero()))
	}
}

func TestF32Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF32().Set(F32Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF32()
		if x.Equal(F32Zero()) {
			continue
		}

		assert.True(t, NewF32().Set(x).Inv().Mul(x).Equal(F32One()))
	}
}
