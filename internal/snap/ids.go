package snap

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// NewNumericID returns a random numeric string of exactly n digits (leading zeros allowed).
// Used for X-EXTERNAL-ID (36 digits) and referenceNo (20 digits).
func NewNumericID(n int) string {
	limit := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
	v, err := rand.Int(rand.Reader, limit)
	if err != nil {
		return strings.Repeat("0", n)
	}
	s := v.String()
	if len(s) < n {
		s = strings.Repeat("0", n-len(s)) + s
	}
	return s
}
