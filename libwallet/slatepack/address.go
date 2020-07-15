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
	"crypto/rand"
	"encoding/json"
	"errors"

	"github.com/blockcypher/libgrin/core/consensus"
	"github.com/btcsuite/btcutil/bech32"
)

// SlatepackAddress is the definition of a Slatepack address
type SlatepackAddress struct {
	// Human-readable prefix
	HRP string
	// ed25519 Public key, to be bech32 encoded,
	// interpreted as tor address or converted
	// to an X25519 public key for encrypting
	// slatepacks
	PubKey ed25519.PublicKey
}

// New slatepack with default hrp
func New(pubKey ed25519.PublicKey, chainType consensus.ChainType) SlatepackAddress {
	var hrp string
	switch chainType {
	case consensus.Mainnet:
		hrp = "grin"
	default:
		hrp = "tgrin"
	}
	return SlatepackAddress{
		HRP:    hrp,
		PubKey: pubKey,
	}
}

// Random create a new slatepack address with a random key
func Random(chainType consensus.ChainType) SlatepackAddress {
	var hrp string
	switch chainType {
	case consensus.Mainnet:
		hrp = "grin"
	default:
		hrp = "tgrin"
	}
	bytes := make([]byte, 32)
	rand.Read(bytes)
	privKey := ed25519.NewKeyFromSeed(bytes)
	pubKey := privKey.Public().(ed25519.PublicKey)

	return SlatepackAddress{
		HRP:    hrp,
		PubKey: pubKey,
	}
}

// MarshalJSON is a custom marshaler for slatepack address
func (sa SlatepackAddress) MarshalJSON() ([]byte, error) {
	var pubKeyBase32 []byte
	pubKeyBase32, err := bech32.ConvertBits(sa.PubKey, 8, 5, true)
	if err != nil {
		return nil, err
	}
	encoded, err := bech32.Encode(sa.HRP, pubKeyBase32)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(encoded)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalJSON is a custom unnmarshaler for slatepack address
func (sa *SlatepackAddress) UnmarshalJSON(b []byte) error {
	var encoded string
	if err := json.Unmarshal(b, &encoded); err != nil {
		return err
	}
	hrp, decoded, err := bech32.Decode(encoded)
	if err != nil {
		return err
	}
	if hrp != "grin" && hrp != "tgrin" {
		return errors.New("incorrect hrp for slatepack address")
	}
	sa.HRP = hrp
	dec, err := bech32.ConvertBits(decoded, 5, 8, false)
	if err != nil {
		return err
	}
	sa.PubKey = dec
	return nil
}
