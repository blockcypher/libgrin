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

package slatepack

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalSlatepackAddress(t *testing.T) {
	bytes := []byte(`"tgrin1ntqks7jl77u00q9x4fjt4ydupx8fad0py9c5kx86v30snqxn2mvqacwa3d"`)
	var sa SlatepackAddress
	if err := json.Unmarshal(bytes, &sa); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, "tgrin", sa.HRP)
	assert.Equal(t, ed25519.PublicKey{154, 193, 104, 122, 95, 247, 184, 247, 128, 166, 170, 100, 186, 145, 188, 9, 142, 158, 181, 225, 33, 113, 75, 24, 250, 100, 95, 9, 128, 211, 86, 216}, sa.PubKey)
}

func TestMarshalSlatepackAddress(t *testing.T) {
	sa := SlatepackAddress{
		HRP:    "tgrin",
		PubKey: ed25519.PublicKey{154, 193, 104, 122, 95, 247, 184, 247, 128, 166, 170, 100, 186, 145, 188, 9, 142, 158, 181, 225, 33, 113, 75, 24, 250, 100, 95, 9, 128, 211, 86, 216},
	}
	saEncoded, err := json.Marshal(sa)
	assert.NoError(t, err)
	assert.Equal(t, "\"tgrin1ntqks7jl77u00q9x4fjt4ydupx8fad0py9c5kx86v30snqxn2mvqacwa3d\"", string(saEncoded))
}
