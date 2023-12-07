package el_gamal

import "testing"

func TestEncryptionRaw(t *testing.T) {
	prv, err := GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	point, err := GeneratePrivateKey()
	if err != nil {
		panic(err)
	}

	msg := point.PublicKey

	cypher, err := Encrypt(msg.X, msg.Y, prv.PublicKey)
	if err != nil {
		panic(err)
	}

	x, y := Decrypt(cypher, prv)

	if msg.X.Cmp(x) != 0 {
		panic("x result is not equal")
	}

	if msg.Y.Cmp(y) != 0 {
		panic("y result is not equal")
	}
}
