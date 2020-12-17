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

package core

import (
	"github.com/blockcypher/libgrin/v5/core/consensus"
	"github.com/blockcypher/libgrin/v5/core/pow"
)

const maxSols uint32 = 10

func createPoWContext(chainType consensus.ChainType, height uint64, edgeBits uint8, proofSize int, nonces []uint64, maxSols uint32) (pow.PowContext, error) {
	if chainType == consensus.Mainnet || chainType == consensus.Testnet {
		// Mainnet and Testnet have Cuckatoo31+ for AF and Cuckaroo{,d,m,z}29 for AR
		if edgeBits > 29 {
			return pow.NewCuckatooCtx(chainType, edgeBits, proofSize, maxSols)
		}
		switch consensus.HeaderVersion(chainType, height) {
		case 1:
			return pow.NewCuckarooCtx(chainType, edgeBits, proofSize)
		case 2:
			return pow.NewCuckaroodCtx(chainType, edgeBits, proofSize)
		case 3:
			return pow.NewCuckaroomCtx(chainType, edgeBits, proofSize)
		case 4:
			return pow.NewCuckaroozCtx(chainType, edgeBits, proofSize)
		default:
			return pow.NoCuckarooCtx()
		}
	}
	// Everything else is Cuckatoo only
	return pow.NewCuckatooCtx(chainType, edgeBits, proofSize, maxSols)
}

// VerifySize validates the proof of work of a given header, and that the proof of work
// satisfies the requirements of the header.
func VerifySize(chainType consensus.ChainType, prePoW []uint8, bh *BlockHeader) error {
	ctx, err := createPoWContext(chainType, bh.Height, bh.PoW.EdgeBits(), len(bh.PoW.Proof.Nonces), bh.PoW.Proof.Nonces, maxSols)
	if err != nil {
		return err
	}
	ctx.SetHeaderNonce(prePoW, nil)
	if err := ctx.Verify(bh.PoW.Proof); err != nil {
		return err
	}
	return nil
}
