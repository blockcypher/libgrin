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
)

// PeerData is the data stored for any given peer we've encountered.
type PeerData struct {
	// Network address of the peer.
	Addr net.IP `json:"addr"`
	// What capabilities the peer advertises. Unknown until a successful
	// connection.
	Capabilities Capabilities `json:"capabilities"`
	// The peer user agent.
	UserAgent string `json:"user_agent"`
	// State the peer has been detected with.
	Flags peerState `json:"flags"`
	// The time the peer was last banned
	LastBanned int64 `json:"last_banned"`
	// The reason for the ban
	Banreason reasonForBan `json:"ban_reason"`
	// Time when we last connected to this peer.
	LastConnected int64 `json:"1570129317"`
}

// PeerState is the state of a peer
type peerState int

const (
	// HealthyPeerState a healthy peer
	HealthyPeerState peerState = iota
	// BannedPeerState a banned peer
	BannedPeerState
	// DefunctPeerState a banned peer
	DefunctPeerState
)

var toStringPeerState = map[peerState]string{
	HealthyPeerState: "Healthy",
	BannedPeerState:  "Banned",
	DefunctPeerState: "Defunct",
}

var toIDPeerState = map[string]peerState{
	"Healthy": HealthyPeerState,
	"Banned":  BannedPeerState,
	"Defunct": DefunctPeerState,
}

// MarshalJSON marshals the enum as a quoted json string
func (s peerState) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringPeerState[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *peerState) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'HealthyPeerState' in this case.
	*s = toIDPeerState[j]
	return nil
}
