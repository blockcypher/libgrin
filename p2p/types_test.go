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

func TestUnmarshalReasonForBan(t *testing.T) {
	noneb := []byte(`"None"`)
	var none reasonForBan
	if err := json.Unmarshal(noneb, &none); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, none, NoneReasonForBan)

	badblockb := []byte(`"BadBlock"`)
	var badblock reasonForBan
	if err := json.Unmarshal(badblockb, &badblock); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, badblock, BadBlockReasonForBan)

	badcompactblockb := []byte(`"BadCompactBlock"`)
	var badcompactblock reasonForBan
	if err := json.Unmarshal(badcompactblockb, &badcompactblock); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, badcompactblock, BadCompactBlockReasonForBan)

	badblockheaderb := []byte(`"BadBlockHeader"`)
	var badblockheader reasonForBan
	if err := json.Unmarshal(badblockheaderb, &badblockheader); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, badblockheader, BadBlockHeaderReasonForBan)

	badtxhashsetb := []byte(`"BadTxHashSet"`)
	var badtxhashset reasonForBan
	if err := json.Unmarshal(badtxhashsetb, &badtxhashset); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, badtxhashset, BadTxHashSetReasonForBan)

	manualbanb := []byte(`"ManualBan"`)
	var manualban reasonForBan
	if err := json.Unmarshal(manualbanb, &manualban); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, manualban, ManualBanReasonForBan)

	fraudheightb := []byte(`"FraudHeight"`)
	var fraudheight reasonForBan
	if err := json.Unmarshal(fraudheightb, &fraudheight); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, fraudheight, FraudHeightReasonForBan)

	badhanshakeb := []byte(`"BadHandshake"`)
	var badhanshake reasonForBan
	if err := json.Unmarshal(badhanshakeb, &badhanshake); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, badhanshake, BadHanshakeReasonForBan)

}

func TestMarshalReasonForBan(t *testing.T) {
	noneb, err := json.Marshal(NoneReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(noneb), "\"None\"")

	badblockb, err := json.Marshal(BadBlockReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(badblockb), "\"BadBlock\"")

	badcompactblockb, err := json.Marshal(BadCompactBlockReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(badcompactblockb), "\"BadCompactBlock\"")

	badblockheaderb, err := json.Marshal(BadBlockHeaderReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(badblockheaderb), "\"BadBlockHeader\"")

	badtxhashsetb, err := json.Marshal(BadTxHashSetReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(badtxhashsetb), "\"BadTxHashSet\"")

	manualbanb, err := json.Marshal(ManualBanReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(manualbanb), "\"ManualBan\"")

	fraudheightb, err := json.Marshal(FraudHeightReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(fraudheightb), "\"FraudHeight\"")

	badhanshakeb, err := json.Marshal(BadHanshakeReasonForBan)
	assert.Nil(t, err)
	assert.Equal(t, string(badhanshakeb), "\"BadHandshake\"")
}

func TestUnmarshalConnectionDirection(t *testing.T) {
	inboundb := []byte(`"Inbound"`)
	var inbound connectionDirection
	if err := json.Unmarshal(inboundb, &inbound); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, inbound, InboundConnectionDirection)

	outboundb := []byte(`"Outbound"`)
	var outbound connectionDirection
	if err := json.Unmarshal(outboundb, &outbound); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, outbound, OutboundConnectionDirection)
}

func TestMarshalConnectionDirection(t *testing.T) {
	inboundb, err := json.Marshal(InboundConnectionDirection)
	assert.Nil(t, err)
	assert.Equal(t, string(inboundb), "\"Inbound\"")

	outboundb, err := json.Marshal(OutboundConnectionDirection)
	assert.Nil(t, err)
	assert.Equal(t, string(outboundb), "\"Outbound\"")
}
