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
	"encoding/base32"
	"errors"
	"strings"

	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/sha3"
)

// PubKeyFromOnionV3 returns the ed25519 public key represented in an onion address
func PubKeyFromOnionV3(onionAddress string) (ed25519.PublicKey, error) {
	input := strings.ToUpper(onionAddress)
	if strings.HasPrefix(input, "HTTP://") || strings.HasPrefix(input, "HTTPS://") {
		input = strings.Replace(input, "HTTP://", "", 1)
		input = strings.Replace(input, "HTTPS://", "", 1)
	}
	if strings.HasSuffix(input, ".ONION") {
		input = strings.Replace(input, ".ONION", "", 1)
	}

	// for now, just check input is the right length and try and decode from base32
	if len(input) != 56 {
		return nil, errors.New("input address is wrong length")
	}
	address, err := base32.StdEncoding.DecodeString(strings.ToUpper(input))
	if err != nil {
		return nil, errors.New("input address is not base 32")
	}

	var key ed25519.PublicKey = address[0:32]

	testV3, err := OnionV3FromPubKey(key)
	if err != nil {
		return nil, errors.New("provided onion V3 address is invalid (converting from pubkey)")
	}

	if strings.ToUpper(testV3) != input {
		return nil, errors.New("provided onion V3 address is invalid (no match)")
	}

	return key, nil
}

// OnionV3FromPubKey generates an onion address from an ed25519 public key
func OnionV3FromPubKey(pubKey ed25519.PublicKey) (string, error) {
	// calculate checksum
	checksum := make([]byte, 32)
	hasher := sha3.New256()
	hasher.Write([]byte(".onion checksum"))
	hasher.Write(pubKey)
	hasher.Write([]byte{3})
	hasher.Sum(checksum[:0])

	addressBytes := pubKey
	addressBytes = append(addressBytes, checksum[0:2]...)
	addressBytes = append(addressBytes, []byte{3}[:]...)
	ret := base32.StdEncoding.EncodeToString(addressBytes)
	return strings.ToLower(ret), nil
}
