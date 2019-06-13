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

//! Contains V1 of the slate (grin 1.0.1, 1.0.2)
//! Changes from V0:
//! * Addition of a version field to Slate struct

// SlateV1 represents a slate version 1
type SlateV1 struct {
	// The number of participants intended to take part in this transaction
	NumParticipants uint `json:"num_participants"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	Transaction TransactionV1 `json:"tx"`
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
	ParticipantData []ParticipantDataV1 `json:"participant_data"`
	/// Version
	Version uint64 `json:"version"`
	// Unexported since just used locally
	origVersion uint64
}

// GetOrigVersion get the origin version
func (v1 *SlateV1) GetOrigVersion() uint64 {
	return v1.origVersion
}

// GetOrigVersion get the origin version
func (v1 *SlateV1) SetOrigVersion(v uint64) {
	v1.origVersion = v
}

// ParticipantDataV1 is the participant data slate version 0
type ParticipantDataV1 struct {
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

// TransactionV1 is a v1 transaction
type TransactionV1 struct {
	/// The kernel "offset" k2
	/// excess is k1G after splitting the key k = k1 + k2
	Offset JSONableSlice `json:"offset"`
	/// The transaction body - inputs/outputs/kernels
	Body TransactionBodyV1 `json:"body"`
}

// TransactionBodyV1 represent a v1 transaction body
type TransactionBodyV1 struct {
	// List of inputs spent by the transaction.
	Inputs []InputV1 `json:"inputs"`
	// List of outputs the transaction produces.
	Outputs []OutputV1 `json:"outputs"`
	// List of kernels that make up this transaction (usually a single kernel).
	Kernels []TxKernelV1 `json:"kernels"`
}

// InputV1 is a v1 input
type InputV1 struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features core.OutputFeatures `json:"features"`
	// The commit referencing the output being spent.
	Commit JSONableSlice `json:"commit"`
}

// OutputV1 is a v1 output
type OutputV1 struct {
	// Options for an output's structure or use
	Features core.OutputFeatures `json:"features"`
	// The homomorphic commitment representing the output amount
	Commit JSONableSlice `json:"commit"`
	// A proof that the commitment is in the right range
	Proof JSONableSlice `json:"proof"`
}

// TxKernelV1 is a v1 tx kernel
type TxKernelV1 struct {
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

// Upgrade V0 to V1
func (v0 *SlateV0) Upgrade() SlateV1 {
	var participantDataV1 []ParticipantDataV1
	for i := range v0.ParticipantData {
		participantDataV1 = append(participantDataV1, v0.ParticipantData[i].upgrade())
	}
	slateV1 := SlateV1{
		NumParticipants: v0.NumParticipants,
		ID:              v0.ID,
		Transaction:     v0.Transaction.upgrade(),
		Amount:          v0.Amount,
		Fee:             v0.Fee,
		Height:          v0.Height,
		LockHeight:      v0.LockHeight,
		ParticipantData: participantDataV1,
		Version:         1,
		origVersion:     0,
	}
	return slateV1
}

func (v0 *ParticipantDataV0) upgrade() ParticipantDataV1 {
	participantDataV1 := ParticipantDataV1{
		ID:                v0.ID,
		PublicBlindExcess: v0.PublicBlindExcess,
		PublicNonce:       v0.PublicNonce,
		PartSig:           v0.PartSig,
		Message:           v0.Message,
		MessageSig:        v0.MessageSig,
	}
	return participantDataV1
}

func (v0 *TransactionV0) upgrade() TransactionV1 {
	transactionV1 := TransactionV1{
		Offset: v0.Offset,
		Body:   v0.Body.upgrade(),
	}
	return transactionV1
}

func (v0 *TransactionBodyV0) upgrade() TransactionBodyV1 {
	var inputsV1 []InputV1
	var outputV1 []OutputV1
	var kernelsV1 []TxKernelV1
	for i := range v0.Inputs {
		inputsV1 = append(inputsV1, v0.Inputs[i].upgrade())
	}
	for i := range v0.Outputs {
		outputV1 = append(outputV1, v0.Outputs[i].upgrade())
	}
	for i := range v0.Kernels {
		kernelsV1 = append(kernelsV1, v0.Kernels[i].upgrade())
	}
	transactionBodyV1 := TransactionBodyV1{
		Inputs:  inputsV1,
		Outputs: outputV1,
		Kernels: kernelsV1,
	}
	return transactionBodyV1
}

func (v0 *InputV0) upgrade() InputV1 {
	inputV1 := InputV1{
		Features: v0.Features,
		Commit:   v0.Commit,
	}
	return inputV1
}

func (v0 *OutputV0) upgrade() OutputV1 {
	outputV1 := OutputV1{
		Features: v0.Features,
		Commit:   v0.Commit,
		Proof:    v0.Proof,
	}
	return outputV1
}

func (v0 *TxKernelV0) upgrade() TxKernelV1 {
	txKernelV1 := TxKernelV1{
		Features:   v0.Features,
		Fee:        v0.Fee,
		LockHeight: v0.LockHeight,
		Excess:     v0.Excess,
		ExcessSig:  v0.ExcessSig,
	}
	return txKernelV1
}

// Downgrade V1 to V0
func (v1 *SlateV1) Downgrade() SlateV0 {
	var participantDataV0 []ParticipantDataV0
	for i := range v1.ParticipantData {
		participantDataV0 = append(participantDataV0, v1.ParticipantData[i].downgrade())
	}
	slateV0 := SlateV0{
		NumParticipants: v1.NumParticipants,
		ID:              v1.ID,
		Transaction:     v1.Transaction.downgrade(),
		Amount:          v1.Amount,
		Fee:             v1.Fee,
		Height:          v1.Height,
		LockHeight:      v1.LockHeight,
		ParticipantData: participantDataV0,
	}
	return slateV0
}

func (v1 *ParticipantDataV1) downgrade() ParticipantDataV0 {
	participantDataV0 := ParticipantDataV0{
		ID:                v1.ID,
		PublicBlindExcess: v1.PublicBlindExcess,
		PublicNonce:       v1.PublicNonce,
		PartSig:           v1.PartSig,
		Message:           v1.Message,
		MessageSig:        v1.MessageSig,
	}
	return participantDataV0
}

func (v1 *TransactionV1) downgrade() TransactionV0 {
	transactionV0 := TransactionV0{
		Offset: v1.Offset,
		Body:   v1.Body.downgrade(),
	}
	return transactionV0
}

func (v1 *TransactionBodyV1) downgrade() TransactionBodyV0 {
	var inputsV0 []InputV0
	var outputV0 []OutputV0
	var kernelsV0 []TxKernelV0
	for i := range v1.Inputs {
		inputsV0 = append(inputsV0, v1.Inputs[i].downgrade())
	}
	for i := range v1.Outputs {
		outputV0 = append(outputV0, v1.Outputs[i].downgrade())
	}
	for i := range v1.Kernels {
		kernelsV0 = append(kernelsV0, v1.Kernels[i].downgrade())
	}
	transactionBodyV0 := TransactionBodyV0{
		Inputs:  inputsV0,
		Outputs: outputV0,
		Kernels: kernelsV0,
	}
	return transactionBodyV0
}

func (v1 *InputV1) downgrade() InputV0 {
	inputV0 := InputV0{
		Features: v1.Features,
		Commit:   v1.Commit,
	}
	return inputV0
}

func (v1 *OutputV1) downgrade() OutputV0 {
	outputV0 := OutputV0{
		Features: v1.Features,
		Commit:   v1.Commit,
		Proof:    v1.Proof,
	}
	return outputV0
}

func (v1 *TxKernelV1) downgrade() TxKernelV0 {
	txKernelV0 := TxKernelV0{
		Features:   v1.Features,
		Fee:        v1.Fee,
		LockHeight: v1.LockHeight,
		Excess:     v1.Excess,
		ExcessSig:  v1.ExcessSig,
	}
	return txKernelV0
}
