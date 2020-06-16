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
	"bytes"
	"encoding/json"

	"github.com/blockcypher/libgrin/core"

	"github.com/google/uuid"
)

// SlateV4 slate v4
type SlateV4 struct {
	// Versioning info
	VersionInfo VersionCompatInfoV4 `json:"version_info"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	// Slate state
	Sta SlateStateV4 `json:"sta"`
	// Offset, modified by each participant inserting inputs
	Offset string `json:"offset"`
	// The number of participants intended to take part in this transaction, optional
	NumParts *uint32 `json:"num_parts"`
	// base amount (excluding fee)
	Amt core.Uint64 `json:"amt"`
	// fee amount
	Fee core.Uint64 `json:"fee"`
	// kernel features, if any
	Feat uint8 `json:"feat"`
	// TTL, the block height at which wallets
	// should refuse to process the transaction and unlock all
	TTL core.Uint64 `json:"ttl"`
	// Structs always required
	// Participant data, each participant in the transaction will
	// insert their public data here. For now, 0 is sender and 1
	// is receiver, though this will change for multi-party
	Sigs []ParticipantDataV4 `json:"sigs"`
	// Situational, but required at some point in the tx
	// Inputs/Output commits added to slate
	Coms *[]CommitsV4 `json:"coms"`
	// Optional Structs
	// Payment Proof
	Proof *PaymentInfoV4 `json:"proof"`
	// Kernel features arguments
	FeatArgs *KernelFeaturesArgsV4 `json:"feat_args"`
}

// SlateStateV4 state definition
type SlateStateV4 int

const (
	// UnknownSlateState coming from earlier versions of the slate
	UnknownSlateState SlateStateV4 = iota
	// Standard1SlateState flow, freshly init
	Standard1SlateState
	// Standard2SlateState flow, return journey
	Standard2SlateState
	// Standard3SlateState flow, ready for transaction posting
	Standard3SlateState
	// Invoice1SlateState flow, freshly init
	Invoice1SlateState
	// Invoice2SlateState flow, return journey
	Invoice2SlateState
	// Invoice3SlateState flow, ready for tranasction posting
	Invoice3SlateState
)

func (s SlateStateV4) String() string {
	return toStringSlateStateV4[s]
}

var toStringSlateStateV4 = map[SlateStateV4]string{
	UnknownSlateState:   "NA",
	Standard1SlateState: "S1",
	Standard2SlateState: "S2",
	Standard3SlateState: "S3",
	Invoice1SlateState:  "I1",
	Invoice2SlateState:  "I2",
	Invoice3SlateState:  "I3",
}

var toIDSlateStateV4 = map[string]SlateStateV4{
	"NA": UnknownSlateState,
	"S1": Standard1SlateState,
	"S2": Standard2SlateState,
	"S3": Standard3SlateState,
	"I1": Invoice1SlateState,
	"I2": Invoice2SlateState,
	"I3": Invoice3SlateState,
}

// MarshalJSON marshals the enum as a quoted json string
func (s SlateStateV4) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringSlateStateV4[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *SlateStateV4) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'UnknownSlateState' in this case.
	*s = toIDSlateStateV4[j]
	return nil
}

// KernelFeaturesArgsV4 are the kernel features arguments definition
type KernelFeaturesArgsV4 struct {
	/// Lock height, for HeightLocked
	LockHgt core.Uint64 `json:"lock_hgt"`
}

// VersionCompatInfoV4 is a V4 version compat info
type VersionCompatInfoV4 struct {
	// The current version of the slate format
	Version uint16 `json:"version"`
	// The grin block header version this slate is intended for
	BlockHeaderVersion uint16 `json:"block_header_version"`
}

// ParticipantDataV4 is a v4 participant data
type ParticipantDataV4 struct {
	// Public key corresponding to private blinding factor
	Xs string `json:"xs"`
	// Public key corresponding to private nonce
	Nonce string `json:"nonce"`
	// Public partial signature
	Part *string `json:"part"`
}

// PaymentInfoV4 is a v4 payment info
type PaymentInfoV4 struct {
	Saddr string  `json:"saddr"`
	Raddr string  `json:"raddr"`
	Rsig  *string `json:"rsig"`
}

// CommitsV4 is a v4 commit
type CommitsV4 struct {
	// Options for an output's structure or use
	F OutputFeaturesV4 `json:"f"`
	///The homomorphic commitment representing the output amount
	C string `json:"c"`
	// A proof that the commitment is in the right range
	// Only applies for transaction outputs
	P *string `json:"p"`
}

// OutputFeaturesV4 is a v4 output features
type OutputFeaturesV4 struct{ uint8 }

// TransactionV4 is a v4 transaction
type TransactionV4 struct {
	// The kernel "offset" k2
	// excess is k1G after splitting the key k = k1 + k2
	Offset string `json:"offset"`
	// The transaction body - inputs/outputs/kernels
	Body TransactionBodyV4 `json:"body"`
}

// TransactionBodyV4 is a common abstraction for transaction and block
type TransactionBodyV4 struct {
	// List of inputs spent by the transaction.
	Ins []InputV4 `json:"ins"`
	// List of outputs the transaction produces.
	Outs []OutputV4 `json:"outs"`
	// List of kernels that make up this transaction (usually a single kernel).
	Kers []TxKernelV4 `json:"kers"`
}

// InputV4 is a v4 input
type InputV4 struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features OutputFeaturesV4 `json:"features"`
	// The commit referencing the output being spent.
	Commit string `json:"commit"`
}

// OutputV4 is a v4 output
type OutputV4 struct {
	// Options for an output's structure or use
	Features OutputFeaturesV4 `json:"features"`
	// The homomorphic commitment representing the output amount
	Com string `json:"com"`
	// A proof that the commitment is in the right range
	Prf string `json:"prf"`
}

// TxKernelV4 is a v4 tx kernel
type TxKernelV4 struct {
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
