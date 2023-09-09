// Package types
// Copyright 2032 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package types

import (
	"errors"
	"math/big"
)

var Zero = big.NewInt(0)

type InputData struct {
	R *BigInt `json:"r"`
	A *BigInt `json:"a"`
}

func (i InputData) Validate() error {
	if i.R.Cmp(Zero) <= 0 {
		return errors.New("invalid r value")
	}

	if i.A.Cmp(Zero) <= 0 {
		return errors.New("invalid r value")
	}

	return nil
}

func (i InputData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"r": i.R.String(),
		"a": i.A.String(),
	}
}
