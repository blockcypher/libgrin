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

// NewCuckaroomCtx instantiates a new CuckaroomContext as a PowContext. Note that this can't
/// be moved in the PoWContext trait as this particular trait needs to be
/// convertible to an object trait.
func NewCuckaroomCtx(chainType consensus.ChainType, edgeBits uint8, proofSize int) *CuckaroomContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckaroomContext{chainType, params}
}

// CuckaroomContext is a Cuckaroom cycle context. Only includes the verifier for now.
type CuckaroomContext struct {
	chainType consensus.ChainType
	params    CuckooParams
}

// SetHeaderNonce sets the header nonce.
func (c *CuckaroomContext) SetHeaderNonce(header []uint8, nonce *uint32) {
	c.params.resetHeaderNonce(header, nonce)
}

// Verify verifies the Cuckaroom context.
func (c *CuckaroomContext) Verify(proof Proof) error {
	if proof.proofSize() != consensus.ChainTypeProofSize(c.chainType) {
		return errors.New("wrong cycle length")
	}
	nonces := proof.Nonces
	from := make([]uint32, proof.proofSize())
	to := make([]uint32, proof.proofSize())
	var xorFrom uint32 = 0
	var xorTo uint32 = 0
	nodemask := c.params.edgeMask >> 1

	for n := 0; n < proof.proofSize(); n++ {
		if nonces[n] > c.params.edgeMask {
			return errors.New("edge too big")
		}
		if n > 0 && nonces[n] <= nonces[n-1] {
			return errors.New("edges not ascending")
		}
		edge := SipHashBlock(c.params.siphashKeys, nonces[n], 21, true)
		from[n] = uint32(edge & nodemask)
		xorFrom ^= from[n]
		to[n] = uint32((edge >> 32) & nodemask)
		xorTo ^= to[n]
	}
	if xorFrom != xorTo {
		return errors.New("endpoints don't match up")
	}
	visited := make([]bool, proof.proofSize())
	n := 0
	i := 0
	for {
		// follow cycle
		if visited[i] {
			return errors.New("branch in cycle")
		}
		visited[i] = true
		nexti := 0
		for from[nexti] != to[i] {
			nexti++
			if nexti == proof.proofSize() {
				return errors.New("cycle dead ends")
			}
		}
		i = nexti
		n++
		if i == 0 {
			// must cycle back to start or find branch
			break
		}
	}
	if n == c.params.proofSize {
		return nil
	}
	return errors.New("cycle too short")
}
