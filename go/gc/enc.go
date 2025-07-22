package gc

import (
	"crypto/rand"
	"crypto/sha256"
)

type Label [32]byte

func RL() Label {
	var buf [32]byte
	rand.Read(buf[:])
	return buf
}

func Xor(x Label, y Label) Label {
	res := Label{}

	for i := range 32 {
		res[i] = x[i] ^ y[i]
	}

	return res
}

func Hash(l Label) [32]byte {
	return sha256.Sum256(l[:])
}

func Encrypt(x Label, y Label, output Label) [32]byte {
	key := Hash(Xor(x, y))
	return Xor(key, output)
}
