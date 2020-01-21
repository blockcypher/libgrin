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

// NewCuckatooCtx instantiates a new CuckatooContext as a PowContext
func NewCuckatooCtx(chainType consensus.ChainType, edgeBits uint8, proofSize int, maxSols uint32) *CuckatooContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckatooContext{chainType, params}
}

// CuckatooContext is a Cuckatoo solver context.
type CuckatooContext struct {
	chainType consensus.ChainType
	params    CuckooParams
}

// SetHeaderNonce sets the header nonce.
func (c *CuckatooContext) SetHeaderNonce(header []uint8, nonce *uint32) {
	c.params.resetHeaderNonce(header, nonce)
}

// Return siphash masked for type.
func (c *CuckatooContext) sipnode(edge, uorv uint64) uint64 {
	return c.params.sipnode(edge, uorv, false)
}

// Verify verifies the Cuckatoo context.
func (c *CuckatooContext) Verify(proof Proof) error {
	if proof.proofSize() != consensus.ChainTypeProofSize(c.chainType) {
		return errors.New("wrong cycle length")
	}
	nonces := proof.Nonces
	uvs := make([]uint64, 2*proof.proofSize())
	xor0 := (uint64(c.params.proofSize) / 2) & 1
	xor1 := xor0

	for n := 0; n < proof.proofSize(); n++ {
		if nonces[n] > c.params.edgeMask {
			return errors.New("edge too big")
		}
		if n > 0 && nonces[n] <= nonces[n-1] {
			return errors.New("edges not ascending")
		}
		uvs[2*n] = c.sipnode(nonces[n], 0)
		uvs[2*n+1] = c.sipnode(nonces[n], 1)
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
			if uvs[k]>>1 == uvs[i]>>1 {
				// find other edge endpoint matching one at i
				if j != i {
					return errors.New("branch in cycle")
				}
				j = k
			}
		}
		if j == i || uvs[j] == uvs[i] {
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
