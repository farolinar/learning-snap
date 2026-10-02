package snap

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// MinifyJSON returns the compact, no-whitespace JSON encoding of v.
// Encoding/json emits compact JSON by default, which is exactly what SNAP means by "minify".
func MinifyJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

// SHA256HexLower returns the lowercase hex-encoded SHA-256 digest of b
func SHA256HexLower(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
