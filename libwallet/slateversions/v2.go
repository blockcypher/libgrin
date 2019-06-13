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
	"encoding/hex"

	"github.com/blockcypher/libgrin/core"
	"github.com/google/uuid"
)

//! Contains V2 of the slate (grin-wallet 1.1.0)
//! Changes from V1:
//! * ParticipantData struct fields serialized as hex strings instead of arrays:
//!    * public_blind_excess
//!    * public_nonce
//!    * part_sig
//!    * message_sig
//! * Transaction fields serialized as hex strings instead of arrays:
//!    * offset
//! * Input field serialized as hex strings instead of arrays:
//!    commit
//! * Output fields serialized as hex strings instead of arrays:
//!    commit
//!    proof
//! * TxKernel fields serialized as hex strings instead of arrays:
//!    commit
//!    signature
//! * version field removed
//! * VersionCompatInfo struct created with fields and added to beginning of struct
//!    version: u16
//!    orig_version: u16,
//!    block_header_version: u16,

// SlateV2 represents a slate version 2
type SlateV2 struct {
	// Versioning info
	VersionInfo VersionCompatInfoV2 `json:"version_info"`
	// The number of participants intended to take part in this transaction
	NumParticipants uint `json:"num_participants"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	Transaction TransactionV2 `json:"tx"`
	// base amount (excluding fee)
	Amount uint64 `json:"amount,string"`
	// fee amount
	Fee uint64 `json:"fee,string"`
	// Block height for the transaction
	Height uint64 `json:"height,string"`
	// Lock height
	LockHeight uint64 `json:"lock_height,string"`
	// Participant data, each participant in the transaction will
	// insert their public data here. For now, 0 is sender and 1
	// is receiver, though this will change for multi-party
	ParticipantData []ParticipantDataV2 `json:"participant_data"`
}

// VersionCompatInfoV2 is a V2 version compat info
type VersionCompatInfoV2 struct {
	// The current version of the slate format
	Version uint16 `json:"version"`
	// Original version this slate was converted from
	OrigVersion uint16 `json:"orig_version"`
	// The grin block header version this slate is intended for
	BlockHeaderVersion uint16 `json:"block_header_version"`
}

// ParticipantDataV2 is the participant data slate version 0
type ParticipantDataV2 struct {
	// Id of participant in the transaction. (For now, 0=sender, 1=rec)
	ID uint64 `json:"id,string"`
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

// TransactionV2 is a v1 transaction
type TransactionV2 struct {
	/// The kernel "offset" k2
	/// excess is k1G after splitting the key k = k1 + k2
	Offset string `json:"offset"`
	/// The transaction body - inputs/outputs/kernels
	Body TransactionBodyV2 `json:"body"`
}

// TransactionBodyV2 represent a v1 transaction body
type TransactionBodyV2 struct {
	// List of inputs spent by the transaction.
	Inputs []InputV2 `json:"inputs"`
	// List of outputs the transaction produces.
	Outputs []OutputV2 `json:"outputs"`
	// List of kernels that make up this transaction (usually a single kernel).
	Kernels []TxKernelV2 `json:"kernels"`
}

// InputV2 is a v1 input
type InputV2 struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features core.OutputFeatures `json:"features"`
	// The commit referencing the output being spent.
	Commit string `json:"commit"`
}

// OutputV2 is a v1 output
type OutputV2 struct {
	// Options for an output's structure or use
	Features core.OutputFeatures `json:"features"`
	// The homomorphic commitment representing the output amount
	Commit string `json:"commit"`
	// A proof that the commitment is in the right range
	Proof string `json:"proof"`
}

// TxKernelV2 is a v1 tx kernel
type TxKernelV2 struct {
	// Options for a kernel's structure or use
	Features core.KernelFeatures `json:"features"`
	// Fee originally included in the transaction this proof is for.
	Fee uint64 `json:"fee,string"`
	// This kernel is not valid earlier than lock_height blocks
	// The max lock_height of all *inputs* to this transaction
	LockHeight uint64 `json:"lock_height,string"`
	// Remainder of the sum of all transaction commitments. If the transaction
	// is well formed, amounts components should sum to zero and the excess
	// is hence a valid public key.
	Excess string `json:"excess"`
	// The signature proving the excess is a valid public key, which signs
	// the transaction fee.
	ExcessSig string `json:"excess_sig"`
}

// Upgrade V1 to V2
func (v1 *SlateV1) Upgrade() SlateV2 {
	var participantDataV2 []ParticipantDataV2
	for i := range v1.ParticipantData {
		participantDataV2 = append(participantDataV2, v1.ParticipantData[i].upgrade())
	}
	var version uint16 = 2
	origVersion := v1.origVersion
	var blockHeaderVersion uint16 = 1
	versionInfo := VersionCompatInfoV2{
		Version:            version,
		OrigVersion:        uint16(origVersion),
		BlockHeaderVersion: blockHeaderVersion,
	}
	slateV2 := SlateV2{
		VersionInfo:     versionInfo,
		NumParticipants: v1.NumParticipants,
		ID:              v1.ID,
		Transaction:     v1.Transaction.upgrade(),
		Amount:          v1.Amount,
		Fee:             v1.Fee,
		Height:          v1.Height,
		LockHeight:      v1.LockHeight,
		ParticipantData: participantDataV2,
	}
	return slateV2
}

func (v1 *ParticipantDataV1) upgrade() ParticipantDataV2 {
	var partSig *string
	if v1.PartSig != nil {
		hexPartSig := hex.EncodeToString(v1.PartSig)
		partSig = &hexPartSig
	}
	var messageSig *string
	if v1.MessageSig != nil {
		hexMessageSig := hex.EncodeToString(v1.MessageSig)
		messageSig = &hexMessageSig
	}
	participantDataV2 := ParticipantDataV2{
		ID:                v1.ID,
		PublicBlindExcess: hex.EncodeToString(v1.PublicBlindExcess),
		PublicNonce:       hex.EncodeToString(v1.PublicNonce),
		PartSig:           partSig,
		Message:           v1.Message,
		MessageSig:        messageSig,
	}
	return participantDataV2
}

func (v1 *TransactionV1) upgrade() TransactionV2 {
	transactionV2 := TransactionV2{
		Offset: hex.EncodeToString(v1.Offset),
		Body:   v1.Body.upgrade(),
	}
	return transactionV2
}

func (v1 *TransactionBodyV1) upgrade() TransactionBodyV2 {
	var inputsV2 []InputV2
	var outputV2 []OutputV2
	var kernelsV2 []TxKernelV2
	for i := range v1.Inputs {
		inputsV2 = append(inputsV2, v1.Inputs[i].upgrade())
	}
	for i := range v1.Outputs {
		outputV2 = append(outputV2, v1.Outputs[i].upgrade())
	}
	for i := range v1.Kernels {
		kernelsV2 = append(kernelsV2, v1.Kernels[i].upgrade())
	}
	transactionBodyV2 := TransactionBodyV2{
		Inputs:  inputsV2,
		Outputs: outputV2,
		Kernels: kernelsV2,
	}
	return transactionBodyV2
}

func (v1 *InputV1) upgrade() InputV2 {
	inputV2 := InputV2{
		Features: v1.Features,
		Commit:   hex.EncodeToString(v1.Commit),
	}
	return inputV2
}

func (v1 *OutputV1) upgrade() OutputV2 {
	outputV2 := OutputV2{
		Features: v1.Features,
		Commit:   hex.EncodeToString(v1.Commit),
		Proof:    hex.EncodeToString(v1.Proof),
	}
	return outputV2
}

func (v1 *TxKernelV1) upgrade() TxKernelV2 {
	txKernelV2 := TxKernelV2{
		Features:   v1.Features,
		Fee:        v1.Fee,
		LockHeight: v1.LockHeight,
		Excess:     hex.EncodeToString(v1.Excess),
		ExcessSig:  hex.EncodeToString(v1.ExcessSig),
	}
	return txKernelV2
}

// Downgrade V2 to V1
func (v2 *SlateV2) Downgrade() SlateV1 {
	var participantDataV1 []ParticipantDataV1
	for i := range v2.ParticipantData {
		participantDataV1 = append(participantDataV1, v2.ParticipantData[i].downgrade())
	}
	slateV1 := SlateV1{
		NumParticipants: v2.NumParticipants,
		ID:              v2.ID,
		Transaction:     v2.Transaction.downgrade(),
		Amount:          v2.Amount,
		Fee:             v2.Fee,
		Height:          v2.Height,
		LockHeight:      v2.LockHeight,
		ParticipantData: participantDataV1,
		Version:         1,
		origVersion:     2,
	}
	return slateV1
}

func (v2 *ParticipantDataV2) downgrade() ParticipantDataV1 {
	publicBlindExcess, _ := hex.DecodeString(v2.PublicBlindExcess)
	publicNonce, _ := hex.DecodeString(v2.PublicNonce)
	var partSig JSONableSlice
	if v2.PartSig != nil {
		partSig, _ = hex.DecodeString(*v2.PartSig)
	}
	var messageSig JSONableSlice
	if v2.MessageSig != nil {
		messageSig, _ = hex.DecodeString(*v2.MessageSig)
	}
	participantDataV1 := ParticipantDataV1{
		ID:                v2.ID,
		PublicBlindExcess: publicBlindExcess,
		PublicNonce:       publicNonce,
		PartSig:           partSig,
		Message:           v2.Message,
		MessageSig:        messageSig,
	}
	return participantDataV1
}

func (v2 *TransactionV2) downgrade() TransactionV1 {
	offset, _ := hex.DecodeString(v2.Offset)
	transactionV1 := TransactionV1{
		Offset: offset,
		Body:   v2.Body.downgrade(),
	}
	return transactionV1
}

func (v2 *TransactionBodyV2) downgrade() TransactionBodyV1 {
	var inputsV1 []InputV1
	var outputV1 []OutputV1
	var kernelsV1 []TxKernelV1
	for i := range v2.Inputs {
		inputsV1 = append(inputsV1, v2.Inputs[i].downgrade())
	}
	for i := range v2.Outputs {
		outputV1 = append(outputV1, v2.Outputs[i].downgrade())
	}
	for i := range v2.Kernels {
		kernelsV1 = append(kernelsV1, v2.Kernels[i].downgrade())
	}
	transactionBodyV1 := TransactionBodyV1{
		Inputs:  inputsV1,
		Outputs: outputV1,
		Kernels: kernelsV1,
	}
	return transactionBodyV1
}

func (v2 *InputV2) downgrade() InputV1 {
	commit, _ := hex.DecodeString(v2.Commit)

	inputV1 := InputV1{
		Features: v2.Features,
		Commit:   commit,
	}
	return inputV1
}

func (v2 *OutputV2) downgrade() OutputV1 {
	commit, _ := hex.DecodeString(v2.Commit)
	proof, _ := hex.DecodeString(v2.Proof)
	outputV1 := OutputV1{
		Features: v2.Features,
		Commit:   commit,
		Proof:    proof,
	}
	return outputV1
}

func (v2 *TxKernelV2) downgrade() TxKernelV1 {
	excess, _ := hex.DecodeString(v2.Excess)
	excessSig, _ := hex.DecodeString(v2.ExcessSig)
	txKernelV1 := TxKernelV1{
		Features:   v2.Features,
		Fee:        v2.Fee,
		LockHeight: v2.LockHeight,
		Excess:     excess,
		ExcessSig:  excessSig,
	}
	return txKernelV1
}
