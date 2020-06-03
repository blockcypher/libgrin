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

var v1_19HashCuckarooz = [4]uint64{
	0xd129f63fba4d9a85,
	0x457dcb3666c5e09c,
	0x045247a2e2ee75f7,
	0x1a0f2e1bcb9d93ff,
}

var v1_19SolCuckarooz = []uint64{
	0x33b6, 0x487b, 0x88b7, 0x10bf6, 0x15144, 0x17cb7, 0x22621, 0x2358e, 0x23775, 0x24fb3,
	0x26b8a, 0x2876c, 0x2973e, 0x2f4ba, 0x30a62, 0x3a36b, 0x3ba5d, 0x3be67, 0x3ec56, 0x43141,
	0x4b9c5, 0x4fa06, 0x51a5c, 0x523e5, 0x53d08, 0x57d34, 0x5c2de, 0x60bba, 0x62509, 0x64d69,
	0x6803f, 0x68af4, 0x6bd52, 0x6f041, 0x6f900, 0x70051, 0x7097d, 0x735e8, 0x742c2, 0x79ae5,
	0x7f64d, 0x7fd49,
}

var v2_29HashCuckarooz = [4]uint64{
	0x34bb4c75c929a2f5,
	0x21df13263aa81235,
	0x37d00939eae4be06,
	0x473251cbf6941553,
}

var v2_29SolCuckarooz = []uint64{
	0x49733a, 0x1d49107, 0x253d2ca, 0x5ad5e59, 0x5b671bd, 0x5dcae1c, 0x5f9a589, 0x65e9afc,
	0x6a59a45, 0x7d9c6d3, 0x7df96e4, 0x8b26174, 0xa17b430, 0xa1c8c0d, 0xa8a0327, 0xabd7402,
	0xacb7c77, 0xb67524f, 0xc1c15a6, 0xc7e2c26, 0xc7f5d8d, 0xcae478a, 0xdea9229, 0xe1ab49e,
	0xf57c7db, 0xfb4e8c5, 0xff314aa, 0x110ccc12, 0x143e546f, 0x17007af8, 0x17140ea2,
	0x173d7c5d, 0x175cd13f, 0x178b8880, 0x1801edc5, 0x18c8f56b, 0x18c8fe6d, 0x19f1a31a,
	0x1bb028d1, 0x1caaa65a, 0x1cf29bc2, 0x1dbde27d,
}

func TestCuckarooz19Vectors(t *testing.T) {
	proof := new(Proof)
	ctx := newCuckaroomImpl(consensus.Mainnet, 19, 42)
	ctx.params.siphashKeys = v1_19HashCuckaroom
	assert.Nil(t, ctx.Verify(proof.new(v1_19SolCuckaroom)))
	assert.NotNil(t, ctx.Verify(proof.zero(42)))
}

func TestCuckarooz29Vectors(t *testing.T) {
	proof := new(Proof)
	ctx := newCuckaroomImpl(consensus.Mainnet, 29, 42)
	ctx.params.siphashKeys = v2_29HashCuckaroom
	assert.Nil(t, ctx.Verify(proof.new(v2_29SolCuckaroom)))
	assert.NotNil(t, ctx.Verify(proof.zero(42)))
}

func newCuckaroozImpl(chainType consensus.ChainType, edgeBits uint8, proofSize int) *CuckaroomContext {
	cp := new(CuckooParams)
	params := cp.new(edgeBits, proofSize)
	return &CuckaroomContext{chainType, params}
}
