// Package pedersen
// Copyright 2023 Oleg Fomenko. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pedersen

import (
	"github.com/iden3/go-iden3-crypto/babyjub"
)

func BJJPointEqual(p1, p2 *babyjub.Point) bool {
	if p1 == nil && p2 == nil {
		return true
	}

	if p1 == nil || p2 == nil {
		return false
	}

	return p1.X.Cmp(p2.X) == 0 && p1.Y.Cmp(p2.Y) == 0
}
