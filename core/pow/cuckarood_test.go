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

// empty header, nonce 64
var v1_19HashCuckarood = [4]uint64{
	0x89f81d7da5e674df,
	0x7586b93105a5fd13,
	0x6fbe212dd4e8c001,
	0x8800c93a8431f938,
}
var v1_19SolCuckarood = []uint64{
	0xa00, 0x3ffb, 0xa474, 0xdc27, 0x182e6, 0x242cc, 0x24de4, 0x270a2, 0x28356, 0x2951f,
	0x2a6ae, 0x2c889, 0x355c7, 0x3863b, 0x3bd7e, 0x3cdbc, 0x3ff95, 0x430b6, 0x4ba1a, 0x4bd7e,
	0x4c59f, 0x4f76d, 0x52064, 0x5378c, 0x540a3, 0x5af6b, 0x5b041, 0x5e9d3, 0x64ec7, 0x6564b,
	0x66763, 0x66899, 0x66e80, 0x68e4e, 0x69133, 0x6b20a, 0x6c2d7, 0x6fd3b, 0x79a8a, 0x79e29,
	0x7ae52, 0x7defe,
}

// empty header, nonce 15
var v2_29HashCuckarood = [4]uint64{
	0xe2f917b2d79492ed,
	0xf51088eaaa3a07a0,
	0xaf4d4288d36a4fa8,
	0xc8cdfd30a54e0581,
}
var v2_29SolCuckarood = []uint64{
	0x1a9629, 0x1fb257, 0x5dc22a, 0xf3d0b0, 0x200c474, 0x24bd68f, 0x48ad104, 0x4a17170,
	0x4ca9a41, 0x55f983f, 0x6076c91, 0x6256ffc, 0x63b60a1, 0x7fd5b16, 0x985bff8, 0xaae71f3,
	0xb71f7b4, 0xb989679, 0xc09b7b8, 0xd7601da, 0xd7ab1b6, 0xef1c727, 0xf1e702b, 0xfd6d961,
	0xfdf0007, 0x10248134, 0x114657f6, 0x11f52612, 0x12887251, 0x13596b4b, 0x15e8d831,
	0x16b4c9e5, 0x17097420, 0x1718afca, 0x187fc40c, 0x19359788, 0x1b41d3f1, 0x1bea25a7,
	0x1d28df0f, 0x1ea6c4a0, 0x1f9bf79f, 0x1fa005c6,
}

var zero [42]uint64

func TestCuckarood19Vectors(t *testing.T) {
	proof := new(Proof)
	ctx := newCuckaroodImpl(consensus.Mainnet, 19, 42)
	ctx.params.siphashKeys = v1_19HashCuckarood
	assert.Nil(t, ctx.Verify(proof.new(v1_19SolCuckarood)))
	assert.NotNil(t, ctx.Verify(proof.zero(42)))
}

func TestCuckarood29Vectors(t *testing.T) {
	proof := new(Proof)
	ctx := newCuckaroodImpl(consensus.Mainnet, 29, 42)
	ctx.params.siphashKeys = v2_29HashCuckarood
	assert.Nil(t, ctx.Verify(proof.new(v2_29SolCuckarood)))
	assert.NotNil(t, ctx.Verify(proof.zero(42)))
}

func newCuckaroodImpl(chainType consensus.ChainType, edgeBits uint8, proofSize int) *CuckaroodContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckaroodContext{chainType, params}
}
