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
	"github.com/blockcypher/libgrin/core"
	"github.com/blockcypher/libgrin/keychain"
	"github.com/blockcypher/libgrin/libwallet/slateversions"
)

// SendTXArgs Send TX API Args
// TODO: This is here to ensure the legacy V1 API remains intact
// remove this when v1 api is removed
type SendTXArgs struct {
	// amount to send
	Amount uint64 `json:"amount"`
	// minimum confirmations
	MinimumConfirmations uint64 `json:"minimum_confirmations"`
	// payment method
	Method string `json:"method"`
	// destination url
	Dest string `json:"dest"`
	// Max number of outputs
	MaxOutputs uint `json:"max_outputs"`
	// Number of change outputs to generate
	NumChangeOutputs uint `json:"num_change_outputs"`
	// whether to use all outputs (combine)
	SelectionStrategyIsUseAll bool `json:"selection_strategy_is_use_all"`
	// Optional message, that will be signed
	Message *string `json:"message"`
	// Optional slate version to target when sending
	TargetSlateVersion *uint16 `json:"target_slate_version"`
}

// InitTxArgs is V3 Init / Send TX API Args
type InitTxArgs struct {
	// The human readable account name from which to draw outputs
	// for the transaction, overriding whatever the active account is as set via the
	// [`set_active_account`](../grin_wallet_api/owner/struct.Owner.html#method.set_active_account) method.
	SrcAcctName *string `json:"src_acct_name"`
	// The amount to send, in nanogrins. (`1 G = 1_000_000_000nG`)
	Amount core.Uint64 `json:"amount"`
	// The minimum number of confirmations an output
	// should have in order to be included in the transaction.
	MinimumConfirmations core.Uint64 `json:"minimum_confirmations"`
	// By default, the wallet selects as many inputs as possible in a
	// transaction, to reduce the Output set and the fees. The wallet will attempt to spend
	// include up to `max_outputs` in a transaction, however if this is not enough to cover
	// the whole amount, the wallet will include more outputs. This parameter should be considered
	// a soft limit.
	MaxOutputs uint32 `json:"max_outputs"`
	// The target number of change outputs to create in the transaction.
	// The actual number created will be `num_change_outputs` + whatever remainder is needed.
	NumChangeOutputs uint32 `json:"num_change_outputs"`
	// If `true`, attempt to use up as many outputs as
	// possible to create the transaction, up the 'soft limit' of `max_outputs`. This helps
	// to reduce the size of the UTXO set and the amount of data stored in the wallet, and
	// minimizes fees. This will generally result in many inputs and a large change output(s),
	// usually much larger than the amount being sent. If `false`, the transaction will include
	// as many outputs as are needed to meet the amount, (and no more) starting with the smallest
	// value outputs.
	SelectionStrategyIsUseAll bool `json:"selection_strategy_is_use_all"`
	// An optional participant message to include alongside the sender's public
	// ParticipantData within the slate. This message will include a signature created with the
	// sender's private excess value, and will be publically verifiable. Note this message is for
	// the convenience of the participants during the exchange; it is not included in the final
	// transaction sent to the chain. The message will be truncated to 256 characters.
	Message *string `json:"message"`
	// Optionally set the output target slate version (acceptable
	// down to the minimum slate version compatible with the current. If `None` the slate
	// is generated with the latest version.
	TargetSlateVersion *uint16 `json:"target_slate_version"`
	// Number of blocks from current after which TX should be ignored
	TTLBlocks *core.Uint64 `json:"ttl_blocks"`
	// If set, require a payment proof for the particular recipient
	PaymentProofRecipientAddress *string `json:"payment_proof_recipient_address"`
	// If true, just return an estimate of the resulting slate, containing fees and amounts
	// locked without actually locking outputs or creating the transaction. Note if this is set to
	// 'true', the amount field in the slate will contain the total amount locked, not the provided
	// transaction amount
	EstimateOnly *bool `json:"estimate_only"`
	// Sender arguments. If present, the underlying function will also attempt to send the
	// transaction to a destination and optionally finalize the result
	SendArgs *InitTxSendArgs `json:"send_args"`
}

// InitTxSendArgs is the send TX API Args, for convenience functionality that inits the transaction and sends
// in one go
type InitTxSendArgs struct {
	// The transaction method. Can currently be 'http' or 'keybase'.
	Method string `json:"method"`
	// The destination, contents will depend on the particular method
	Dest string `json:"dest"`
	// Whether to finalize the result immediately if the send was successful
	Finalize bool `json:"finalize"`
	// Whether to post the transaction if the send and finalize were successful
	PostTx bool `json:"post_tx"`
	// Whether to use dandelion when posting. If false, skip the dandelion relay
	Fluff bool `json:"fluff"`
}

// IssueInvoiceTxArgs are the v2 Issue Invoice Tx Args
type IssueInvoiceTxArgs struct {
	// The human readable account name to which the received funds should be added
	// overriding whatever the active account is as set via the
	// [`set_active_account`](../grin_wallet_api/owner/struct.Owner.html#method.set_active_account) method.
	DestAcctName *string `json:"dest_acct_name"`
	// The invoice amount in nanogrins. (`1 G = 1_000_000_000nG`)
	Amount core.Uint64 `json:"amount"`
	// Optional message, that will be signed
	Message *string `json:"message"`
	// Optionally set the output target slate version (acceptable
	// down to the minimum slate version compatible with the current. If `None` the slate
	// is generated with the latest version.
	TargetSlateVersion *uint16 `json:"target_slate_version"`
}

// BlockFees are the fees in block to use for coinbase amount calculation
type BlockFees struct {
	// fees
	Fees core.Uint64 `json:"fees"`
	// height
	Height core.Uint64 `json:"height"`
	// key id
	KeyID *keychain.Identifier `json:"key_id"`
}

// CbData is the response to build a coinbase output.
type CbData struct {
	// Output
	Output core.Output `json:"output"`
	// Kernel
	Kernel core.TxKernel `json:"kernel"`
	// Key Id
	KeyID *keychain.Identifier `json:"key_id"`
}

// OutputCommitMapping is the map Outputdata to commits
type OutputCommitMapping struct {
	// Output Data
	Output OutputData `json:"output"`
	// The commit
	Commit string `json:"commit"`
}

// NodeHeightResult is the node height result
type NodeHeightResult struct {
	// Last known height
	Height core.Uint64 `json:"height"`
	// Hash
	HeaderHash string `json:"header_hash"`
	// Whether this height was updated from the node
	UpdatedFromNode bool `json:"updated_from_node"`
}

// VersionInfo is the version request result
type VersionInfo struct {
	// API version
	ForeignAPIVersion uint16 `json:"foreign_api_version"`
	// Slate version
	SupportedSlateVersions []slateversions.SlateVersion `json:"supported_slate_versions"`
}
