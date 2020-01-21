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

package example

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/blockcypher/libgrin/core"
	"github.com/blockcypher/libgrin/libwallet"
	log "github.com/sirupsen/logrus"
)

var url = "http://127.0.0.1:3420/v2/owner"

// RetrieveOutputs returns a list of outputs from the active account in the
func RetrieveOutputs(includeSpent, refreshFromNode bool, txID *uint32) (bool, *[]libwallet.OutputCommitMapping, error) {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{includeSpent, refreshFromNode, txID}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return false, nil, err
	}
	envl, err := client.Request("retrieve_outputs", paramsBytes)
	if err != nil {
		return false, nil, err
	}
	if envl == nil {
		return false, nil, errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during RetrieveOutputs")
		return false, nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return false, nil, err
	}
	if result.Err != nil {
		return false, nil, errors.New(string(result.Err))
	}

	var okArray []json.RawMessage
	if err = json.Unmarshal(result.Ok, &okArray); err != nil {
		return false, nil, err
	}
	if len(okArray) < 2 {
		return false, nil, errors.New("Wrong okArray length")
	}
	var refreshedFromNode bool
	if err = json.Unmarshal(okArray[0], &refreshedFromNode); err != nil {
		return false, nil, err
	}
	var txLogEntries []libwallet.OutputCommitMapping
	if err := json.Unmarshal(okArray[1], &txLogEntries); err != nil {
		return false, nil, err
	}

	return refreshedFromNode, &txLogEntries, nil
}

// RetrieveTxs returns a list of Transaction Log Entries from the active account in the
func RetrieveTxs(refreshFromNode bool, txID *uint32, txSlateID *uuid.UUID) (bool, *[]libwallet.TxLogEntry, error) {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{refreshFromNode, txID, txSlateID}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return false, nil, err
	}
	envl, err := client.Request("retrieve_txs", paramsBytes)
	if err != nil {
		return false, nil, err
	}
	if envl == nil {
		return false, nil, errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during RetrieveTxs")
		return false, nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return false, nil, err
	}
	if result.Err != nil {
		return false, nil, errors.New(string(result.Err))
	}

	var okArray []json.RawMessage
	if err = json.Unmarshal(result.Ok, &okArray); err != nil {
		return false, nil, err
	}
	if len(okArray) < 2 {
		return false, nil, errors.New("Wrong okArray length")
	}
	var refreshedFromNode bool
	if err = json.Unmarshal(okArray[0], &refreshedFromNode); err != nil {
		return false, nil, err
	}
	var txLogEntries []libwallet.TxLogEntry
	if err := json.Unmarshal(okArray[1], &txLogEntries); err != nil {
		return false, nil, err
	}

	return refreshedFromNode, &txLogEntries, nil
}

// RetrieveSummaryInfo returns summary information from the active account in the
func RetrieveSummaryInfo(refreshFromNode bool, minimumConfirmations uint64) (bool, *libwallet.WalletInfo, error) {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{refreshFromNode, minimumConfirmations}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return false, nil, err
	}
	envl, err := client.Request("retrieve_summary_info", paramsBytes)
	if err != nil {
		return false, nil, err
	}
	if envl == nil {
		return false, nil, errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during RetrieveSummaryInfo")
		return false, nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return false, nil, err
	}
	if result.Err != nil {
		return false, nil, errors.New(string(result.Err))
	}

	var okArray []json.RawMessage
	if err = json.Unmarshal(result.Ok, &okArray); err != nil {
		return false, nil, err
	}
	if len(okArray) < 2 {
		return false, nil, errors.New("Wrong okArray length")
	}
	var refreshedFromNode bool
	if err = json.Unmarshal(okArray[0], &refreshedFromNode); err != nil {
		return false, nil, err
	}
	var walletInfo libwallet.WalletInfo
	if err := json.Unmarshal(okArray[1], &walletInfo); err != nil {
		return false, nil, err
	}
	return refreshedFromNode, &walletInfo, nil
}

type initSendTxArgs struct {
	InitTxArgs libwallet.InitTxArgs `json:"args"`
}

// InitSendTx initiates a new transaction as the sender, creating a new Slate
// object containing the sender's inputs, change outputs, and public signature
// data.
func InitSendTx(initTxArgs libwallet.InitTxArgs) (*libwallet.Slate, error) {
	initSendTxArgs := initSendTxArgs{initTxArgs}
	client := RPCHTTPClient{URL: url}
	paramsBytes, err := json.Marshal(initSendTxArgs)
	if err != nil {
		return nil, err
	}
	envl, err := client.Request("init_send_tx", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during InitSendTx")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}

	var slate libwallet.Slate
	if err := json.Unmarshal(result.Ok, &slate); err != nil {
		return nil, err
	}
	return &slate, nil
}

// TxLockOutputs locks the outputs associated with the inputs to the transaction
// in the given Slate, making them unavailable for use in further transactions.
func TxLockOutputs(slate *libwallet.Slate, participantID uint) error {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{slate, participantID}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := client.Request("tx_lock_outputs", paramsBytes)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during TxLockOutputs")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	return nil
}

// FinalizeTx finalizes a transaction, after all parties have filled in both rounds of Slate generation.
func FinalizeTx(slateIn libwallet.Slate) (*libwallet.Slate, error) {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{slateIn}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	envl, err := client.Request("finalize_tx", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during FinalizeTx")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}

	var slate libwallet.Slate
	if err := json.Unmarshal(result.Ok, &slate); err != nil {
		return nil, err
	}
	return &slate, nil
}

// PostTx posts a completed transaction to the listening node for validation and
// inclusion in a block for mining.
func PostTx(tx core.Transaction, fluff bool) error {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{tx, fluff}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := client.Request("post_tx", paramsBytes)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during PostTx")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	return nil
}

// CancelTx cancels a transaction.
func CancelTx(txID *uint32, txSlateID *uuid.UUID) error {
	client := RPCHTTPClient{URL: url}
	params := []interface{}{txID, txSlateID}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := client.Request("cancel_tx", paramsBytes)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during CancelTx")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	return nil
}

// NodeHeight retrieves the last known height known by the node.
func NodeHeight() (*libwallet.NodeHeightResult, error) {
	client := RPCHTTPClient{URL: url}
	envl, err := client.Request("node_height", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during NodeHeight")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var nodeHeightResult libwallet.NodeHeightResult
	if err := json.Unmarshal(result.Ok, &nodeHeightResult); err != nil {
		return nil, err
	}
	return &nodeHeightResult, nil
}
