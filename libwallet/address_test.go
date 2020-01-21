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
	"fmt"
	"github.com/blockcypher/libgrin/libwallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnionV3Conversion(t *testing.T) {
	onionAddress := "2a6at2obto3uvkpkitqp4wxcg6u36qf534eucbskqciturczzc5suyid"
	key, err := libwallet.PubKeyFromOnionV3(onionAddress)
	assert.Nil(t, err)
	assert.NotNil(t, key)
	fmt.Printf("Key: %q\n", key)

	outAddress, err := libwallet.OnionV3FromPubKey(key)
	assert.Nil(t, err)
	fmt.Printf("Address: %s", outAddress)
	assert.Equal(t, onionAddress, outAddress)
}
