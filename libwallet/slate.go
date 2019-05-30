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

package libwallet

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/blockcypher/libgrin/core"
)

// ParticipantData is a public data for each participant in the slate
type ParticipantData struct {
	// Id of participant in the transaction. (For now, 0=sender, 1=rec)
	ID uint64 `json:"id"`
	// Public key corresponding to private blinding factor
	PublicBlindExcess JSONableSlice `json:"public_blind_excess"`
	// Public key corresponding to private nonce
	PublicNonce JSONableSlice `json:"public_nonce"`
	// Public partial signature
	PartSig *JSONableSlice `json:"part_sig"`
	// A message for other participants
	Message *string `json:"message"`
	// Signature, created with private key corresponding to 'public_blind_excess'
	MessageSig JSONableSlice `json:"message_sig"`
}

// ParticipantMessageData is the public message data (for serializing and storage)
type ParticipantMessageData struct {
	// id of the particpant in the tx
	ID uint64 `json:"id"`
	// Public key
	PublicKey JSONableSlice `json:"public_key"`
	// Message,
	Message *string `json:"message"`
	// Signature
	MessageSig *JSONableSlice `json:"message_sig"`
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
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	Transaction core.Transaction `json:"tx"`
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
	ParticipantData []ParticipantData `json:"participant_data"`
}

// VersionCompatInfo is the versioning and compatibility info about this slate
type VersionCompatInfo struct {
	// The current version of the slate format
	Version uint16 `json:"version"`
	// Original version this slate was converted from
	OrigVersion uint16 `json:"orig_version"`
	// The grin block header version this slate is intended for
	BlockHeaderVersion uint16 `json:"block_header_version"`
}

// ParticipantMessages is an helper just to facilitate serialization
type ParticipantMessages struct {
	// included messages
	Messages []ParticipantMessageData `json:"messages"`
}

// TODO Replace this by Signature and public key
// JSONableSlice is a slice that is not represented as a base58 when serialized
type JSONableSlice []uint8

// MarshalJSON is the marshal function for such type
func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}
