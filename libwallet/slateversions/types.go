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
	"bytes"
	"encoding/json"
)

// CurrentSlateVersion is the most recent version of the slate
const CurrentSlateVersion uint16 = 4

// GrinBlockHeaderVersion is the grin block header this slate is intended to be compatible with
const GrinBlockHeaderVersion uint16 = 3

// SlateVersion represents the slate version
type SlateVersion int

const (
	// V4SlateVersion is v4 (most current)
	V4SlateVersion SlateVersion = iota
)

var toStringSlateVersion = map[SlateVersion]string{
	V4SlateVersion: "V4",
}

var toIDSlateVersion = map[string]SlateVersion{
	"V4": V4SlateVersion,
}

// MarshalJSON marshals the enum as a quoted json string
func (s SlateVersion) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringSlateVersion[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *SlateVersion) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'UnknownSlateState' in this case.
	*s = toIDSlateVersion[j]
	return nil
}
