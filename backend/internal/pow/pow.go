// Package pow encodes and measures proofs of work as defined in docs/spec/claims.md section 2.
package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
)

const (
	version      = "hp1"
	preimageSize = 24
	maxWork      = 256
)

// Proof is a claim on cell (X, Y) with Color, at the client-chosen time TS (Unix seconds).
type Proof struct {
	X     uint16
	Y     uint16
	Color uint8
	TS    int64
	Nonce uint64
}

// Preimage returns the byte-exact preimage of p.
func Preimage(p Proof) [preimageSize]byte {
	var b [preimageSize]byte

	copy(b[:], version)
	binary.BigEndian.PutUint16(b[3:], p.X)
	binary.BigEndian.PutUint16(b[5:], p.Y)
	b[7] = p.Color
	binary.BigEndian.PutUint64(b[8:], uint64(p.TS)) //nolint:gosec // reason: int64be is the two's complement bit pattern
	binary.BigEndian.PutUint64(b[16:], p.Nonce)

	return b
}

// Hash returns the SHA-256 of the preimage of p.
func Hash(p Proof) [sha256.Size]byte {
	preimage := Preimage(p)

	return sha256.Sum256(preimage[:])
}

// Work returns 256 - log2(H) in bits, where H is hash read as a big-endian integer, or 256 if H = 0.
func Work(hash [sha256.Size]byte) float64 {
	// float64 holds H to 53 significant bits, far more precision than work needs.
	h := 0.0
	for _, b := range hash {
		h = h*256 + float64(b)
	}

	if h == 0 {
		return maxWork
	}

	return maxWork - math.Log2(h)
}
