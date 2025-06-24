package tower

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF8Mul(t *testing.T) {
	for _ = range 1 {
		x := RandomF8()
		y := RandomF8()

		fmt.Println(x)
		fmt.Println(y)
		fmt.Println(x.Add(y))
	}
}

func TestF8Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF8()
		assert.True(t, NewF8().Set(x).Neg().Add(x).Equal(F8Zero()))
	}
}

func TestF8Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF8().Set(F8Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF8()
		if x.Equal(F8Zero()) {
			continue
		}

		assert.True(t, NewF8().Set(x).Inv().Mul(x).Equal(F8One()))
	}
}
