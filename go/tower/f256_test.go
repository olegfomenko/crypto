package tower

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestF256Neg(t *testing.T) {
	for _ = range 100 {
		x := RandomF256()
		assert.True(t, NewF256().Set(x).Neg().Add(x).Equal(F256Zero()))
	}
}

func TestF256Inv(t *testing.T) {
	assert.Panics(t, func() {
		NewF256().Set(F256Zero()).Inv()
	})

	for _ = range 100 {
		x := RandomF256()
		if x.Equal(F256Zero()) {
			continue
		}

		assert.True(t, NewF256().Set(x).Inv().Mul(x).Equal(F256One()))
	}
}

//func BenchmarkF256_Inv(b *testing.B) {
//	b.Run("Tower field 256 bit", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			b.StopTimer()
//			var x *F256 = RandomF256()
//			for ; ; x = RandomF256() {
//				if !x.Equal(F256Zero()) {
//					break
//				}
//			}
//
//			b.StartTimer()
//
//			x.Inv()
//		}
//	})
//
//	b.Run("BLS12-377 254 bit", func(b *testing.B) {
//		for i := 0; i < b.N; i++ {
//			b.StopTimer()
//			x := new(fr.Element)
//			x.SetRandom()
//
//			for ; ; x.SetRandom() {
//				if !x.IsZero() {
//					break
//				}
//			}
//
//			b.StartTimer()
//			x.Inverse(x)
//		}
//	})
//	//for i := 0; i < b.N; i++ {
//	//	b.StopTimer()
//	//	var x *F256 = RandomF256()
//	//	for ; ; x = RandomF256() {
//	//		if !x.Equal(F256Zero()) {
//	//			break
//	//		}
//	//	}
//	//
//	//	b.StartTimer()
//	//
//	//	x.Inv()
//	//}
//}

func BenchmarkF256_Mul(b *testing.B) {
	b.Run("Tower field 256 bit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			var x *F256 = RandomF256()
			var y *F256 = RandomF256()
			b.StartTimer()
			x.Mul(y)
		}
	})

	b.Run("BLS12-377 254 bit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			x := new(fr.Element)
			y := new(fr.Element)
			x.SetRandom()
			y.SetRandom()
			b.StartTimer()
			x.Mul(x, y)
		}
	})
}

//func BenchmarkBLSFr_Inv(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		b.StopTimer()
//		x := new(fr.Element)
//		x.SetRandom()
//
//		for ; ; x.SetRandom() {
//			if !x.IsZero() {
//				break
//			}
//		}
//
//		b.StartTimer()
//		x.Inverse(x)
//	}
//}
