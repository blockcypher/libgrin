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

// empty header, nonce 71
var v1_19Hash = [4]uint64{
	0x23796193872092ea,
	0xf1017d8a68c4b745,
	0xd312bd53d2cd307b,
	0x840acce5833ddc52,
}

var v1_19Sol = []uint64{
	0x45e9, 0x6a59, 0xf1ad, 0x10ef7, 0x129e8, 0x13e58, 0x17936, 0x19f7f, 0x208df, 0x23704,
	0x24564, 0x27e64, 0x2b828, 0x2bb41, 0x2ffc0, 0x304c5, 0x31f2a, 0x347de, 0x39686, 0x3ab6c,
	0x429ad, 0x45254, 0x49200, 0x4f8f8, 0x5697f, 0x57ad1, 0x5dd47, 0x607f8, 0x66199, 0x686c7,
	0x6d5f3, 0x6da7a, 0x6dbdf, 0x6f6bf, 0x6ffbb, 0x7580e, 0x78594, 0x785ac, 0x78b1d, 0x7b80d,
	0x7c11c, 0x7da35,
}

// empty header, nonce 143
var v2_19Hash = [4]uint64{
	0x6a54f2a35ab7e976,
	0x68818717ff5cd30e,
	0x9c14260c1bdbaf7,
	0xea5b4cd5d0de3cf0,
}

var v2_19Sol = []uint64{
	0x2b1e, 0x67d3, 0xb041, 0xb289, 0xc6c3, 0xd31e, 0xd75c, 0x111d7, 0x145aa, 0x1712e, 0x1a3af,
	0x1ecc5, 0x206b1, 0x2a55c, 0x2a9cd, 0x2b67e, 0x321d8, 0x35dde, 0x3721e, 0x37ac0, 0x39edb,
	0x3b80b, 0x3fc79, 0x4148b, 0x42a48, 0x44395, 0x4bbc9, 0x4f775, 0x515c5, 0x56f97, 0x5aa10,
	0x5bc1b, 0x5c56d, 0x5d552, 0x60a2e, 0x66646, 0x6c3aa, 0x70709, 0x71d13, 0x762a3, 0x79d88,
	0x7e3ae,
}

func TestCuckaroo19Vectors(t *testing.T) {
	proof := new(Proof)
	ctx := newCuckarooImpl(consensus.Mainnet, 19, 42)
	ctx.params.siphashKeys = v1_19Hash
	assert.Nil(t, ctx.Verify(proof.new(v1_19Sol)))
	ctx.params.siphashKeys = v2_19Hash
	assert.Nil(t, ctx.Verify(proof.new(v2_19Sol)))
	assert.NotNil(t, ctx.Verify(proof.new(v1_19Sol)))
}

func newCuckarooImpl(chainType consensus.ChainType, edgeBits uint8, proofSize int) *CuckarooContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckarooContext{chainType, params}
}
