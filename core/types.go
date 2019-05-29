// Copyright 2019 BlockCypher
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// PowContext is a generic interface for a solver/verifier providing common
// interface into Cuckoo-family PoW
type PowContext interface {
	// Sets the header along with an optional nonce at the end solve: whether to
	// set up structures for a solve (true) or just validate (false)
	SetHeaderNonce(header []uint8, nonce *uint32)
	// Verify a solution with the stored parameters
	Verify(proof Proof) error
}

// Proof is a Cuck(at)oo Cycle proof of work, consisting of the edge_bits to get
// the graph size (i.e. the 2-log of the number of edges) and the nonces of the
// graph solution. While being expressed as u64 for simplicity, nonces a.k.a.
// edge indices range from 0 to (1 << edge_bits) - 1
//
// The hash of the `Proof` is the hash of its packed nonces when serializing
// them at their exact bit size. The resulting bit sequence is padded to be
// byte-aligned.
type Proof struct {
	// Power of 2 used for the size of the cuckoo graph
	EdgeBits uint8
	// The nonces
	Nonces []uint64
}

func (p *Proof) new(inNonces []uint64) Proof {
	// No sorting here
	return Proof{EdgeBits: DefaultMinEdgeBits, Nonces: inNonces}
}

// Builds a proof with all bytes zeroed out
func (p *Proof) zero(proofSize int) Proof {
	return Proof{DefaultMinEdgeBits, make([]uint64, proofSize)}
}

// Returns the proof size
func (p *Proof) proofSize() int {
	return len(p.Nonces)
}

// Difficulty achieved by this proof with given scaling factor
func (p *Proof) scaledDifficulty(blockHashString string, scaleUint64 uint64) uint64 {
	hash, _ := strconv.ParseUint(blockHashString[:16], 16, 64)
	var scale big.Int
	scale = *scale.SetUint64(scaleUint64)
	scaleShifted := scale.Lsh(&scale, 64)

	var maxUint64 big.Int
	maxUint64 = *maxUint64.SetUint64(math.MaxUint64)

	var diff big.Int
	diff = *diff.Div(scaleShifted, big.NewInt(int64(hash)))
	if diff.Cmp(&maxUint64) == 1 {
		return math.MaxUint64
	}
	return diff.Uint64()
}

// ProofOfWork is a block header information pertaining to the proof of work
type ProofOfWork struct {
	// Total accumulated difficulty since genesis block
	TotalDifficulty uint64
	// Variable difficulty scaling factor fo secondary proof of work
	SecondaryScaling uint32
	// Nonce increment used to mine this block.
	Nonce uint64
	// Proof of work data.
	Proof Proof
}

// EdgeBits is the edge bits used for the cuckoo cycle size on this proof
func (p *ProofOfWork) EdgeBits() uint8 {
	return p.Proof.EdgeBits
}

// IsPrimary is whether this proof of work is for the primary algorithm (as opposed / to
//secondary). Only depends on the edge_bits at this time.
func (p *ProofOfWork) IsPrimary() bool {
	// 2 conditions are redundant right now but not necessarily in the future
	return p.Proof.EdgeBits != SecondPoWEdgeBits && p.Proof.EdgeBits >= DefaultMinEdgeBits
}

// IsSecondary is whether this proof of work is for the secondary algorithm (as opposed / to
// primary). Only depends on the edge_bits at this time.
func (p *ProofOfWork) IsSecondary() bool {
	return p.Proof.EdgeBits == SecondPoWEdgeBits
}

// The Difficulty is defined as the maximum target divided by the block hash.
type Difficulty struct {
	num uint64
}

func (p *Difficulty) zero() Difficulty {
	return Difficulty{num: 0}
}

// Difficulty of MIN_DIFFICULTY
func (p *Difficulty) min() Difficulty {
	return Difficulty{num: MinDifficulty}
}

// Unit is a difficulty unit, which is the graph weight of minimal graph
func (p *Difficulty) Unit(chainType ChainType) Difficulty {
	return Difficulty{num: uint64(initialGraphWeight(chainType))}
}

// FromNum converts a `uint64` into a `Difficulty`
func (p *Difficulty) FromNum(num uint64) Difficulty {
	// can't have difficulty lower than 1
	return Difficulty{num: uint64(math.Max(float64(num), 1))}
}

// JSONableSlice is a slice that is not represented as a base58 when serialized
type JSONableSlice []uint8

// MarshalJSON is the marshal function for such type
func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}

// saturatingSubUint8 is a saturating uint8 subtraction. Computes a - b, saturating at the numeric bounds instead of overflowing.
func saturatingSubUint8(a, b uint8) uint8 {
	if a < b {
		return 0
	}
	return a - b
}

// saturatingSubUint16 is a saturating uint16 subtraction. Computes a - b, saturating at the numeric bounds instead of overflowing.
func saturatingSubUint16(a, b uint16) uint16 {
	if a < b {
		return 0
	}
	return a - b
}

// saturatingSubUint32 is a saturating uint32 subtraction. Computes a - b, saturating at the numeric bounds instead of overflowing.
func saturatingSubUint32(a, b uint32) uint32 {
	if a < b {
		return 0
	}
	return a - b
}

// saturatingSubUint64 is a saturating uint64 subtraction. Computes a - b, saturating at the numeric bounds instead of overflowing.
func saturatingSubUint64(a, b uint64) uint64 {
	if a < b {
		return 0
	}
	return a - b
}

// saturatingAddUint8 is a saturating uint8 addition. Computes a + b, saturating at the numeric bounds instead of overflowing.
func saturatingAddUint8(a, b uint8) uint8 {
	c := a + b
	if c < a {
		return math.MaxUint8
	}
	return a + b
}

// saturatingAddUint16 is a saturating uint16 addition. Computes a + b, saturating at the numeric bounds instead of overflowing.
func saturatingAddUint16(a, b uint16) uint16 {
	c := a + b
	if c < a {
		return math.MaxUint16
	}
	return a + b
}

// saturatingAddUint32 is a saturating uint32 addition. Computes a + b, saturating at the numeric bounds instead of overflowing.
func saturatingAddUint32(a, b uint32) uint32 {
	c := a + b
	if c < a {
		return math.MaxUint32
	}
	return a + b
}

// saturatingAddUint64 is a saturating uint64 addition. Computes a + b, saturating at the numeric bounds instead of overflowing.
func saturatingAddUint64(a, b uint64) uint64 {
	c := a + b
	if c < a {
		return math.MaxUint64
	}
	return a + b
}
