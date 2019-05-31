// Copyright 2019 BlockCypher
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
	"fmt"
	"strings"
)

// CurrentSlateVersion is The most recent version of the slate
const CurrentSlateVersion uint16 = 2

// SlateVersion represents the slate version
type SlateVersion int

const (
	// V0 (first version)
	V0 SlateVersion = iota
	// V1 (like V0 but with version field)
	V1
	// V2 (most current)
	V2
)

// JSONableSlice is a slice that is not represented as a base58 when serialized
type JSONableSlice []uint8

// MarshalJSON is the marshal function for such type
func (u JSONableSlice) MarshalJSON() ([]byte, error) {
	var result string
	if u == nil {
		result = "null"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	}
	return []byte(result), nil
}
