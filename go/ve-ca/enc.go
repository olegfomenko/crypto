package ve_ca

import (
	"math/big"
)

func E(key F, val F) F {
	kBytes := FBytes(key)
	for len(kBytes) < 32 {
		kBytes = append([]byte{0x0}, kBytes...)
	}

	valBytes := FBytes(val)
	for len(valBytes) < 32 {
		valBytes = append([]byte{0x0}, valBytes...)
	}

	res := make([]byte, 0, 256)
	for i := range 32 {
		res = append(res, kBytes[i]^valBytes[i])
	}

	return new(big.Int).SetBytes(res)
}

func D(key F, cip F) F {
	kBytes := FBytes(key)
	for len(kBytes) < 32 {
		kBytes = append([]byte{0x0}, kBytes...)
	}

	cipBytes := FBytes(cip)
	for len(cipBytes) < 32 {
		cipBytes = append([]byte{0x0}, cipBytes...)
	}

	res := make([]byte, 0, 32)
	for i := range 32 {
		res = append(res, kBytes[i]^cipBytes[i])
	}

	return new(big.Int).SetBytes(res)
}
