// Copyright 2020 BlockCypher
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http//www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package slateversions

import (
	"github.com/blockcypher/libgrin/core"
	"github.com/google/uuid"
)

//! Contains V3 of the slate (grin-wallet 3.0.0)
//! Changes from V2:
//! * Addition of payment_proof (PaymentInfo struct)
//! * Addition of a u64 ttl_cutoff_height field

type SlateV3 struct {
	// Versioning info
	VersionInfo VersionCompatInfoV3 `json:"version_info"`
	// The number of participants intended to take part in this transaction
	NumParticipants uint `json:"num_participants"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	Transaction TransactionV3 `json:"tx"`
	// base amount (excluding fee)
	Amount core.Uint64 `json:"amount"`
	// fee amount
	Fee core.Uint64 `json:"fee"`
	// Block height for the transaction
	Height core.Uint64 `json:"height"`
	// Lock height
	LockHeight core.Uint64 `json:"lock_height"`
	// TTL, the block height at which wallets
	// should refuse to process the transaction and unlock all
	// associated outputs
	TTLCutoffHeight *core.Uint64 `json:"ttl_cutoff_height"`
	// Participant data, each participant in the transaction will
	// insert their public data here. For now, 0 is sender and 1
	// is receiver, though this will change for multi-party
	ParticipantData []ParticipantDataV3 `json:"participant_data"`
	// Payment Proof
	PaymentProof *PaymentInfoV3 `json:"payment_proof"`
}

// VersionCompatInfoV3 is a V3 version compat info
type VersionCompatInfoV3 struct {
	// The current version of the slate format
	Version uint16 `json:"version"`
	// Original version this slate was converted from
	OrigVersion uint16 `json:"orig_version"`
	// The grin block header version this slate is intended for
	BlockHeaderVersion uint16 `json:"block_header_version"`
}

// ParticipantDataV3 is the participant data slate version 3
type ParticipantDataV3 struct {
	// Id of participant in the transaction. (For now, 0=sender, 1=rec)
	ID core.Uint64 `json:"id"`
	// Public key corresponding to private blinding factor
	PublicBlindExcess string `json:"public_blind_excess"`
	// Public key corresponding to private nonce
	PublicNonce string `json:"public_nonce"`
	// Public partial signature
	PartSig *string `json:"part_sig"`
	// A message for other participants
	Message *string `json:"message"`
	// Signature, created with private key corresponding to 'public_blind_excess'
	MessageSig *string `json:"message_sig"`
}

// PaymentInfoV3 is a v3 payment info
type PaymentInfoV3 struct {
	SenderAddress     string  `json:"sender_address"`
	ReceiverAddress   string  `json:"receiver_address"`
	ReceiverSignature *string `json:"receiver_signature"`
}

// TransactionV3 is a v3 transaction
type TransactionV3 struct {
	/// The kernel "offset" k2
	/// excess is k1G after splitting the key k = k1 + k2
	Offset string `json:"offset"`
	/// The transaction body - inputs/outputs/kernels
	Body TransactionBodyV3 `json:"body"`
}

// TransactionBodyV3 represent a v3 transaction body
type TransactionBodyV3 struct {
	// List of inputs spent by the transaction.
	Inputs []InputV3 `json:"inputs"`
	// List of outputs the transaction produces.
	Outputs []OutputV3 `json:"outputs"`
	// List of kernels that make up this transaction (usually a single kernel).
	Kernels []TxKernelV3 `json:"kernels"`
}

// InputV3 is a v3 input
type InputV3 struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features core.OutputFeatures `json:"features"`
	// The commit referencing the output being spent.
	Commit string `json:"commit"`
}

// OutputV3 is a v3 output
type OutputV3 struct {
	// Options for an output's structure or use
	Features core.OutputFeatures `json:"features"`
	// The homomorphic commitment representing the output amount
	Commit string `json:"commit"`
	// A proof that the commitment is in the right range
	Proof string `json:"proof"`
}

// TxKernelV3 is a v3 tx kernel
type TxKernelV3 struct {
	// Options for a kernel's structure or use
	Features core.KernelFeatures `json:"features"`
	// Fee originally included in the transaction this proof is for.
	Fee core.Uint64 `json:"fee"`
	// This kernel is not valid earlier than lock_height blocks
	// The max lock_height of all *inputs* to this transaction
	LockHeight core.Uint64 `json:"lock_height"`
	// Remainder of the sum of all transaction commitments. If the transaction
	// is well formed, amounts components should sum to zero and the excess
	// is hence a valid public key.
	Excess string `json:"excess"`
	// The signature proving the excess is a valid public key, which signs
	// the transaction fee.
	ExcessSig string `json:"excess_sig"`
}
