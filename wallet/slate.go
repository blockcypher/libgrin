// Copyright 2018 BlockCypher
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wallet

import "github.com/blockcypher/libgrin/core"

// A Slate is passed around to all parties to build up all of the public
// transaction data needed to create a finalized transaction. Callers can pass
// the slate around by whatever means they choose, (but we can provide some
// binary or JSON serialization helpers here).
type Slate struct {
	// The number of participants intended to take part in this transaction
	NumParticipants uint `json:"num_participants"`
	// Unique transaction ID, selected by sender
	ID string `json:"id"`
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

// ParticipantData is a public data for each participant in the slate
type ParticipantData struct {
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
