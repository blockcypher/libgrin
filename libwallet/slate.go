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

package libwallet

import (
	"bytes"
	"encoding/json"

	"github.com/blockcypher/libgrin/core"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// PaymentInfo is a payment info
type PaymentInfo struct {
	SenderAddress     string  `json:"sender_address"`
	ReceiverAddress   string  `json:"receiver_address"`
	ReceiverSignature *string `json:"receiver_signature"`
}

// ParticipantData is a public data for each participant in the slate
type ParticipantData struct {
	// Public key corresponding to private blinding factor
	PublicBlindExcess string `json:"public_blind_excess"`
	// Public key corresponding to private nonce
	PublicNonce string `json:"public_nonce"`
	// Public partial signature
	PartSig *string `json:"part_sig"`
}

// ParticipantMessageData is the public message data (for serializing and storage)
type ParticipantMessageData struct {
	// id of the particpant in the tx
	ID core.Uint64 `json:"id"`
	// Public key
	PublicKey string `json:"public_key"`
	// Message,
	Message *string `json:"message"`
	// Signature
	MessageSig *string `json:"message_sig"`
}

// A Slate is passed around to all parties to build up all of the public
// transaction data needed to create a finalized transaction. Callers can pass
// the slate around by whatever means they choose, (but we can provide some
// binary or JSON serialization helpers here).
type Slate struct {
	// Versioning info
	VersionInfo VersionCompatInfo `json:"version_info"`
	// The number of participants intended to take part in this transaction
	NumParticipants uint `json:"num_participants"`
	// Unique transaction ID, selected by sender
	ID uuid.UUID `json:"id"`
	/// Slate state
	State SlateState `json:"state"`
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	// Optional as of V4 to allow for a compact
	// transaction initiation
	Transaction *core.Transaction `json:"tx"`
	// base amount (excluding fee)
	Amount core.Uint64 `json:"amount"`
	// fee amount
	Fee core.Uint64 `json:"fee"`
	// TTL, the block height at which wallets
	// should refuse to process the transaction and unlock all
	// associated outputs
	TTLCutoffHeight *core.Uint64 `json:"ttl_cutoff_height"`
	// Kernel Features flag -
	// 	0: plain
	// 	1: coinbase (invalid)
	// 	2: height_locked
	// 	3: NRD
	KernelFeatures uint8 `json:"kernel_features"`
	// Participant data, each participant in the transaction will
	// insert their public data here. For now, 0 is sender and 1
	// is receiver, though this will change for multi-party
	ParticipantData []ParticipantData `json:"participant_data"`
	// Payment Proof
	PaymentProof *PaymentInfo `json:"payment_proof"`
	// Kernel features arguments
	KernelFeaturesArgs *KernelFeaturesArgs `json:"kernel_features_args"`
	// TODO: Remove post HF3
	// participant ID, only stored for compatibility with V3 slates
	// not serialized anywhere
	ParticipantID *string `json:"participant_id"`
}

// SlateState state definition
type SlateState int

const (
	// UnknownSlateState coming from earlier versions of the slate
	UnknownSlateState SlateState = iota
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

func (s SlateState) String() string {
	return toStringSlateState[s]
}

var toStringSlateState = map[SlateState]string{
	UnknownSlateState:   "NA",
	Standard1SlateState: "S1",
	Standard2SlateState: "S2",
	Standard3SlateState: "S3",
	Invoice1SlateState:  "I1",
	Invoice2SlateState:  "I2",
	Invoice3SlateState:  "I3",
}

var toIDSlateState = map[string]SlateState{
	"NA": UnknownSlateState,
	"S1": Standard1SlateState,
	"S2": Standard2SlateState,
	"S3": Standard3SlateState,
	"I1": Invoice1SlateState,
	"I2": Invoice2SlateState,
	"I3": Invoice3SlateState,
}

// MarshalJSON marshals the enum as a quoted json string
func (s SlateState) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringSlateState[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *SlateState) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'UnknownSlateState' in this case.
	*s = toIDSlateState[j]
	return nil
}

// KernelFeaturesArgs are the kernel features arguments definition
type KernelFeaturesArgs struct {
	/// Lock height, for HeightLocked (also NRD relative lock height)
	LockHeight core.Uint64 `json:"lock_hgt"`
}

// VersionCompatInfo is the versioning and compatibility info about this slate
type VersionCompatInfo struct {
	// The current version of the slate format
	Version uint16 `json:"version"`
	// The grin block header version this slate is intended for
	BlockHeaderVersion uint16 `json:"block_header_version"`
}

func parseSlateVersion(slateBytes []byte) (uint16, error) {
	var version uint16
	slate := make(map[string]interface{})
	if err := json.Unmarshal(slateBytes, &slate); err != nil {
		return 0, err
	}
	// First check for version info
	if _, ok := slate["version_info"]; ok {
		var versionCompatInfo VersionCompatInfo
		if err := mapstructure.Decode(slate["version_info"], &versionCompatInfo); err == nil {
			return versionCompatInfo.Version, nil
		}
	}

	if _, ok := slate["version"].(float64); ok {
		return 1, nil
	}
	return version, nil
}
