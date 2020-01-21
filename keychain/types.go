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
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// IdentifierSize is the size of an identifier in bytes
const IdentifierSize uint = 17

// Identifier is a  semi-opaque structure (just bytes) to track keys
// within the Keychain.
type Identifier [IdentifierSize]uint8

// UnmarshalJSON is a custom unmarshal function for Identifier
func (i *Identifier) UnmarshalJSON(b []byte) error {
	if len(b) != 36 {
		return errors.New("wrong identifier length")
	}
	b = b[1:]
	b = b[:len(b)-1]
	identifierBytes, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}
	if len(identifierBytes) != 17 {
		return errors.New("wrong identifier length")
	}
	copy(i[:], identifierBytes)
	return nil
}

// MarshalJSON is a custom marshal function for Identifier
func (i Identifier) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	str := strings.Join(strings.Fields(fmt.Sprintf("%x", i)), "")
	buffer.WriteString(str)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}
