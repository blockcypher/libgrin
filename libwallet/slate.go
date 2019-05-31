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
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"

	"github.com/blockcypher/libgrin/core"
	"github.com/blockcypher/libgrin/libwallet/slateversions"
)

// ParticipantData is a public data for each participant in the slate
type ParticipantData struct {
	// Id of participant in the transaction. (For now, 0=sender, 1=rec)
	ID uint64 `json:"id"`
	// Public key corresponding to private blinding factor
	PublicBlindExcess string `json:"public_blind_excess"`
	// Public key corresponding to private nonce
	PublicNonce string `json:"public_nonce"`
	// Public partial signature
	PartSig *string `json:"part_sig"`
	// A message for other participants
	Message *string `json:"message"`
	// Signature, created with private key corresponding to 'public_blind_excess'
	MessageSig string `json:"message_sig"`
}

// ParticipantMessageData is the public message data (for serializing and storage)
type ParticipantMessageData struct {
	// id of the particpant in the tx
	ID uint64 `json:"id"`
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

func parseSlateVersion(slateBytes []byte) (uint16, error) {
	var version uint16
	slate := make(map[string]interface{})
	err := json.Unmarshal(slateBytes, &slate)
	if err != nil {
		return 0, err
	}
	if versionCompatInfo, ok := slate["version_info"].(VersionCompatInfo); !ok {
		return versionCompatInfo.Version, nil
	}
	if _, ok := slate["version"].(uint16); !ok {
		return 1, nil
	}
	return version, nil
}

func DeserializeUpgrade(slateBytes []byte) (*Slate, error) {
	// check version
	version, err := parseSlateVersion(slateBytes)
	if err != nil {
		return nil, errors.New("Can't parse slate version")
	}

	switch version {
	case 2:
		var slate Slate
		if err := json.Unmarshal(slateBytes, &slate); err != nil {
			return nil, err
		}
		return &slate, nil
	case 1:
		var v1 slateversions.SlateV1
		if err := json.Unmarshal(slateBytes, &v1); err != nil {
			return nil, err
		}
		v1.OrigVersion = 1
		return nil, errors.New("Can't parse slate version")
	case 0:
		var v0 slateversions.SlateV0
		if err := json.Unmarshal(slateBytes, &v0); err != nil {
			return nil, err
		}
		return nil, errors.New("Can't parse slate version")
	default:
		return nil, errors.New("Can't parse slate version")
	}
}
