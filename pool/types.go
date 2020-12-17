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

package pool

import (
	"bytes"
	"encoding/json"

	"github.com/blockcypher/libgrin/v5/core"
)

// TxSource is used to make decisions based on transaction acceptance priority from
// various sources. For example, a node may want to bypass pool size
// restrictions when accepting a transaction from a local wallet.
//
// Most likely this will evolve to contain some sort of network identifier,
// once we get a better sense of what transaction building might look like.
type TxSource int

const (
	// PushAPITxSource when the tx is pushed via API
	PushAPITxSource TxSource = iota
	// BroadcastTxSource when the source is a broadcats
	BroadcastTxSource
	// FluffTxSource zhen the source is our own fluff
	FluffTxSource
	// EmbargoExpiredTxSource when the source is an expired embargo
	EmbargoExpiredTxSource
	// DeaggregateTxSource when the source is the deaggregation
	DeaggregateTxSource
)

var toStringTxSource = map[TxSource]string{
	PushAPITxSource:        "PushApi",
	BroadcastTxSource:      "Broadcast",
	FluffTxSource:          "Fluff",
	EmbargoExpiredTxSource: "EmbargoExpired",
	DeaggregateTxSource:    "Deaggregate",
}

var toIDTxSource = map[string]TxSource{
	"PushApi":        PushAPITxSource,
	"Broadcast":      BroadcastTxSource,
	"Fluff":          FluffTxSource,
	"EmbargoExpired": EmbargoExpiredTxSource,
	"Deaggregate":    DeaggregateTxSource,
}

// MarshalJSON marshals the enum as a quoted json string
func (s TxSource) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringTxSource[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *TxSource) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'NoneReasonForBan' in this case.
	*s = toIDTxSource[j]
	return nil
}

// PoolEntry represents a single entry in the pool.
// A single (possibly aggregated) transaction.
type PoolEntry struct {
	// Info on where this tx originated from.
	Src TxSource `json:"src"`
	// Timestamp of when this tx was originally added to the pool.
	TxAt string `json:"tx_at"`
	// The transaction itself.
	Tx core.Transaction `json:"tx"`
}
