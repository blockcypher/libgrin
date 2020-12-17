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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalPeerState(t *testing.T) {
	healthyb := []byte(`"Healthy"`)
	var healthy peerState
	if err := json.Unmarshal(healthyb, &healthy); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, healthy, HealthyPeerState)

	bannedb := []byte(`"Banned"`)
	var banned peerState
	if err := json.Unmarshal(bannedb, &banned); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, banned, BannedPeerState)

	defunctb := []byte(`"Defunct"`)
	var defunct peerState
	if err := json.Unmarshal(defunctb, &defunct); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, defunct, DefunctPeerState)
}

func TestMarshalPeerState(t *testing.T) {
	healthyb, err := json.Marshal(HealthyPeerState)
	assert.Nil(t, err)
	assert.Equal(t, string(healthyb), "\"Healthy\"")

	bannedb, err := json.Marshal(BannedPeerState)
	assert.Nil(t, err)
	assert.Equal(t, string(bannedb), "\"Banned\"")

	defunctb, err := json.Marshal(DefunctPeerState)
	assert.Nil(t, err)
	assert.Equal(t, string(defunctb), "\"Defunct\"")
}
