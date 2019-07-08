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

func TestMarshal(t *testing.T) {
	slateV2JSON, _ := ioutil.ReadFile("slateversions/test_data/v2_raw.slate")
	var slateV2 libwallet.Slate
	err := libwallet.UnmarshalUpgrade(slateV2JSON, &slateV2)
	assert.Nil(t, err)

	// First test that if we put nothing it serialize as a Slate V2
	serializedSlateV2, err := json.Marshal(slateV2)
	assert.Nil(t, err)
	assert.Equal(t, slateV2JSON, serializedSlateV2)
}
