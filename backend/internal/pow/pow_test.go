package pow_test

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/nem0z/hashplace/backend/internal/pow"
)

const (
	vectorsPath   = "../../../docs/spec/vectors/pow-v1.json"
	workTolerance = 1e-9
)

// vector is one entry of the golden file. TS and Nonce may be JSON numbers or strings.
type vector struct {
	X        uint16      `json:"x"`
	Y        uint16      `json:"y"`
	Color    uint8       `json:"color"`
	TS       json.Number `json:"ts"`
	Nonce    json.Number `json:"nonce"`
	Preimage string      `json:"preimage"`
	Hash     string      `json:"hash"`
	Work     float64     `json:"work"`
}

func loadVectors(t *testing.T) []vector {
	t.Helper()

	data, err := os.ReadFile(vectorsPath)
	if err != nil {
		t.Fatal(err)
	}

	var vectors []vector
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}

	if len(vectors) == 0 {
		t.Fatal("no vectors")
	}

	return vectors
}

func (v vector) proof(t *testing.T) pow.Proof {
	t.Helper()

	ts, err := strconv.ParseInt(v.TS.String(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	nonce, err := strconv.ParseUint(v.Nonce.String(), 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	return pow.Proof{X: v.X, Y: v.Y, Color: v.Color, TS: ts, Nonce: nonce}
}

func TestVectors(t *testing.T) {
	t.Parallel()

	for _, v := range loadVectors(t) {
		t.Run(v.Preimage, func(t *testing.T) {
			t.Parallel()

			p := v.proof(t)

			preimage := pow.Preimage(p)
			if got := hex.EncodeToString(preimage[:]); got != v.Preimage {
				t.Errorf("preimage = %s, want %s", got, v.Preimage)
			}

			hash := pow.Hash(p)
			if got := hex.EncodeToString(hash[:]); got != v.Hash {
				t.Errorf("hash = %s, want %s", got, v.Hash)
			}

			if got := pow.Work(hash); math.Abs(got-v.Work) > workTolerance {
				t.Errorf("work = %v, want %v", got, v.Work)
			}
		})
	}
}

func TestWorkBounds(t *testing.T) {
	t.Parallel()

	one := [32]byte{31: 1}
	half := [32]byte{0: 0x80}

	var maxHash [32]byte
	for i := range maxHash {
		maxHash[i] = 0xff
	}

	tests := []struct {
		name string
		hash [32]byte
		want float64
	}{
		{"zero", [32]byte{}, 256},
		{"one", one, 256},
		{"half", half, 1},
		{"max", maxHash, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pow.Work(tt.hash); math.Abs(got-tt.want) > workTolerance {
				t.Errorf("Work = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzWork(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0xff, 0xff, 0xff, 0xff})
	f.Add(make([]byte, 31))

	f.Fuzz(func(t *testing.T, data []byte) {
		var hash [32]byte
		copy(hash[:], data)

		work := pow.Work(hash)
		if math.IsNaN(work) || work < 0 || work > 256 {
			t.Fatalf("Work(%x) = %v, want a value in 0..256", hash, work)
		}
	})
}

func BenchmarkHashWork(b *testing.B) {
	p := pow.Proof{X: 12, Y: 34, Color: 5, TS: 1791072000}

	for b.Loop() {
		p.Nonce++
		pow.Work(pow.Hash(p))
	}
}
