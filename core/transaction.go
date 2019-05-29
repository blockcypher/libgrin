// Copyright 2019 BlockCypher
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

// Transaction represents a transaction
type Transaction struct {
	// The kernel "offset" k2
	// excess is k1G after splitting the key k = k1 + k2
	Offset JSONableSlice `json:"offset"`
	// The transaction body - inputs/outputs/kernels
	Body TransactionBody `json:"body"`
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

// Input is a transaction input.
//
// Primarily a reference to an output being spent by the transaction.
type Input struct {
	// The features of the output being spent.
	// We will check maturity for coinbase output.
	Features string `json:"features"`
	// The commit referencing the output being spent.
	Commit JSONableSlice `json:"commit"`
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
	Features string `json:"features"`
	// The homomorphic commitment representing the output amount
	Commit JSONableSlice `json:"commit"`
	// A proof that the commitment is in the right range
	Proof JSONableSlice `json:"proof"`
}

// TxKernel is a proof that a transaction sums to zero. Includes both the transaction's
// Pedersen commitment and the signature, that guarantees that the commitments
// amount to zero.
// The signature signs the fee and the lock_height, which are retained for
// signature validation.
type TxKernel struct {
	// Options for a kernel's structure or use
	Features string `json:"features"`
	// Fee originally included in the transaction this proof is for.
	Fee uint64 `json:"fee"`
	// This kernel is not valid earlier than lock_height blocks
	// The max lock_height of all *inputs* to this transaction
	LockHeight uint64 `json:"lock_height"`
	// Remainder of the sum of all transaction commitments. If the transaction
	// is well formed, amounts components should sum to zero and the excess
	// is hence a valid public key.
	Excess JSONableSlice `json:"excess"`
	// The signature proving the excess is a valid public key, which signs
	// the transaction fee.
	ExcessSig JSONableSlice `json:"excess_sig"`
}
