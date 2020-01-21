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

import "github.com/blockcypher/libgrin/core/pow"

// BlockHeader is a block header, fairly standard compared to other blockchains.
type BlockHeader struct {
	// Version of the block
	Version uint16
	// Height of this block since the genesis block (height 0)
	Height uint64
	// Hash of the block previous to this in the chain.
	PrevHash string
	// Root hash of the header MMR at the previous header.
	PrevRoot string
	// Timestamp at which the block was built.
	Timestamp string
	// Merklish root of all the commitments in the TxHashSet
	OutputRoot string
	// Merklish root of all range proofs in the TxHashSet
	RangeProofRoot string
	// Merklish root of all transaction kernels in the TxHashSet
	KernelRoot string
	// Total accumulated sum of kernel offsets since genesis block.
	// We can derive the kernel offset sum for *this* block from
	// the total kernel offset of the previous block header.
	TotalKernelOffset string
	// Total size of the output MMR after applying this block
	OutputMmrSize uint64
	// Total size of the kernel MMR after applying this block
	KernelMmrSize uint64
	// Proof of work and related
	PoW pow.ProofOfWork
}
