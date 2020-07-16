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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/blockcypher/libgrin/core"

	"github.com/google/uuid"
)

const defaultNumParticipants uint8 = 2

// SlateV4 slate v4
type SlateV4 struct {
	// Versioning info
	Ver VersionCompatInfoV4 `json:"ver"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	// Slate state
	Sta slateStateV4 `json:"sta"`
	// Optional fields depending on state
	// Offset, modified by each participant inserting inputs
	Off string `json:"off,omitempty"`
	// The number of participants intended to take part in this transaction, optional
	NumParts uint8 `json:"num_parts"`
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
	Coms *[]CommitsV4 `json:"coms,omitempty"`
	// Optional Structs
	// Payment Proof
	Proof *PaymentInfoV4 `json:"proof,omitempty"`
	// Kernel features arguments
	FeatArgs *KernelFeaturesArgsV4 `json:"feat_args,omitempty"`
}

// MarshalJSON is a custom marshaller to account for default field
func (s SlateV4) MarshalJSON() ([]byte, error) {
	type TempSlateV4 SlateV4
	var tempS TempSlateV4 = TempSlateV4(s)
	bytes, err := json.Marshal(tempS)
	if err != nil {
		return nil, err
	}
	var tempSlateMapV4 map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &tempSlateMapV4); err != nil {
		return nil, err
	}
	// Default numParts is 2
	if s.NumParts == defaultNumParticipants {
		delete(tempSlateMapV4, "num_parts")
	}
	// We could keep the following fields
	// But save some space by having
	if s.Amt == 0 {
		delete(tempSlateMapV4, "amt")
	}
	if s.Fee == 0 {
		delete(tempSlateMapV4, "fee")
	}
	if s.Feat == 0 {
		delete(tempSlateMapV4, "feat")
	}
	if s.TTL == 0 {
		delete(tempSlateMapV4, "ttl")
	}
	return json.Marshal(tempSlateMapV4)
}

// UnmarshalJSON is a custom unmarshaller that respect default value
func (s *SlateV4) UnmarshalJSON(b []byte) error {
	type TempSlateV4 SlateV4
	var tempS TempSlateV4
	if err := json.Unmarshal(b, &tempS); err != nil {
		return err
	}
	// this indicates that the num parts was not in the bytes slice
	// replace by the default num parts
	if tempS.NumParts == 0 {
		tempS.NumParts = defaultNumParticipants
	}
	*s = SlateV4(tempS)
	return nil
}

// SlateStateV4 state definition
type slateStateV4 int

const (
	// UnknownSlateState coming from earlier versions of the slate
	UnknownSlateState slateStateV4 = iota
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

var toStringSlateStateV4 = map[slateStateV4]string{
	UnknownSlateState:   "NA",
	Standard1SlateState: "S1",
	Standard2SlateState: "S2",
	Standard3SlateState: "S3",
	Invoice1SlateState:  "I1",
	Invoice2SlateState:  "I2",
	Invoice3SlateState:  "I3",
}

var toIDSlateStateV4 = map[string]slateStateV4{
	"NA": UnknownSlateState,
	"S1": Standard1SlateState,
	"S2": Standard2SlateState,
	"S3": Standard3SlateState,
	"I1": Invoice1SlateState,
	"I2": Invoice2SlateState,
	"I3": Invoice3SlateState,
}

// MarshalJSON marshals the enum as a quoted json string
func (s slateStateV4) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringSlateStateV4[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *slateStateV4) UnmarshalJSON(b []byte) error {
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

// MarshalJSON marshals the VersionCompatInfoV4 as a quoted version like {}:{}
func (v VersionCompatInfoV4) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%d:%d", v.Version, v.BlockHeaderVersion)
	bytes, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalJSON unmarshals a quoted version to a VersionCompatInfoV4
func (v *VersionCompatInfoV4) UnmarshalJSON(bs []byte) error {
	var verString string
	if err := json.Unmarshal(bs, &verString); err != nil {
		return err
	}
	ver := strings.Split(verString, ":")
	if len(ver) != 2 {
		return errors.New("cannot parse version")
	}

	version, err := strconv.ParseUint(ver[0], 10, 16)
	if err != nil {
		return errors.New("cannot parse version")
	}

	v.Version = uint16(version)
	blockHeaderVersion, err := strconv.ParseUint(ver[1], 10, 16)
	if err != nil {
		return errors.New("cannot parse version")
	}
	v.BlockHeaderVersion = uint16(blockHeaderVersion)
	return nil
}

// ParticipantDataV4 is a v4 participant data
type ParticipantDataV4 struct {
	// Public key corresponding to private blinding factor
	Xs string `json:"xs"`
	// Public key corresponding to private nonce
	Nonce string `json:"nonce"`
	// Public partial signature
	Part *string `json:"part,omitempty"`
}

// PaymentInfoV4 is a v4 payment info
type PaymentInfoV4 struct {
	Saddr string  `json:"saddr"`
	Raddr string  `json:"raddr"`
	Rsig  *string `json:"rsig,omitempty"`
}

// CommitsV4 is a v4 commit
type CommitsV4 struct {
	// Options for an output's structure or use
	F OutputFeaturesV4 `json:"f"`
	// The homomorphic commitment representing the output amount
	C string `json:"c"`
	// A proof that the commitment is in the right range
	// Only applies for transaction outputs
	P *string `json:"p,omitempty"`
}

// OutputFeaturesV4 is a v4 output features
type OutputFeaturesV4 uint8

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
