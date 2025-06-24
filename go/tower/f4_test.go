package tower

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF4Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF4()
		assert.True(t, NewF4().Set(x).Neg().Add(x).Equal(F4Zero()))
	}
}

func TestF4Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF4().Set(F4Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF4()
		if x.Equal(F4Zero()) {
			continue
		}

		assert.True(t, NewF4().Set(x).Inv().Mul(x).Equal(F4One()))
	}
}
