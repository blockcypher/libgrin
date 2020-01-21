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

package pow

import (
	"errors"

	"github.com/blockcypher/libgrin/core/consensus"
)

// https://github.com/mimblewimble/grin/blob/master/core/src/pow/cuckaroo.rs

// NewCuckarooCtx instantiates a new CuckarooContext as a PowContext
func NewCuckarooCtx(chainType consensus.ChainType, edgeBits uint8, proofSize int) *CuckarooContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckarooContext{chainType, params}
}

// CuckarooContext is a Cuckatoo cycle context. Only includes the verifier for now.
type CuckarooContext struct {
	chainType consensus.ChainType
	params    CuckooParams
}

// SetHeaderNonce sets the header nonce.
func (c *CuckarooContext) SetHeaderNonce(header []uint8, nonce *uint32) {
	c.params.resetHeaderNonce(header, nonce)
}

// Verify verifies the Cuckatoo context.
func (c *CuckarooContext) Verify(proof Proof) error {
	if proof.proofSize() != consensus.ChainTypeProofSize(c.chainType) {
		return errors.New("wrong cycle length")
	}
	nonces := proof.Nonces
	uvs := make([]uint64, 2*proof.proofSize())
	var xor0, xor1 uint64

	for n := 0; n < proof.proofSize(); n++ {
		if nonces[n] > c.params.edgeMask {
			return errors.New("edge too big")
		}
		if n > 0 && nonces[n] <= nonces[n-1] {
			return errors.New("edges not ascending")
		}
		// 21 is standard siphash rotation constant
		edge := SipHashBlock(c.params.siphashKeys, nonces[n], 21, false)
		uvs[2*n] = edge & c.params.edgeMask
		uvs[2*n+1] = (edge >> 32) & c.params.edgeMask
		xor0 ^= uvs[2*n]
		xor1 ^= uvs[2*n+1]
	}

	if xor0|xor1 != 0 {
		return errors.New("endpoints don't match up")
	}

	var i, j, n int
	for {
		// follow cycle
		j = i
		k := j
		for {
			k = (k + 2) % (2 * c.params.proofSize)
			if k == i {
				break
			}
			if uvs[k] == uvs[i] {
				// find other edge endpoint matching one at i
				if j != i {
					return errors.New("branch in cycle")
				}
				j = k
			}
		}
		if j == i {
			return errors.New("cycle dead ends")
		}
		i = j ^ 1
		n++
		if i == 0 {
			break
		}
	}
	if n == c.params.proofSize {
		return nil
	}
	return errors.New("cycle too short")
}
