// Copyright 2018 BlockCypher
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wallet

import (
	"fmt"
	"strings"
)

// Output represent the outputs in the wallet
type Output struct {
	RootKeyID  string
	KeyID      string
	NChild     uint
	Value      uint
	Status     string
	Height     uint
	LockHeight uint
	Coinbase   bool
	Block      string
}

// SendTXArgs are the send parameters
type SendTXArgs struct {
	Amount                    uint64 `json:"amount"`
	MinimumConfirmations      int    `json:"minimum_confirmations"`
	Method                    string `json:"method"`
	Dest                      string `json:"dest"`
	MaxOutputs                int    `json:"max_outputs"`
	NumChangeOutputs          int    `json:"num_change_outputs"`
	SelectionStrategyIsUseAll bool   `json:"selection_strategy_is_use_all"`
	Message                   string `json:"message,omitempty"`
}

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
