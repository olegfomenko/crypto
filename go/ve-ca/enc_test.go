package ve_ca

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDecryption(t *testing.T) {
	for testCase := range 100000 {
		key := MustRandomF()
		value := MustRandomF()

		cipher := E(key, value)
		res := D(key, cipher)
		if !bytes.Equal(FBytes(res), FBytes(value)) {
			fmt.Println("Case ", testCase)
			fmt.Println(key)
			fmt.Println(value)
			fmt.Println(cipher)
			fmt.Println(res)
			panic("invalid decryption")
		}
	}
}
