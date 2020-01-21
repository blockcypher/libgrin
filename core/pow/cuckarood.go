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

// NewCuckaroodCtx instantiates a new CuckaroodContext as a PowContext. Note that this can't
// be moved in the PoWContext trait as this particular trait needs to be
// convertible to an object trait.
func NewCuckaroodCtx(chainType consensus.ChainType, edgeBits uint8, proofSize int) *CuckaroodContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckaroodContext{chainType, params}
}

// CuckaroodContext is a Cuckatoo cycle context. Only includes the verifier for now.
type CuckaroodContext struct {
	chainType consensus.ChainType
	params    CuckooParams
}

// SetHeaderNonce sets the header nonce.
func (c *CuckaroodContext) SetHeaderNonce(header []uint8, nonce *uint32) {
	c.params.resetHeaderNonce(header, nonce)
}

// Verify verifies the Cuckatoo context.
func (c *CuckaroodContext) Verify(proof Proof) error {
	if proof.proofSize() != consensus.ChainTypeProofSize(c.chainType) {
		return errors.New("wrong cycle length")
	}
	nonces := proof.Nonces
	uvs := make([]uint64, 2*proof.proofSize())
	ndir := make([]uint64, 2)
	var xor0, xor1 uint64
	nodemask := c.params.edgeMask >> 1

	for n := 0; n < proof.proofSize(); n++ {
		dir := uint(nonces[n] & 1)
		if ndir[dir] >= uint64(proof.proofSize())/2 {
			return errors.New("edges not balanced")
		}
		if nonces[n] > c.params.edgeMask {
			return errors.New("edge too big")
		}
		if n > 0 && nonces[n] <= nonces[n-1] {
			return errors.New("edges not ascending")
		}
		edge := SipHashBlock(c.params.siphashKeys, nonces[n], 25, false)
		idx := 4*ndir[dir] + 2*uint64(dir)
		uvs[idx] = edge & nodemask
		uvs[idx+1] = (edge >> 32) & nodemask
		xor0 ^= uvs[idx]
		xor1 ^= uvs[idx+1]
		ndir[dir]++
	}

	if xor0|xor1 != 0 {
		return errors.New("endpoints don't match up")
	}
	var i, j, n int

	for {
		// follow cycle
		j = i
		for k := ((i % 4) ^ 2); k < 2*c.params.proofSize; k += 4 {
			if uvs[k] == uvs[i] {
				// find reverse edge endpoint identical to one at i
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
