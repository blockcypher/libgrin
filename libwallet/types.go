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
	"time"

	"github.com/blockcypher/libgrin/core"
	"github.com/blockcypher/libgrin/keychain"
	"github.com/google/uuid"
)

// OutputData is the information about an output that's being tracked by the wallet. Must be
// enough to reconstruct the commitment associated with the output when the
// root private key is known.
type OutputData struct {
	// Root key_id that the key for this output is derived from
	RootKeyID keychain.Identifier `json:"root_key_id"`
	// Derived key for this output
	KeyID keychain.Identifier `json:"key_id"`
	// How many derivations down from the root key
	NChild uint32 `json:"n_child"`
	// The actual commit optionally stored
	Commit *string `json:"commit"`
	// PMMR Index, used on restore in case of duplicate wallets using the same
	// key_id (2 wallets using same seed, for instance
	MMRIndex *core.Uint64 `json:"mmr_index"`
	// Value of the output, necessary to rebuild the commitment
	Value core.Uint64 `json:"value"`
	// Current status of the output
	Status OutputStatus `json:"status"`
	// Height of the output
	Height core.Uint64 `json:"height"`
	// Height we are locked until
	LockHeight core.Uint64 `json:"lock_height"`
	// Is this a coinbase output? Is it subject to coinbase locktime?
	IsCoinbase bool `json:"is_coinbase"`
	// Optional corresponding internal entry in tx entry log
	TxLogEntry *uint32 `json:"tx_log_entry"`
}

// OutputStatus is the status of an output that's being tracked by the wallet.
// Can either be unconfirmed, spent, unspent, or locked (when it's been used
//to generate a transaction but we don't have confirmation that the transaction
// was broadcasted or mined).
type OutputStatus int

const (
	// Unconfirmed output
	Unconfirmed OutputStatus = iota
	// Unspent output
	Unspent
	// Locked output
	Locked
	// Spent output
	Spent
)

func (s OutputStatus) String() string {
	return toStringOutputStatus[s]
}

var toStringOutputStatus = map[OutputStatus]string{
	Unconfirmed: "Unconfirmed",
	Unspent:     "Unspent",
	Locked:      "Locked",
	Spent:       "Spent",
}

var toIDOutputStatus = map[string]OutputStatus{
	"Unconfirmed": Unconfirmed,
	"Unspent":     Unspent,
	"Locked":      Locked,
	"Spent":       Spent,
}

// MarshalJSON marshals the enum as a quoted json string
func (s OutputStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringOutputStatus[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *OutputStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toIDOutputStatus[j]
	return nil
}

// WalletInfo is a contained wallet info struct, so automated tests can parse
// wallet info can add more fields here over time as needed
type WalletInfo struct {
	// height from which info was taken
	LastConfirmedHeight core.Uint64 `json:"last_confirmed_height"`
	// Minimum number of confirmations for an output to be treated as "spendable".
	MinimumConfirmations core.Uint64 `json:"minimum_confirmations"`
	// total amount in the wallet
	Total core.Uint64 `json:"total"`
	// amount awaiting finalization
	AmountAwaitingFinalization core.Uint64 `json:"amount_awaiting_finalization"`
	// amount awaiting confirmation
	AmountAwaitingConfirmation core.Uint64 `json:"amount_awaiting_confirmation"`
	// coinbases waiting for lock height
	AmountImmature core.Uint64 `json:"amount_immature"`
	// amount currently spendable
	AmountCurrentlySpendable core.Uint64 `json:"amount_currently_spendable"`
	// amount locked via previous transactions
	AmountLocked core.Uint64 `json:"amount_locked"`
}

// TxLogEntryType represent the type of transactions that can be contained
// within a TXLog entry
type TxLogEntryType int

const (
	// ConfirmedCoinbase is a coinbase transaction becomes confirmed
	ConfirmedCoinbase TxLogEntryType = iota
	// TxReceived are outputs created when a transaction is received
	TxReceived
	// TxSent are inputs locked + change outputs when a transaction is created
	TxSent
	// TxReceivedCancelled is a received transaction that was rolled back by user
	TxReceivedCancelled
	// TxSentCancelled is a sent transaction that was rolled back by user
	TxSentCancelled
)

func (s TxLogEntryType) String() string {
	return toStringTxLogEntryType[s]
}

var toStringTxLogEntryType = map[TxLogEntryType]string{
	ConfirmedCoinbase:   "ConfirmedCoinbase",
	TxReceived:          "TxReceived",
	TxSent:              "TxSent",
	TxReceivedCancelled: "TxReceivedCancelled",
	TxSentCancelled:     "TxSentCancelled",
}

var toIDTxLogEntryType = map[string]TxLogEntryType{
	"ConfirmedCoinbase":   ConfirmedCoinbase,
	"TxReceived":          TxReceived,
	"TxSent":              TxSent,
	"TxReceivedCancelled": TxReceivedCancelled,
	"TxSentCancelled":     TxSentCancelled,
}

// MarshalJSON marshals the enum as a quoted json string
func (s TxLogEntryType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringTxLogEntryType[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *TxLogEntryType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toIDTxLogEntryType[j]
	return nil
}

// ParticipantMessages is an helper just to facilitate serialization
type ParticipantMessages struct {
	// included messages
	Messages []ParticipantMessageData `json:"messages"`
}

// TxLogEntry is an optional transaction information, recorded when an event
// happens to add or remove funds from a wallet. One Transaction log entry maps
// to one or many outputs
type TxLogEntry struct {
	// BIP32 account path used for creating this tx
	ParentKeyID keychain.Identifier `json:"parent_key_id"`
	// Local id for this transaction (distinct from a slate transaction id)
	ID uint32 `json:"id"`
	// Slate transaction this entry is associated with, if any
	TxSlateID *uuid.UUID `json:"tx_slate_id"`
	// Transaction type (as above)
	TxType TxLogEntryType `json:"tx_type"`
	// Time this tx entry was created
	// #[serde(with = "tx_date_format")]
	CreationTs time.Time `json:"creation_ts"`
	// Time this tx was confirmed (by this wallet)
	// #[serde(default, with = "opt_tx_date_format")]
	ConfirmationTs *time.Time `json:"confirmation_ts"`
	// Whether the inputs+outputs involved in this transaction have been
	// confirmed (In all cases either all outputs involved in a tx should be
	// confirmed, or none should be; otherwise there's a deeper problem)
	Confirmed bool `json:"confirmed"`
	// number of inputs involved in TX
	NumInputs uint `json:"num_inputs"`
	// number of outputs involved in TX
	NumOutputs uint `json:"num_outputs"`
	// Amount credited via this transaction
	AmountCredited core.Uint64 `json:"amount_credited"`
	// Amount debited via this transaction
	AmountDebited core.Uint64 `json:"amount_debited"`
	// Fee
	Fee *core.Uint64 `json:"fee"`
	// Message data, stored as json
	Messages *ParticipantMessages `json:"messages"`
	// Location of the store transaction, (reference or resending)
	StoredTx *string `json:"stored_tx"`
}
