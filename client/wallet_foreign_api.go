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

package client

import (
	"encoding/json"
	"errors"

	"github.com/blockcypher/libgrin/v5/libwallet"
	"github.com/blockcypher/libgrin/v5/libwallet/slateversions"
	log "github.com/sirupsen/logrus"
)

// WalletForeignAPI represents the wallet foreign API (v2)
type WalletForeignAPI struct {
	client RPCHTTPClient
}

// NewWalletForeignAPI creates a new wallet foreign API
func NewWalletForeignAPI(url string) *WalletForeignAPI {
	return &WalletForeignAPI{client: RPCHTTPClient{URL: url}}
}

// CheckVersion returns the version capabilities of the running ForeignApi Node
func (foreign *WalletForeignAPI) CheckVersion() (*libwallet.VersionInfo, error) {
	envl, err := foreign.client.Request("check_version", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("WalletForeignAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("WalletForeignAPI: RPC Error during CheckVersion")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var versionInfo libwallet.VersionInfo
	if err := json.Unmarshal(result.Ok, &versionInfo); err != nil {
		return nil, err
	}
	return &versionInfo, nil
}

// BuildCoinbase builds a new unconfirmed coinbase output in the wallet, generally for inclusion
// in a potential new block's coinbase output during mining.
//
// All potential coinbase outputs are created as 'Unconfirmed' with the coinbase flag set.
// If a potential coinbase output is found on the chain after a wallet update, it status
// is set to Unsent and a Transaction Log Entry will be created. Note the output will be
// unspendable until the coinbase maturity period has expired.
func (foreign *WalletForeignAPI) BuildCoinbase(blockFees libwallet.BlockFees) (*libwallet.CbData, error) {
	paramsBytes, err := json.Marshal(blockFees)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("build_coinbase", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("WalletForeignAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("WalletForeignAPI: RPC Error during BuildCoinbase")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var cbData libwallet.CbData
	if err := json.Unmarshal(result.Ok, &cbData); err != nil {
		return nil, err
	}
	return &cbData, nil
}

// FinalizeTx finalizes a (standard or invoice) transaction initiated by this wallet's Owner api.
// This step assumes the paying party has completed round 1 and 2 of slate creation,
// and added their partial signatures. This wallet will verify and add their partial sig,
// then create the finalized transaction, ready to post to a node.
//
// This function also stores the final transaction in the user's wallet files for retrieval
// via the get_stored_tx function.
func (foreign *WalletForeignAPI) FinalizeTx(slate *slateversions.SlateV4) (*slateversions.SlateV4, error) {
	paramsBytes, err := json.Marshal(slate)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("finalize_Tx", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("WalletForeignAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("WalletForeignAPI: RPC Error during FinalizeTx")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var finalizedSlate slateversions.SlateV4
	if err := json.Unmarshal(result.Ok, &finalizedSlate); err != nil {
		return nil, err
	}
	return &finalizedSlate, nil
}

// ReceiveTx receives a transaction created by another party, returning the modified Slate object,
// modified with the recipient's output for the transaction amount, and public signature data.
// This slate can then be sent back to the sender to finalize the transaction via the Owner API's
// finalize_tx method.
//
// This function creates a single output for the full amount, set to a status of 'Awaiting finalization'.
// It will remain in this state until the wallet finds the corresponding output on the chain, at which point
// it will become 'Unspent'. The slate will be updated with the results of Signing round 1 and 2, adding
// the recipient's public nonce, public excess value, and partial signature to the slate.
//
// Also creates a corresponding Transaction Log Entry in the wallet's transaction log.
func (foreign *WalletForeignAPI) ReceiveTx(slate slateversions.SlateV4, destAcctName *string, dest *string) (*slateversions.SlateV4, error) {
	// TODO this is broken
	params := struct {
		Slate        slateversions.SlateV4 `json:"slate"`
		DestAcctName *string               `json:"dest_acct_name"`
		Dest         *string               `json:"dest"`
	}{
		Slate:        slate,
		DestAcctName: destAcctName,
		Dest:         dest,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("receive_tx", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("WalletForeignAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("WalletForeignAPI: RPC Error during ReceiveTx")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var receivedSlate slateversions.SlateV4
	if err := json.Unmarshal(result.Ok, &receivedSlate); err != nil {
		return nil, err
	}
	return &receivedSlate, nil
}
