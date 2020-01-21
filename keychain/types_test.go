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

package keychain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifier(t *testing.T) {
	// Identifier must be serialized/marshaled to string
	identifier := Identifier{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	b, err := json.Marshal(identifier)
	assert.Nil(t, err)
	assert.Equal(t, "\"0200000000000000000000000000000000\"", string(b))

	// Identifier must be deserialized/unmarshaled from string
	identifierString := "\"0200000000000000000000000000000000\""
	var unmarshaledIdentifier Identifier
	err = json.Unmarshal([]byte(identifierString), &unmarshaledIdentifier)
	assert.Nil(t, err)
	assert.Equal(t, identifier, unmarshaledIdentifier)

}
