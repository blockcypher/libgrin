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

	"github.com/blockcypher/libgrin/v4/core/consensus"
)

// NewCuckaroozCtx instantiates a new CuckaroozContext as a PowContext. Note that this can't
/// be moved in the PoWContext trait as this particular trait needs to be
/// convertible to an object trait.
func NewCuckaroozCtx(chainType consensus.ChainType, edgeBits uint8, proofSize int) (*CuckaroozContext, error) {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, edgeBits+1, proofSize)
	return &CuckaroozContext{chainType, params}, nil
}

// CuckaroozContext is a Cuckarooz cycle context. Only includes the verifier for now.
type CuckaroozContext struct {
	chainType consensus.ChainType
	params    CuckooParams
}

// SetHeaderNonce sets the header nonce.
func (c *CuckaroozContext) SetHeaderNonce(header []uint8, nonce *uint32) {
	c.params.resetHeaderNonce(header, nonce)
}

// Verify verifies the Cuckaroom context.
func (c *CuckaroozContext) Verify(proof Proof) error {
	if proof.proofSize() != consensus.ChainTypeProofSize(c.chainType) {
		return errors.New("wrong cycle length")
	}
	nonces := proof.Nonces
	uvs := make([]uint64, 2*proof.proofSize())
	var xoruv uint64 = 0

	for n := 0; n < proof.proofSize(); n++ {
		if nonces[n] > c.params.edgeMask {
			return errors.New("edge too big")
		}
		if n > 0 && nonces[n] <= nonces[n-1] {
			return errors.New("edges not ascending")
		}
		// 21 is standard siphash rotation constant
		edge := SipHashBlock(c.params.siphashKeys, nonces[n], 21, true)
		uvs[2*n] = edge & c.params.nodeMask
		uvs[2*n+1] = (edge >> 32) & c.params.nodeMask
		xoruv ^= uvs[2*n] ^ uvs[2*n+1]
	}
	if xoruv != 0 {
		return errors.New("endpoints don't match up")
	}

	n := 0
	i := 0
	j := 0
	for {
		// follow cycle
		j = i
		k := j
		for {
			k = (k + 1) % (2 * c.params.proofSize)
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
