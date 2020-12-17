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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalTxSource(t *testing.T) {
	pushapib := []byte(`"PushApi"`)
	var pushapi TxSource
	if err := json.Unmarshal(pushapib, &pushapi); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, pushapi, PushAPITxSource)

	broadcastb := []byte(`"Broadcast"`)
	var broadcast TxSource
	if err := json.Unmarshal(broadcastb, &broadcast); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, broadcast, BroadcastTxSource)

	fluffb := []byte(`"Fluff"`)
	var fluff TxSource
	if err := json.Unmarshal(fluffb, &fluff); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, fluff, FluffTxSource)

	eeb := []byte(`"EmbargoExpired"`)
	var ee TxSource
	if err := json.Unmarshal(eeb, &ee); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, ee, EmbargoExpiredTxSource)

	deaggregateb := []byte(`"Deaggregate"`)
	var deaggregate TxSource
	if err := json.Unmarshal(deaggregateb, &deaggregate); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, deaggregate, DeaggregateTxSource)
}

func TestMarshalTxSource(t *testing.T) {
	pushapib, err := json.Marshal(PushAPITxSource)
	assert.Nil(t, err)
	assert.Equal(t, string(pushapib), "\"PushApi\"")

	broadcastb, err := json.Marshal(BroadcastTxSource)
	assert.Nil(t, err)
	assert.Equal(t, string(broadcastb), "\"Broadcast\"")

	fluffb, err := json.Marshal(FluffTxSource)
	assert.Nil(t, err)
	assert.Equal(t, string(fluffb), "\"Fluff\"")

	eeb, err := json.Marshal(EmbargoExpiredTxSource)
	assert.Nil(t, err)
	assert.Equal(t, string(eeb), "\"EmbargoExpired\"")

	deaggregateb, err := json.Marshal(DeaggregateTxSource)
	assert.Nil(t, err)
	assert.Equal(t, string(deaggregateb), "\"Deaggregate\"")

}
