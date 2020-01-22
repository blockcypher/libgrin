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
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/blockcypher/libgrin/core"
	"github.com/blockcypher/libgrin/libwallet"
	"github.com/btcsuite/btcd/btcec"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

// SecureOwnerAPI represent the wallet owner API (v3)
type SecureOwnerAPI struct {
	client          RPCHTTPClient
	token           string
	privateKey      btcec.PrivateKey
	PublicKey       btcec.PublicKey
	ServerPublicKey *btcec.PublicKey
	sharedSecret    []byte
}

// NewSecureOwnerAPI creates a new owner API
func NewSecureOwnerAPI(url string) *SecureOwnerAPI {
	return &SecureOwnerAPI{client: RPCHTTPClient{URL: url}}
}

func newKey() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	priv, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return priv, &priv.PublicKey, nil
}

// Init initalize the secure owner API
func (owner *SecureOwnerAPI) Init() error {
	ecdsaPrivateKey, ecdsaPublicKey, err := newKey()
	if err != nil {
		return err
	}
	privkey := btcec.PrivateKey(*ecdsaPrivateKey)
	pubkey := btcec.PublicKey(*ecdsaPublicKey)
	serverPubKeyHex, err := owner.InitSecureAPI(pubkey.SerializeCompressed())
	if err != nil {
		return err
	}
	serverPubKey, err := hex.DecodeString(serverPubKeyHex)
	if err != nil {
		return err
	}
	owner.ServerPublicKey, err = btcec.ParsePubKey(serverPubKey, btcec.S256())
	if err != nil {
		return err
	}

	owner.privateKey = privkey
	owner.PublicKey = pubkey
	owner.sharedSecret = btcec.GenerateSharedSecret(&privkey, owner.ServerPublicKey)
	return nil
}

// Open is an helper function to open the wallet and set the token
func (owner *SecureOwnerAPI) Open(name *string, password string) error {
	token, err := owner.OpenWallet(name, password)
	if err != nil {
		return err
	}
	owner.token = token
	return nil
}

// Close is an helper function to close the wallet and free the token from memory
func (owner *SecureOwnerAPI) Close(name *string) error {
	if err := owner.CloseWallet(name); err != nil {
		return err
	}
	owner.token = ""
	return nil
}

// InitSecureAPI Initializes the secure JSON-RPC API. This function must be called and a shared key
// established before any other OwnerAPI JSON-RPC function can be called.
func (owner *SecureOwnerAPI) InitSecureAPI(pubKey []byte) (string, error) {
	hexPubKey := hex.EncodeToString(pubKey)
	params := struct {
		PublicKey string `json:"ecdh_pubkey"`
	}{
		PublicKey: hexPubKey,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	envl, err := owner.client.Request("init_secure_api", paramsBytes)
	if err != nil {
		return "", err
	}

	if envl == nil {
		return "", errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during InitSecureAPI")
		return "", errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return "", err
	}
	if result.Err != nil {
		return "", errors.New(string(result.Err))
	}
	serverPubKey := strings.Trim(string(result.Ok), "\"")
	return serverPubKey, nil
}

// OpenWallet `opens` a wallet, populating the internal keychain with the encrypted seed, and optionally
// returning a `keychain_mask` token to the caller to provide in all future calls.
// If using a mask, the seed will be stored in-memory XORed against the `keychain_mask`, and
// will not be useable if the mask is not provided.
func (owner *SecureOwnerAPI) OpenWallet(name *string, password string) (string, error) {
	params := struct {
		Name     *string `json:"name"`
		Password string  `json:"password"`
	}{
		Name:     name,
		Password: password,
	}

	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	envl, err := owner.client.EncryptedRequest("open_wallet", paramsBytes, owner.sharedSecret)
	if err != nil {
		return "", err
	}

	if envl == nil {
		return "", errors.New("OwnerAPI: Empty RPC Response from grin-wallet")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("OwnerAPI: RPC Error during OpenWallet")
		return "", errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return "", err
	}
	if result.Err != nil {
		return "", errors.New(string(result.Err))
	}
	token := strings.Trim(string(result.Ok), "\"")
	return token, nil

}

// CloseWallet close a wallet, removing the master seed from memory.
func (owner *SecureOwnerAPI) CloseWallet(name *string) error {
	params := struct {
		Name *string `json:"name"`
	}{
		Name: name,
	}

	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}

	envl, err := owner.client.EncryptedRequest("close_wallet", paramsBytes, owner.sharedSecret)
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
		}).Error("OwnerAPI: RPC Error during CloseWallet")
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

// RetrieveOutputs returns a list of outputs from the active account in the
func (owner *SecureOwnerAPI) RetrieveOutputs(includeSpent, refreshFromNode bool, txID *uint32) (bool, *[]libwallet.OutputCommitMapping, error) {
	params := struct {
		Token           string  `json:"token"`
		IncludeSpent    bool    `json:"include_spent"`
		RefreshFromNode bool    `json:"refresh_from_node"`
		TxID            *uint32 `json:"tx_id"`
	}{
		Token:           owner.token,
		IncludeSpent:    includeSpent,
		RefreshFromNode: refreshFromNode,
		TxID:            txID,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return false, nil, err
	}
	envl, err := owner.client.EncryptedRequest("retrieve_outputs", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) RetrieveTxs(refreshFromNode bool, txID *uint32, txSlateID *uuid.UUID) (bool, *[]libwallet.TxLogEntry, error) {
	params := struct {
		Token           string     `json:"token"`
		RefreshFromNode bool       `json:"refresh_from_node"`
		TxID            *uint32    `json:"tx_id"`
		TxSlateID       *uuid.UUID `json:"tx_slate_id"`
	}{
		Token:           owner.token,
		RefreshFromNode: refreshFromNode,
		TxID:            txID,
		TxSlateID:       txSlateID,
	}

	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return false, nil, err
	}
	envl, err := owner.client.EncryptedRequest("retrieve_txs", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) RetrieveSummaryInfo(refreshFromNode bool, minimumConfirmations uint64) (bool, *libwallet.WalletInfo, error) {
	params := struct {
		Token                string `json:"token"`
		RefreshFromNode      bool   `json:"refresh_from_node"`
		MinimumConfirmations uint64 `json:"minimum_confirmations"`
	}{
		Token:                owner.token,
		RefreshFromNode:      refreshFromNode,
		MinimumConfirmations: minimumConfirmations,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return false, nil, err
	}
	envl, err := owner.client.EncryptedRequest("retrieve_summary_info", paramsBytes, owner.sharedSecret)
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

// InitSendTx initiates a new transaction as the sender, creating a new Slate
// object containing the sender's inputs, change outputs, and public signature
// data.
func (owner *SecureOwnerAPI) InitSendTx(initTxArgs libwallet.InitTxArgs) (*libwallet.Slate, error) {
	params := struct {
		Token string               `json:"token"`
		Args  libwallet.InitTxArgs `json:"args"`
	}{
		Token: owner.token,
		Args:  initTxArgs,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	envl, err := owner.client.EncryptedRequest("init_send_tx", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) TxLockOutputs(slate libwallet.Slate, participantID uint) error {
	params := struct {
		Token         string          `json:"token"`
		Slate         libwallet.Slate `json:"slate"`
		ParticipantID uint            `json:"participant_id"`
	}{
		Token:         owner.token,
		Slate:         slate,
		ParticipantID: participantID,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := owner.client.EncryptedRequest("tx_lock_outputs", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) FinalizeTx(slateIn libwallet.Slate) (*libwallet.Slate, error) {
	params := struct {
		Token string          `json:"token"`
		Slate libwallet.Slate `json:"slate"`
	}{
		Token: owner.token,
		Slate: slateIn,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	envl, err := owner.client.EncryptedRequest("finalize_tx", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) PostTx(tx core.Transaction, fluff bool) error {
	params := struct {
		Token string           `json:"token"`
		Tx    core.Transaction `json:"tx"`
		Fluff bool             `json:"fluff"`
	}{
		Token: owner.token,
		Tx:    tx,
		Fluff: fluff,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := owner.client.EncryptedRequest("post_tx", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) CancelTx(txID *uint32, txSlateID *uuid.UUID) error {
	params := struct {
		Token     string     `json:"token"`
		TxID      *uint32    `json:"tx_id"`
		TxSlateID *uuid.UUID `json:"tx_slate_id"`
	}{
		Token:     owner.token,
		TxID:      txID,
		TxSlateID: txSlateID,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := owner.client.EncryptedRequest("cancel_tx", paramsBytes, owner.sharedSecret)
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
func (owner *SecureOwnerAPI) NodeHeight() (*libwallet.NodeHeightResult, error) {
	params := struct {
		Token string `json:"token"`
	}{
		Token: owner.token,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	envl, err := owner.client.EncryptedRequest("node_height", paramsBytes, owner.sharedSecret)
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

// SetTorConfig set the TOR configuration for this instance of the OwnerAPI,
// used during InitSendTx when send args are present and a TOR address is specified
func (owner *SecureOwnerAPI) SetTorConfig(torConfig libwallet.TorConfig) error {
	params := struct {
		TorConfig libwallet.TorConfig `json:"tor_config"`
	}{
		TorConfig: torConfig,
	}
	paramsBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	envl, err := owner.client.EncryptedRequest("set_tor_config", paramsBytes, owner.sharedSecret)
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
		}).Error("OwnerAPI: RPC Error during SetTorConfig")
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
