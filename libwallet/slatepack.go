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
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Slatepack is a basic slatepack definition
type Slatepack struct {
	// Required Fields
	// Versioning info
	Slatepack SlatepackVersion `json:"slatepack"`
	// Delivery Mode, 0 = plain_text, 1 = age encrypted
	Mode uint8 `json:"mode"`

	// Optional Fields
	// Optional Sender address
	Sender *string `json:"sender,omitempty"`

	// Encrypted metadata, to be serialized into payload only
	// shouldn't be accessed directly
	// Encrypted metadata
	EncryptedMeta SlatepackEncMetadata `json:"encrypted_meta,omitempty"`

	// Payload (e.g. slate), including encrypted metadata, if present
	// Binary payload, can be encrypted or plaintext
	Payload string `json:"payload"`
}

// SlatepackVersion slatepack versionn
type SlatepackVersion struct {
	// Major
	Major uint8 `json:"major"`
	// Minor
	Minor uint8 `json:"minor"`
}

// MarshalJSON marshals the VersionCompatInfoV4 as a quoted version like {}:{}
func (v SlatepackVersion) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%d.%d", v.Major, v.Minor)
	bytes, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalJSON unmarshals a quoted version to a VersionCompatInfoV4
func (v *SlatepackVersion) UnmarshalJSON(bs []byte) error {
	var verString string
	if err := json.Unmarshal(bs, &verString); err != nil {
		return err
	}
	ver := strings.Split(verString, ".")
	if len(ver) != 2 {
		return errors.New("cannot parse version")
	}

	version, err := strconv.ParseUint(ver[0], 10, 8)
	if err != nil {
		return errors.New("cannot parse version")
	}

	v.Major = uint8(version)
	blockHeaderVersion, err := strconv.ParseUint(ver[1], 10, 8)
	if err != nil {
		return errors.New("cannot parse version")
	}
	v.Minor = uint8(blockHeaderVersion)
	return nil
}

// SlatepackEncMetadata encapsulates encrypted metadata fields
type SlatepackEncMetadata struct {
	// Encrypted Sender address, if desired
	Sender *string `json:"sender,omitempty"`
	// Recipients list, if desired (mostly for future multiparty needs)
	Recipients []string `json:"recipients,omitempty"`
}
