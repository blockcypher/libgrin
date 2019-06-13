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
