// Copyright 2018 BlockCypher
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

const maxSols uint32 = 10

// VerifySize validates the proof of work of a given header, and that the proof of work
// satisfies the requirements of the header.
func VerifySize(chainType ChainType, prePoW []uint8, bh *BlockHeader) error {
	ctx := createPoWContext(chainType, bh.Height, bh.PoW.EdgeBits(), len(bh.PoW.Proof.Nonces), bh.PoW.Proof.Nonces, maxSols)
	ctx.SetHeaderNonce(prePoW, nil)
	if err := ctx.Verify(bh.PoW.Proof); err != nil {
		return err
	}
	return nil
}
