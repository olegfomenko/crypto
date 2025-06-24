package tower

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF128Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF128()
		assert.True(t, NewF128().Set(x).Neg().Add(x).Equal(F128Zero()))
	}
}

func TestF128Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF128().Set(F128Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF128()
		if x.Equal(F128Zero()) {
			continue
		}

		assert.True(t, NewF128().Set(x).Inv().Mul(x).Equal(F128One()))
	}
}

//func BenchmarkF128_Inv(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		b.StopTimer()
//		var x *F128 = RandomF128()
//		for ; ; x = RandomF128() {
//			if !x.Equal(F128Zero()) {
//				break
//			}
//		}
//
//		b.StartTimer()
//
//		NewF128().Set(x).Inv()
//	}
//}
