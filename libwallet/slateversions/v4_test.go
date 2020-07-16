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

package slateversions

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSlateState(t *testing.T) {
	s1b := []byte(`"S1"`)
	var s1 slateStateV4
	if err := json.Unmarshal(s1b, &s1); err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, s1, Standard1SlateState)

	s2b := []byte(`"S2"`)
	var s2 slateStateV4
	if err := json.Unmarshal(s2b, &s2); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, s2, Standard2SlateState)

	s3b := []byte(`"S3"`)
	var s3 slateStateV4
	if err := json.Unmarshal(s3b, &s3); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, s3, Standard3SlateState)

	i1b := []byte(`"I1"`)
	var i1 slateStateV4
	if err := json.Unmarshal(i1b, &i1); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, i1, Invoice1SlateState)

	i2b := []byte(`"I2"`)
	var i2 slateStateV4
	if err := json.Unmarshal(i2b, &i2); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, i2, Invoice2SlateState)

	i3b := []byte(`"I3"`)
	var i3 slateStateV4
	if err := json.Unmarshal(i3b, &i3); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, i3, Invoice3SlateState)
}

func TestMarshalSlateState(t *testing.T) {
	s1b, err := json.Marshal(Standard1SlateState)
	assert.Nil(t, err)
	assert.Equal(t, string(s1b), "\"S1\"")

	s2b, err := json.Marshal(Standard2SlateState)
	assert.Nil(t, err)
	assert.Equal(t, string(s2b), "\"S2\"")

	s3b, err := json.Marshal(Standard3SlateState)
	assert.Nil(t, err)
	assert.Equal(t, string(s3b), "\"S3\"")

	i1b, err := json.Marshal(Invoice1SlateState)
	assert.Nil(t, err)
	assert.Equal(t, string(i1b), "\"I1\"")

	i2b, err := json.Marshal(Invoice2SlateState)
	assert.Nil(t, err)
	assert.Equal(t, string(i2b), "\"I2\"")

	i3b, err := json.Marshal(Invoice3SlateState)
	assert.Nil(t, err)
	assert.Equal(t, string(i3b), "\"I3\"")
}

func TestUnmarshalSlateV4(t *testing.T) {
	slateV4JSON, _ := ioutil.ReadFile("test_data/v4.slate")
	var slateV4 SlateV4
	if err := json.Unmarshal(slateV4JSON, &slateV4); err != nil {
		assert.Error(t, err)
	}
	assert.True(t, slateV4.Ver.Version != 0)

	// Check default for num parts
	assert.Equal(t, slateV4.NumParts, uint8(2))
}

func TestMarshalSlateV4(t *testing.T) {
	slateV4JSON, _ := ioutil.ReadFile("test_data/v4_raw.slate")
	var slateV4 SlateV4
	if err := json.Unmarshal(slateV4JSON, &slateV4); err != nil {
		assert.Error(t, err)
	}
	serializedSlateV4, err := json.Marshal(slateV4)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, slateV4JSON, serializedSlateV4)
}
