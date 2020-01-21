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
	"encoding/json"
	"errors"

	"github.com/blockcypher/libgrin/core"
	"github.com/blockcypher/libgrin/libwallet/slateversions"
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
	// The core transaction data:
	// inputs, outputs, kernels, kernel offset
	Transaction core.Transaction `json:"tx"`
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
	ParticipantData []ParticipantData `json:"participant_data"`
	// Payment Proof
	PaymentProof *PaymentInfo `json:"payment_proof"`
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

// UnmarshalUpgrade unmarshal and upgrade a slate to v3 if necessary
func UnmarshalUpgrade(slateBytes []byte, slate *Slate) error {
	// check version
	version, err := parseSlateVersion(slateBytes)
	if err != nil {
		return errors.New("can't parse slate version")
	}

	switch version {
	case 3:
		if err := json.Unmarshal(slateBytes, &slate); err != nil {
			return err
		}
		return nil
	case 2:
		var slateV2 slateversions.SlateV2
		if err := json.Unmarshal(slateBytes, &slateV2); err != nil {
			return err
		}
		*slate = slateV2ToSlate(slateV2)
		return nil
	default:
		return errors.New("can't parse slate version")
	}
}

func slateV2ToSlate(v2 slateversions.SlateV2) Slate {
	var slate Slate
	slate.VersionInfo = VersionCompatInfo(v2.VersionInfo)
	slate.NumParticipants = v2.NumParticipants
	slate.ID = v2.ID
	var inputs []core.Input
	for i := range v2.Transaction.Body.Inputs {
		inputs = append(inputs, core.Input(v2.Transaction.Body.Inputs[i]))
	}
	var outputs []core.Output
	for i := range v2.Transaction.Body.Outputs {
		outputs = append(outputs, core.Output(v2.Transaction.Body.Outputs[i]))
	}
	var kernels []core.TxKernel
	for i := range v2.Transaction.Body.Kernels {
		kernels = append(kernels, core.TxKernel(v2.Transaction.Body.Kernels[i]))
	}
	slate.Transaction.Body.Inputs = inputs
	slate.Transaction.Body.Outputs = outputs
	slate.Transaction.Body.Kernels = kernels
	slate.Transaction.Offset = v2.Transaction.Offset
	slate.Amount = v2.Amount
	slate.Fee = v2.Fee
	slate.Height = v2.Height
	slate.LockHeight = v2.LockHeight
	var participantData []ParticipantData
	for i := range v2.ParticipantData {
		participantData = append(participantData, ParticipantData(v2.ParticipantData[i]))
	}
	slate.ParticipantData = participantData
	return slate
}

func slateToSlateV2(slate Slate) slateversions.SlateV2 {
	var slateV2 slateversions.SlateV2
	slateV2.VersionInfo = slateversions.VersionCompatInfoV2(slate.VersionInfo)
	slateV2.NumParticipants = slate.NumParticipants
	slateV2.ID = slate.ID
	var inputs []slateversions.InputV2
	for i := range slate.Transaction.Body.Inputs {
		inputs = append(inputs, slateversions.InputV2(slate.Transaction.Body.Inputs[i]))
	}
	var outputs []slateversions.OutputV2
	for i := range slate.Transaction.Body.Outputs {
		outputs = append(outputs, slateversions.OutputV2(slate.Transaction.Body.Outputs[i]))
	}
	var kernels []slateversions.TxKernelV2
	for i := range slate.Transaction.Body.Kernels {
		kernels = append(kernels, slateversions.TxKernelV2(slate.Transaction.Body.Kernels[i]))
	}
	slateV2.Transaction.Body.Inputs = inputs
	slateV2.Transaction.Body.Outputs = outputs
	slateV2.Transaction.Body.Kernels = kernels
	slateV2.Transaction.Offset = slate.Transaction.Offset
	slateV2.Amount = slate.Amount
	slateV2.Fee = slate.Fee
	slateV2.Height = slate.Height
	slateV2.LockHeight = slate.LockHeight
	var participantData []slateversions.ParticipantDataV2
	for i := range slate.ParticipantData {
		participantData = append(participantData, slateversions.ParticipantDataV2(slate.ParticipantData[i]))
	}
	slateV2.ParticipantData = participantData
	return slateV2
}
