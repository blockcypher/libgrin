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

package p2p

import (
	"bytes"
	"encoding/json"
	"net"

	"github.com/blockcypher/libgrin/v5/core/consensus"
)

// Capabilities are options for what type of interaction a peer supports
type Capabilities struct {
	Bits uint32 `json:"bits"`
}

// reasonForBan is the ban reason
type reasonForBan int

const (
	// NoneReasonForBan not banned
	NoneReasonForBan reasonForBan = iota
	// BadBlockReasonForBan banned for sending a bad block
	BadBlockReasonForBan
	// BadCompactBlockReasonForBan banned for sending a bad compact block
	BadCompactBlockReasonForBan
	// BadBlockHeaderReasonForBan banned for sending a bad block header
	BadBlockHeaderReasonForBan
	// BadTxHashSetReasonForBan banned for sending a bad txhashset
	BadTxHashSetReasonForBan
	// ManualBanReasonForBan banned manualy
	ManualBanReasonForBan
	// FraudHeightReasonForBan banned for fraud height
	FraudHeightReasonForBan
	// BadHanshakeReasonForBan banned for sending a bad handshake
	BadHanshakeReasonForBan
)

var toStringReasonForBan = map[reasonForBan]string{
	NoneReasonForBan:            "None",
	BadBlockReasonForBan:        "BadBlock",
	BadCompactBlockReasonForBan: "BadCompactBlock",
	BadBlockHeaderReasonForBan:  "BadBlockHeader",
	BadTxHashSetReasonForBan:    "BadTxHashSet",
	ManualBanReasonForBan:       "ManualBan",
	FraudHeightReasonForBan:     "FraudHeight",
	BadHanshakeReasonForBan:     "BadHandshake",
}

var toIDReasonForBan = map[string]reasonForBan{
	"None":            NoneReasonForBan,
	"BadBlock":        BadBlockReasonForBan,
	"BadCompactBlock": BadCompactBlockReasonForBan,
	"BadBlockHeader":  BadBlockHeaderReasonForBan,
	"BadTxHashSet":    BadTxHashSetReasonForBan,
	"ManualBan":       ManualBanReasonForBan,
	"FraudHeight":     FraudHeightReasonForBan,
	"BadHandshake":    BadHanshakeReasonForBan,
}

// MarshalJSON marshals the enum as a quoted json string
func (s reasonForBan) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringReasonForBan[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *reasonForBan) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'NoneReasonForBan' in this case.
	*s = toIDReasonForBan[j]
	return nil
}

// PeerInfoDisplay is a flatten out PeerInfo and nested PeerLiveInfo (taking a read lock on it)
// so we can serialize/deserialize the data for the API and the TUI.
type PeerInfoDisplay struct {
	Capabilities    Capabilities         `json:"capabilities"`
	UserAgent       string               `json:"user_agent"`
	Version         uint32               `json:"version"`
	Addr            net.IP               `json:"addr"`
	Direction       connectionDirection  `json:"direction"`
	TotalDifficulty consensus.Difficulty `json:"total_difficulty"`
	Height          uint64               `json:"height"`
}

// Types of connection direction
type connectionDirection int

const (
	// InboundConnectionDirection is an inbound connection
	InboundConnectionDirection connectionDirection = iota
	// OutboundConnectionDirection is an outbound connection
	OutboundConnectionDirection
)

var toStringConnection = map[connectionDirection]string{
	InboundConnectionDirection:  "Inbound",
	OutboundConnectionDirection: "Outbound",
}

var toIDConnection = map[string]connectionDirection{
	"Inbound":  InboundConnectionDirection,
	"Outbound": OutboundConnectionDirection,
}

// MarshalJSON marshals the enum as a quoted json string
func (s connectionDirection) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringConnection[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *connectionDirection) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'InboundConnection' in this case.
	*s = toIDConnection[j]
	return nil
}
