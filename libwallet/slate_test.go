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

package libwallet_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/blockcypher/libgrin/libwallet"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalUpgradeV2(t *testing.T) {
	slateV2JSON, _ := ioutil.ReadFile("slateversions/test_data/v2.slate")
	var slateV2 libwallet.Slate
	err := libwallet.UnmarshalUpgrade(slateV2JSON, &slateV2)
	assert.Nil(t, err)

	// Compare with a direct unmarshal
	slateV2JSONReference, _ := ioutil.ReadFile("slateversions/test_data/v2.slate")
	var slateV2Reference libwallet.Slate
	assert.Nil(t, json.Unmarshal(slateV2JSONReference, &slateV2Reference))
	assert.Exactly(t, slateV2Reference, slateV2)

}

func TestUnmarshalUpgradeV1(t *testing.T) {
	slateV1JSON, _ := ioutil.ReadFile("slateversions/test_data/v1.slate")
	var slateV2 libwallet.Slate
	err := libwallet.UnmarshalUpgrade(slateV1JSON, &slateV2)
	assert.Equal(t, uint16(1), slateV2.VersionInfo.OrigVersion)
	assert.Nil(t, err)

	// Compare with a direct unmarshal from a slate V2
	slateV2JSONReference, _ := ioutil.ReadFile("slateversions/test_data/v2.slate")
	var slateV2Reference libwallet.Slate
	assert.Nil(t, json.Unmarshal(slateV2JSONReference, &slateV2Reference))
	// just for this test to pass
	slateV2.VersionInfo.OrigVersion = 2
	assert.Exactly(t, slateV2Reference, slateV2)
}

func TestUnmarshalUpgradeV0(t *testing.T) {
	slateV0JSON, _ := ioutil.ReadFile("slateversions/test_data/v0.slate")
	var slateV2 libwallet.Slate
	err := libwallet.UnmarshalUpgrade(slateV0JSON, &slateV2)
	assert.Nil(t, err)

	// Compare with a direct unmarshal from a slate V2
	slateV2JSONReference, _ := ioutil.ReadFile("slateversions/test_data/v2.slate")
	var slateV2Reference libwallet.Slate
	assert.Nil(t, json.Unmarshal(slateV2JSONReference, &slateV2Reference))
	// just for this test to pass
	slateV2.VersionInfo.OrigVersion = 2
	assert.Exactly(t, slateV2Reference, slateV2)
}

func TestMarshal(t *testing.T) {
	slateV2JSON, _ := ioutil.ReadFile("slateversions/test_data/v2_raw.slate")
	var slateV2 libwallet.Slate
	err := libwallet.UnmarshalUpgrade(slateV2JSON, &slateV2)
	assert.Nil(t, err)

	// First test that if we put nothing it serialize as a Slate V2
	serializedSlateV2, err := libwallet.Marshal(slateV2)
	assert.Nil(t, err)
	assert.Equal(t, slateV2JSON, serializedSlateV2)

	// Then test that it does serialize correctly with V1
	slateV2.VersionInfo.OrigVersion = 1
	serializedSlateV1, err := libwallet.Marshal(slateV2)
	assert.Nil(t, err)
	slateV1Ref, _ := ioutil.ReadFile("slateversions/test_data/v1_raw.slate")
	assert.Equal(t, slateV1Ref, serializedSlateV1)

	// Then test that it does serialize correctly with V1
	slateV2.VersionInfo.OrigVersion = 0
	serializedSlateV0, err := libwallet.Marshal(slateV2)
	assert.Nil(t, err)
	slateV0Ref, _ := ioutil.ReadFile("slateversions/test_data/v0_raw.slate")
	assert.Equal(t, slateV0Ref, serializedSlateV0)
}
