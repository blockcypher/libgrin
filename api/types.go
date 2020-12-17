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

package api

import (
	"bytes"
	"encoding/json"

	"github.com/blockcypher/libgrin/v5/core"
)

// BlockPrintable is the result of the Grin block API
type BlockPrintable struct {
	Header  BlockHeaderPrintable  `json:"header"`
	Inputs  []string              `json:"inputs"`
	Outputs []OutputPrintable     `json:"outputs"`
	Kernels []TxKernelsPrintables `json:"kernels"`
}

// BlockHeaderPrintable is the header of the BlockPrintable
type BlockHeaderPrintable struct {
	// Hash
	Hash string `json:"hash"`
	// Version of the block
	Version uint16 `json:"version"`
	// Height of this block since the genesis block (height 0)
	Height uint64 `json:"height"`
	// Hash of the block previous to this in the chain.
	Previous string `json:"previous"`
	// Root hash of the header MMR at the previous header.
	PrevRoot string `json:"prev_root"`
	// rfc3339 timestamp at which the block was built.
	Timestamp string `json:"timestamp"`
	// Merklish root of all the commitments in the TxHashSet
	OutputRoot string `json:"output_root"`
	// Merklish root of all range proofs in the TxHashSet
	RangeProofRoot string `json:"range_proof_root"`
	// Merklish root of all transaction kernels in the TxHashSet
	KernelRoot string `json:"kernel_root"`
	// Nonce increment used to mine this block.
	Nonce uint64 `json:"nonce"`
	// Size of the cuckoo graph
	EdgeBits uint8 `json:"edge_bits"`
	// Nonces of the cuckoo solution
	CuckooSolution []uint64 `json:"cuckoo_solution"`
	// Total accumulated difficulty since genesis block
	TotalDifficulty uint64 `json:"total_difficulty"`
	// Network secondary PoW factor or factor to use
	SecondaryScaling uint64 `json:"secondary_scaling"`
	// Total kernel offset since genesis block
	TotalKernelOffset string `json:"total_kernel_offset"`
}

// OutputPrintable represents the output of a block
type OutputPrintable struct {
	// The type of output Coinbase|Transaction
	OutputType outputType `json:"output_type"`
	// The homomorphic commitment representing the output's amount
	// (as hex string)
	Commit string `json:"commit"`
	// Whether the output has been spent
	Spent bool `json:"spent"`
	// Rangeproof (as hex string)
	Proof *string `json:"proof"`
	// Rangeproof hash (as hex string)
	ProofHash string `json:"proof_hash"`
	// BlockHeight at which the output is found
	BlockHeight *uint64 `json:"block_height"`
	// Merkle Proof
	MerkleProof *string `json:"merkle_proof"`
	// MMR Position
	MMRIndex uint64 `json:"mmr_index"`
}

// TxKernelsPrintables is the tx kernel
type TxKernelsPrintables struct {
	Features   string `json:"features"`
	Fee        uint64 `json:"fee"`
	LockHeight uint64 `json:"lock_height"`
	Excess     string `json:"excess"`
	ExcessSig  string `json:"excess_sig"`
}

// The Status represents various statistics about the network
type Status struct {
	ProtocolVersion uint32  `json:"protocol_version"`
	UserAgent       string  `json:"user_agent"`
	Connections     uint32  `json:"connections"`
	Tip             Tip     `json:"tip"`
	SyncStatus      string  `json:"sync_status"`
	SyncInfo        *string `json:"sync_info"`
}

// LocatedTxKernel is a located TxKernel
type LocatedTxKernel struct {
	TxKernel core.TxKernel `json:"tx_kernel"`
	Height   uint64        `json:"height"`
	MMRIndex uint64        `json:"mmr_index"`
}

// outputType is the type of output
type outputType int

const (
	// CoinbaseOutputType is a coinbase output type
	CoinbaseOutputType outputType = iota
	// TransactionOutputType is a transaction output type
	TransactionOutputType
)

var toStringOutputType = map[outputType]string{
	CoinbaseOutputType:    "Coinbase",
	TransactionOutputType: "Transaction",
}

var toIDOutputType = map[string]outputType{
	"Coinbase":    CoinbaseOutputType,
	"Transaction": TransactionOutputType,
}

// MarshalJSON marshals the enum as a quoted json string
func (s outputType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringOutputType[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *outputType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'CoinbaseOutputType' in this case.
	*s = toIDOutputType[j]
	return nil
}

// OutputListing is for traversing all outputs in the UTXO set with the
// transactions in the block
type OutputListing struct {
	// The last available output index
	HighestIndex uint64 `json:"highest_index"`
	// The last insertion index retrieved
	LastRetrievedIndex uint64 `json:"last_retrieved_index"`
	// A printable version of the outputs
	Outputs []OutputPrintable `json:"outputs"`
}

// Tip is the state of the current fork tip
type Tip struct {
	/// Height of the tip (max height of the fork)
	Height uint64 `json:"height"`
	// Last block pushed to the fork
	LastBlockPushed string `json:"last_block_pushed"`
	// Block previous to last
	PrevBlockToLast string `json:"prev_block_to_last"`
	// Total difficulty accumulated on that fork
	TotalDifficulty uint64 `json:"total_difficulty"`
}

// Version is the API Version Information
type Version struct {
	// Current node API Version (api crate version)
	NodeVersion string `json:"node_version"`
	// Block header version
	BlockHeaderVersion uint16 `json:"block_header_version"`
}
