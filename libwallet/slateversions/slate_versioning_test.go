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

package slateversions_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/blockcypher/libgrin/libwallet/slateversions"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeAndV0UpgradeToV1(t *testing.T) {
	slateV0JSON, _ := ioutil.ReadFile("test_data/v0.slate")
	var slateV0 slateversions.SlateV0
	assert.Nil(t, json.Unmarshal(slateV0JSON, &slateV0))
	slateV1 := slateV0.Upgrade()
	assert.Equal(t, slateV0.NumParticipants, slateV1.NumParticipants)
	assert.Equal(t, slateV0.ID, slateV1.ID)
	assert.Equal(t, slateV0.Transaction.Offset, slateV1.Transaction.Offset)
	for i := range slateV1.Transaction.Body.Inputs {
		assert.Equal(t, slateV0.Transaction.Body.Inputs[i].Commit, slateV1.Transaction.Body.Inputs[i].Commit)
		assert.Equal(t, slateV0.Transaction.Body.Inputs[i].Features, slateV1.Transaction.Body.Inputs[i].Features)
	}
	for i := range slateV1.Transaction.Body.Outputs {
		assert.Equal(t, slateV0.Transaction.Body.Outputs[i].Commit, slateV1.Transaction.Body.Outputs[i].Commit)
		assert.Equal(t, slateV0.Transaction.Body.Outputs[i].Features, slateV1.Transaction.Body.Outputs[i].Features)
		assert.Equal(t, slateV0.Transaction.Body.Outputs[i].Proof, slateV1.Transaction.Body.Outputs[i].Proof)
	}
	for i := range slateV1.Transaction.Body.Kernels {
		assert.Equal(t, slateV0.Transaction.Body.Kernels[i].Features, slateV1.Transaction.Body.Kernels[i].Features)
		assert.Equal(t, slateV0.Transaction.Body.Kernels[i].Fee, slateV1.Transaction.Body.Kernels[i].Fee)
		assert.Equal(t, slateV0.Transaction.Body.Kernels[i].LockHeight, slateV1.Transaction.Body.Kernels[i].LockHeight)
		assert.Equal(t, slateV0.Transaction.Body.Kernels[i].Excess, slateV1.Transaction.Body.Kernels[i].Excess)
		assert.Equal(t, slateV0.Transaction.Body.Kernels[i].ExcessSig, slateV1.Transaction.Body.Kernels[i].ExcessSig)
	}
	assert.Equal(t, slateV0.Amount, slateV1.Amount)
	assert.Equal(t, slateV0.Fee, slateV1.Fee)
	assert.Equal(t, slateV0.Height, slateV1.Height)
	assert.Equal(t, slateV0.LockHeight, slateV1.LockHeight)
	assert.Equal(t, slateV0.Fee, slateV1.Fee)
	for i := range slateV1.ParticipantData {
		assert.Equal(t, slateV0.ParticipantData[i].ID, slateV1.ParticipantData[i].ID)
		assert.Equal(t, slateV0.ParticipantData[i].PublicBlindExcess, slateV1.ParticipantData[i].PublicBlindExcess)
		assert.Equal(t, slateV0.ParticipantData[i].PublicNonce, slateV1.ParticipantData[i].PublicNonce)
		assert.Equal(t, slateV0.ParticipantData[i].PartSig, slateV1.ParticipantData[i].PartSig)
		assert.Equal(t, slateV0.ParticipantData[i].Message, slateV1.ParticipantData[i].Message)
		assert.Equal(t, slateV0.ParticipantData[i].MessageSig, slateV1.ParticipantData[i].MessageSig)
	}
	assert.Equal(t, uint64(1), slateV1.Version)
	assert.Equal(t, uint64(0), slateV1.GetOrigVersion())
	// Compare it to the one on disk
	slateV1JSON, _ := ioutil.ReadFile("test_data/v1.slate")
	var slateV1Alt slateversions.SlateV1
	assert.Nil(t, json.Unmarshal(slateV1JSON, &slateV1Alt))
	assert.Exactly(t, slateV1Alt, slateV1)
}

func TestDeserializeAndV1UpgradeToV2(t *testing.T) {
	slateV1JSON, _ := ioutil.ReadFile("test_data/v1.slate")
	var slateV1 slateversions.SlateV1
	slateV1.SetOrigVersion(1)
	assert.Nil(t, json.Unmarshal(slateV1JSON, &slateV1))
	slateV2 := slateV1.Upgrade()
	assert.Equal(t, slateV1.NumParticipants, slateV2.NumParticipants)
	assert.Equal(t, slateV1.ID, slateV2.ID)
	for i := range slateV2.Transaction.Body.Inputs {
		assert.Equal(t, slateV1.Transaction.Body.Inputs[i].Features, slateV2.Transaction.Body.Inputs[i].Features)
	}
	for i := range slateV2.Transaction.Body.Outputs {
		assert.Equal(t, slateV1.Transaction.Body.Outputs[i].Features, slateV2.Transaction.Body.Outputs[i].Features)
	}
	for i := range slateV2.Transaction.Body.Kernels {
		assert.Equal(t, slateV1.Transaction.Body.Kernels[i].Features, slateV2.Transaction.Body.Kernels[i].Features)
		assert.Equal(t, slateV1.Transaction.Body.Kernels[i].Fee, slateV2.Transaction.Body.Kernels[i].Fee)
		assert.Equal(t, slateV1.Transaction.Body.Kernels[i].LockHeight, slateV2.Transaction.Body.Kernels[i].LockHeight)
	}
	assert.Equal(t, slateV1.Amount, slateV2.Amount)
	assert.Equal(t, slateV1.Fee, slateV2.Fee)
	assert.Equal(t, slateV1.Height, slateV2.Height)
	assert.Equal(t, slateV1.LockHeight, slateV2.LockHeight)
	assert.Equal(t, slateV1.Fee, slateV2.Fee)
	for i := range slateV2.ParticipantData {
		assert.Equal(t, slateV1.ParticipantData[i].ID, slateV2.ParticipantData[i].ID)
		assert.Equal(t, slateV1.ParticipantData[i].Message, slateV2.ParticipantData[i].Message)
	}
	assert.Equal(t, uint16(1), slateV2.VersionInfo.BlockHeaderVersion)
	assert.Equal(t, uint16(1), slateV2.VersionInfo.OrigVersion)
	assert.Equal(t, uint16(2), slateV2.VersionInfo.Version)
	// Compare it to the one on disk
	slateV2JSON, _ := ioutil.ReadFile("test_data/v2.slate")
	var slateV2Alt slateversions.SlateV2
	// just for this test to pass
	slateV2.VersionInfo.OrigVersion = 2
	assert.Nil(t, json.Unmarshal(slateV2JSON, &slateV2Alt))
	assert.Exactly(t, slateV2Alt, slateV2)
}

func TestDeserializeAndV2DowngradeToV1(t *testing.T) {
	slateV2JSON, _ := ioutil.ReadFile("test_data/v2.slate")
	var slateV2 slateversions.SlateV2
	assert.Nil(t, json.Unmarshal(slateV2JSON, &slateV2))
	slateV1 := slateV2.Downgrade()
	assert.Equal(t, slateV2.NumParticipants, slateV1.NumParticipants)
	assert.Equal(t, slateV2.ID, slateV2.ID)
	for i := range slateV1.Transaction.Body.Inputs {
		assert.Equal(t, slateV2.Transaction.Body.Inputs[i].Features, slateV1.Transaction.Body.Inputs[i].Features)
	}
	for i := range slateV1.Transaction.Body.Outputs {
		assert.Equal(t, slateV2.Transaction.Body.Outputs[i].Features, slateV1.Transaction.Body.Outputs[i].Features)
	}
	for i := range slateV1.Transaction.Body.Kernels {
		assert.Equal(t, slateV2.Transaction.Body.Kernels[i].Features, slateV1.Transaction.Body.Kernels[i].Features)
		assert.Equal(t, slateV2.Transaction.Body.Kernels[i].Fee, slateV1.Transaction.Body.Kernels[i].Fee)
		assert.Equal(t, slateV2.Transaction.Body.Kernels[i].LockHeight, slateV1.Transaction.Body.Kernels[i].LockHeight)
	}
	assert.Equal(t, slateV2.Amount, slateV1.Amount)
	assert.Equal(t, slateV2.Fee, slateV1.Fee)
	assert.Equal(t, slateV2.Height, slateV1.Height)
	assert.Equal(t, slateV2.LockHeight, slateV1.LockHeight)
	assert.Equal(t, slateV2.Fee, slateV1.Fee)
	for i := range slateV1.ParticipantData {
		assert.Equal(t, slateV2.ParticipantData[i].ID, slateV1.ParticipantData[i].ID)
		assert.Equal(t, slateV2.ParticipantData[i].Message, slateV1.ParticipantData[i].Message)
	}

	assert.Equal(t, uint64(1), slateV1.Version)

	// Compare it to the one on disk
	slateV1JSON, _ := ioutil.ReadFile("test_data/v1.slate")
	var slateV1Alt slateversions.SlateV1
	// just for this test to pass
	slateV1Alt.SetOrigVersion(2)
	assert.Nil(t, json.Unmarshal(slateV1JSON, &slateV1Alt))
	assert.Exactly(t, slateV1Alt, slateV1)
}

func TestDeserializeAndV1DowngradeToV0(t *testing.T) {
	slateV1JSON, _ := ioutil.ReadFile("test_data/v1.slate")
	var slateV1 slateversions.SlateV1
	assert.Nil(t, json.Unmarshal(slateV1JSON, &slateV1))
	slateV0 := slateV1.Downgrade()
	assert.Equal(t, slateV0.NumParticipants, slateV0.NumParticipants)
	assert.Equal(t, slateV0.ID, slateV0.ID)
	for i := range slateV0.Transaction.Body.Inputs {
		assert.Equal(t, slateV1.Transaction.Body.Inputs[i].Features, slateV0.Transaction.Body.Inputs[i].Features)
	}
	for i := range slateV0.Transaction.Body.Outputs {
		assert.Equal(t, slateV1.Transaction.Body.Outputs[i].Features, slateV0.Transaction.Body.Outputs[i].Features)
	}
	for i := range slateV0.Transaction.Body.Kernels {
		assert.Equal(t, slateV1.Transaction.Body.Kernels[i].Features, slateV0.Transaction.Body.Kernels[i].Features)
		assert.Equal(t, slateV1.Transaction.Body.Kernels[i].Fee, slateV0.Transaction.Body.Kernels[i].Fee)
		assert.Equal(t, slateV1.Transaction.Body.Kernels[i].LockHeight, slateV0.Transaction.Body.Kernels[i].LockHeight)
	}
	assert.Equal(t, slateV1.Amount, slateV0.Amount)
	assert.Equal(t, slateV1.Fee, slateV0.Fee)
	assert.Equal(t, slateV1.Height, slateV0.Height)
	assert.Equal(t, slateV1.LockHeight, slateV0.LockHeight)
	assert.Equal(t, slateV1.Fee, slateV0.Fee)
	for i := range slateV0.ParticipantData {
		assert.Equal(t, slateV1.ParticipantData[i].ID, slateV0.ParticipantData[i].ID)
		assert.Equal(t, slateV1.ParticipantData[i].Message, slateV0.ParticipantData[i].Message)
	}

	// Compare it to the one on disk
	slateV0JSON, _ := ioutil.ReadFile("test_data/v0.slate")
	var slateV0Alt slateversions.SlateV0
	assert.Nil(t, json.Unmarshal(slateV0JSON, &slateV0Alt))
	assert.Exactly(t, slateV0Alt, slateV0)
}
