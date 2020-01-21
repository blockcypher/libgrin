// Copyright 2020 BlockCypher
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

package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

// Uint64 is an uint64 that can be unmarshal from a string or uint64 is
// marshal to a string
type Uint64 uint64

// MarshalJSON marshals the Uint64 as a quoted uint64 string
func (u Uint64) MarshalJSON() ([]byte, error) {
	str := strconv.FormatUint(uint64(u), 10)
	bytes, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalJSON unmarshals a quoted an uint64 or a string to an uint64 value
func (u *Uint64) UnmarshalJSON(bs []byte) error {
	var i uint64
	if err := json.Unmarshal(bs, &i); err == nil {
		*u = Uint64(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(bs, &s); err != nil {
		return errors.New("expected a string or an integer")
	}
	if err := json.Unmarshal([]byte(s), &i); err != nil {
		return err
	}
	*u = Uint64(i)
	return nil
}

// KernelFeatures is an enum of various supported kernels "features".
type KernelFeatures int

const (
	// PlainKernel kernel (the default for Grin txs).
	PlainKernel KernelFeatures = iota
	// CoinbaseKernel is a coinbase kernel.
	CoinbaseKernel
	// HeightLockedKernel is a kernel with an explicit lock height.
	HeightLockedKernel
)

func (s KernelFeatures) String() string {
	return toStringKernelFeatures[s]
}

var toStringKernelFeatures = map[KernelFeatures]string{
	PlainKernel:        "Plain",
	CoinbaseKernel:     "Coinbase",
	HeightLockedKernel: "HeightLocked",
}

var toIDKernelFeatures = map[string]KernelFeatures{
	"Plain":        PlainKernel,
	"Coinbase":     CoinbaseKernel,
	"HeightLocked": HeightLockedKernel,
}

// MarshalJSON marshals the enum as a quoted json string
func (s KernelFeatures) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringKernelFeatures[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *KernelFeatures) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toIDKernelFeatures[j]
	return nil
}

// TxKernel is a proof that a transaction sums to zero. Includes both the transaction's
// Pedersen commitment and the signature, that guarantees that the commitments
// amount to zero.
// The signature signs the fee and the lock_height, which are retained for
// signature validation.
type TxKernel struct {
	// Options for a kernel's structure or use
	Features KernelFeatures `json:"features"`
	// Fee originally included in the transaction this proof is for.
	Fee Uint64 `json:"fee"`
	// This kernel is not valid earlier than lock_height blocks
	// The max lock_height of all *inputs* to this transaction
	LockHeight Uint64 `json:"lock_height"`
	// Remainder of the sum of all transaction commitments. If the transaction
	// is well formed, amounts components should sum to zero and the excess
	// is hence a valid public key.
	Excess string `json:"excess"`
	// The signature proving the excess is a valid public key, which signs
	// the transaction fee.
	ExcessSig string `json:"excess_sig"`
}

// TransactionBody is a common abstraction for transaction and block
type TransactionBody struct {
	// List of inputs spent by the transaction.
	Inputs []Input `json:"inputs"`
	// List of outputs the transaction produces.
	Outputs []Output `json:"outputs"`
	// List of kernels that make up this transaction (usually a single kernel).
	Kernels []TxKernel `json:"kernels"`
}

// Transaction represents a transaction
type Transaction struct {
	// The kernel "offset" k2
	// excess is k1G after splitting the key k = k1 + k2
	Offset string `json:"offset"`
	// The transaction body - inputs/outputs/kernels
	Body TransactionBody `json:"body"`
}

// Input is a transaction input.
//
// Primarily a reference to an output being spent by the transaction.
type Input struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features OutputFeatures `json:"features"`
	// The commit referencing the output being spent.
	Commit string `json:"commit"`
}

// OutputFeatures is an enum of various supported outputs "features".
type OutputFeatures int

const (
	// PlainOutput output (the default for Grin txs).
	PlainOutput OutputFeatures = iota
	// CoinbaseOutput is a coinbase output.
	CoinbaseOutput
)

func (s OutputFeatures) String() string {
	return toStringOutputFeatures[s]
}

var toStringOutputFeatures = map[OutputFeatures]string{
	PlainOutput:    "Plain",
	CoinbaseOutput: "Coinbase",
}

var toIDOutputFeatures = map[string]OutputFeatures{
	"Plain":    PlainOutput,
	"Coinbase": CoinbaseOutput,
}

// MarshalJSON marshals the enum as a quoted json string
func (s OutputFeatures) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toStringOutputFeatures[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals a quoted json string to the enum value
func (s *OutputFeatures) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toIDOutputFeatures[j]
	return nil
}

// Output for a transaction, defining the new ownership of coins that are being
// transferred. The commitment is a blinded value for the output while the
// range proof guarantees the commitment includes a positive value without
// overflow and the ownership of the private key. The switch commitment hash
// provides future-proofing against quantum-based attacks, as well as providing
// wallet implementations with a way to identify their outputs for wallet
// reconstruction.
type Output struct {
	// Options for an output's structure or use
	Features OutputFeatures `json:"features"`
	// The homomorphic commitment representing the output amount
	Commit string `json:"commit"`
	// A proof that the commitment is in the right range
	Proof string `json:"proof"`
}

// An OutputIdentifier can be build from either an input _or_ an output and
// contains everything we need to uniquely identify an output being spent.
// Needed because it is not sufficient to pass a commitment around.
type OutputIdentifier struct {
	// Output features (coinbase vs. regular transaction output)
	// We need to include this when hashing to ensure coinbase maturity can be
	// enforced.
	Features OutputFeatures `json:"features"`
	// Output commitment
	Commit string `json:"commit"`
}
