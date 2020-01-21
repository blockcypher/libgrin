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
	"testing"

	"github.com/blockcypher/libgrin/core/consensus"
	"github.com/stretchr/testify/assert"
)

// Cuckatoo 29 Solution for Header [0u8;80] - nonce 20
var v1_29 = []uint64{
	0x48a9e2, 0x9cf043, 0x155ca30, 0x18f4783, 0x248f86c, 0x2629a64, 0x5bad752, 0x72e3569,
	0x93db760, 0x97d3b37, 0x9e05670, 0xa315d5a, 0xa3571a1, 0xa48db46, 0xa7796b6, 0xac43611,
	0xb64912f, 0xbb6c71e, 0xbcc8be1, 0xc38a43a, 0xd4faa99, 0xe018a66, 0xe37e49c, 0xfa975fa,
	0x11786035, 0x1243b60a, 0x12892da0, 0x141b5453, 0x1483c3a0, 0x1505525e, 0x1607352c,
	0x16181fe3, 0x17e3a1da, 0x180b651e, 0x1899d678, 0x1931b0bb, 0x19606448, 0x1b041655,
	0x1b2c20ad, 0x1bd7a83c, 0x1c05d5b0, 0x1c0b9caa,
}

// Cuckatoo 31 Solution for Header [0u8;80] - nonce 99
var v1_31 = []uint64{
	0x1128e07, 0xc181131, 0x110fad36, 0x1135ddee, 0x1669c7d3, 0x1931e6ea, 0x1c0005f3, 0x1dd6ecca,
	0x1e29ce7e, 0x209736fc, 0x2692bf1a, 0x27b85aa9, 0x29bb7693, 0x2dc2a047, 0x2e28650a, 0x2f381195,
	0x350eb3f9, 0x3beed728, 0x3e861cbc, 0x41448cc1, 0x41f08f6d, 0x42fbc48a, 0x4383ab31, 0x4389c61f,
	0x4540a5ce, 0x49a17405, 0x50372ded, 0x512f0db0, 0x588b6288, 0x5a36aa46, 0x5c29e1fe, 0x6118ab16,
	0x634705b5, 0x6633d190, 0x6683782f, 0x6728b6e1, 0x67adfb45, 0x68ae2306, 0x6d60f5e1, 0x78af3c4f,
	0x7dde51ab, 0x7faced21,
}

func TestValidate29Vectors(t *testing.T) {
	ctx := newCuckatooImpl(consensus.Mainnet, 29, 42, 10)
	nonce := uint32(20)
	ctx.SetHeaderNonce(make([]uint8, 80), &nonce)
	proof := new(Proof)
	assert.Nil(t, ctx.Verify(proof.new(v1_29)))
}

func TestValidate31Vectors(t *testing.T) {
	ctx := newCuckatooImpl(consensus.Mainnet, 31, 42, 10)
	nonce := uint32(99)
	ctx.SetHeaderNonce(make([]uint8, 80), &nonce)
	proof := new(Proof)
	assert.Nil(t, ctx.Verify(proof.new(v1_31)))
}

func TestValidateFail(t *testing.T) {
	proof := new(Proof)
	ctx := newCuckatooImpl(consensus.Mainnet, 29, 42, 10)
	header := make([]uint8, 80)
	header[0] = uint8(1)
	nonce := uint32(20)
	ctx.SetHeaderNonce(header, &nonce)
	assert.NotNil(t, ctx.Verify(proof.new(v1_29)))
	header[0] = uint8(0)
	ctx.SetHeaderNonce(header, &nonce)
	assert.Nil(t, ctx.Verify(proof.new(v1_29)))
	badProof := v1_29
	badProof[0] = 0x48a9e1
	assert.NotNil(t, ctx.Verify(proof.new(badProof)))
}

func newCuckatooImpl(chainType consensus.ChainType, edgeBits uint8, proofSize int, maxSols uint32) CuckatooContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return CuckatooContext{chainType, params}
}
