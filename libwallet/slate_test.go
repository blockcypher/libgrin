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

package libwallet_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/blockcypher/libgrin/libwallet"
	"github.com/blockcypher/libgrin/libwallet/slateversions"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalV2(t *testing.T) {
	slateV2JSON, _ := ioutil.ReadFile("slateversions/test_data/v2.slate")
	var slateV2 slateversions.SlateV2
	err := json.Unmarshal(slateV2JSON, &slateV2)
	assert.Nil(t, err)
}

func TestMarshalV2(t *testing.T) {
	slateV2JSON, _ := ioutil.ReadFile("slateversions/test_data/v2_raw.slate")
	var slateV2 slateversions.SlateV2
	err := json.Unmarshal(slateV2JSON, &slateV2)
	assert.Nil(t, err)

	serializedSlateV2, err := json.Marshal(slateV2)
	assert.Equal(t, slateV2JSON, serializedSlateV2)
}

func TestUnmarshalV3(t *testing.T) {
	slateV3JSON, _ := ioutil.ReadFile("slateversions/test_data/v3.slate")
	var slateV3 libwallet.Slate
	err := json.Unmarshal(slateV3JSON, &slateV3)
	assert.Nil(t, err)
}

func TestMarshalV3(t *testing.T) {
	slateV3JSON, _ := ioutil.ReadFile("slateversions/test_data/v3_raw.slate")
	var slateV3 libwallet.Slate
	err := json.Unmarshal(slateV3JSON, &slateV3)
	assert.Nil(t, err)

	serializedSlateV3, err := json.Marshal(slateV3)
	assert.Equal(t, slateV3JSON, serializedSlateV3)

}
