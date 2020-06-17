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

package libwallet

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSlatepackVersion(t *testing.T) {
	verb := []byte(`"1.0"`)
	var ver SlatepackVersion
	if err := json.Unmarshal(verb, &ver); err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, ver.Major, uint8(1))
	assert.Equal(t, ver.Minor, uint8(0))
}

func TestMarshalSlatepackVersion(t *testing.T) {
	ver := SlatepackVersion{Major: 1, Minor: 0}
	bytes, err := json.Marshal(ver)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, string(bytes), "\"1.0\"")
}

func TestUnmarshalSlatepack(t *testing.T) {
	// TODO
}

func TestMarshalSlatepack(t *testing.T) {
	// TODO
}
