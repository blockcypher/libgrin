// Copyright 2019 BlockCypher
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

// SlateV0 represents a slate version 0
type SlateV0 struct {
	// The number of participants intended to take part in this transaction
	NumParticipants uint `json:"num_participants"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	Transaction TransactionV0 `json:"tx"`
	// base amount (excluding fee)
	Amount uint64 `json:"amount"`
	// fee amount
	Fee uint64 `json:"fee"`
	// Block height for the transaction
	Height uint64 `json:"height"`
	// Lock height
	LockHeight uint64 `json:"lock_height"`
	// Participant data, each participant in the transaction will
	// insert their public data here. For now, 0 is sender and 1
	// is receiver, though this will change for multi-party
	ParticipantData []ParticipantDataV0 `json:"participant_data"`
}

// ParticipantDataV0 is the participant data slate version 0
type ParticipantDataV0 struct {
	// Id of participant in the transaction. (For now, 0=sender, 1=rec)
	ID uint64 `json:"id"`
	// Public key corresponding to private blinding factor
	PublicBlindExcess JSONableSlice `json:"public_blind_excess"`
	// Public key corresponding to private nonce
	PublicNonce JSONableSlice `json:"public_nonce"`
	// Public partial signature
	PartSig JSONableSlice `json:"part_sig"`
	// A message for other participants
	Message *string `json:"message"`
	// Signature, created with private key corresponding to 'public_blind_excess'
	MessageSig JSONableSlice `json:"message_sig"`
}

// TransactionV0 is a v0 transaction
type TransactionV0 struct {
	/// The kernel "offset" k2
	/// excess is k1G after splitting the key k = k1 + k2
	Offset JSONableSlice `json:"offset"`
	/// The transaction body - inputs/outputs/kernels
	Body TransactionBodyV0 `json:"body"`
}

// TransactionBodyV0 represent a v0 transaction body
type TransactionBodyV0 struct {
	// List of inputs spent by the transaction.
	Inputs []InputV0 `json:"inputs"`
	// List of outputs the transaction produces.
	Outputs []OutputV0 `json:"outputs"`
	// List of kernels that make up this transaction (usually a single kernel).
	Kernels []TxKernelV0 `json:"kernels"`
}

// InputV0 is a v0 input
type InputV0 struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features core.OutputFeatures `json:"features"`
	// The commit referencing the output being spent.
	Commit JSONableSlice `json:"commit"`
}

// OutputV0 is a v0 output
type OutputV0 struct {
	// Options for an output's structure or use
	Features core.OutputFeatures `json:"features"`
	// The homomorphic commitment representing the output amount
	Commit JSONableSlice `json:"commit"`
	// A proof that the commitment is in the right range
	Proof JSONableSlice `json:"proof"`
}

// TxKernelV0 is a v0 tx kernel
type TxKernelV0 struct {
	// Options for a kernel's structure or use
	Features core.KernelFeatures `json:"features"`
	// Fee originally included in the transaction this proof is for.
	Fee uint64 `json:"fee"`
	// This kernel is not valid earlier than lock_height blocks
	// The max lock_height of all *inputs* to this transaction
	LockHeight uint64 `json:"lock_height"`
	// Remainder of the sum of all transaction commitments. If the transaction
	// is well formed, amounts components should sum to zero and the excess
	// is hence a valid public key.
	Excess JSONableSlice `json:"excess"`
	// The signature proving the excess is a valid public key, which signs
	// the transaction fee.
	ExcessSig JSONableSlice `json:"excess_sig"`
}
