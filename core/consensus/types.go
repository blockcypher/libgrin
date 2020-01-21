// Copyright 2020 BlockCypher
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

package consensus

import "math"

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
