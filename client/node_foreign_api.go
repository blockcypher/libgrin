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

	"github.com/blockcypher/libgrin/v4/api"
	"github.com/blockcypher/libgrin/v4/core"
	"github.com/blockcypher/libgrin/v5/pool"
	log "github.com/sirupsen/logrus"
)

// NodeForeignAPI represents the node foreign API (v2)
type NodeForeignAPI struct {
	client RPCHTTPClient
}

// NewNodeForeignAPI creates a new node foreign API
func NewNodeForeignAPI(url string) *NodeForeignAPI {
	return &NodeForeignAPI{client: RPCHTTPClient{URL: url}}
}

// GetBlock gets block details given either a height, a hash or an unspent output commitment.
// Only one parameters is needed. If multiple parameters are provided only the first one in the list is used.
func (foreign *NodeForeignAPI) GetBlock(height *uint64, hash, commit *string) (*api.BlockPrintable, error) {
	arrayParams := [3]interface{}{height, hash, commit}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("get_block", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetBlock")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var block api.BlockPrintable
	if err := json.Unmarshal(result.Ok, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

// GetHeader gets block header given either a height, a hash or an unspent output commitment.
// Only one parameters is needed. If multiple parameters are provided only the first one in the list is used.
func (foreign *NodeForeignAPI) GetHeader(height *uint64, hash, commit *string) (*api.BlockHeaderPrintable, error) {
	arrayParams := [3]interface{}{height, hash, commit}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("get_header", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetHeader")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var header api.BlockHeaderPrintable
	if err := json.Unmarshal(result.Ok, &header); err != nil {
		return nil, err
	}
	return &header, nil
}

// GetKernel returns a LocatedTxKernel based on the kernel excess. The min_height and max_height parameters are both optional.
// If not supplied, min_height will be set to 0 and max_height will be set to the head of the chain.
// The method will start at the block height max_height and traverse the kernel MMR backwards, until either the kernel
// is found or min_height is reached.
func (foreign *NodeForeignAPI) GetKernel(excess string, minHeight, maxHeight *uint64) (*api.LocatedTxKernel, error) {
	arrayParams := [3]interface{}{excess, minHeight, maxHeight}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("get_kernel", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetKernel")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var locatedTxKernel api.LocatedTxKernel
	if err := json.Unmarshal(result.Ok, &locatedTxKernel); err != nil {
		return nil, err
	}
	return &locatedTxKernel, nil
}

// GetOutputs retrieves details about specifics outputs. Supports retrieval of multiple outputs in a single request.
// Support retrieval by both commitment string and block height.
func (foreign *NodeForeignAPI) GetOutputs(commit *[]string, startHeight, endHeight *uint64, includeProof, includeMerkleProof *bool) (*[]api.OutputPrintable, error) {
	arrayParams := [5]interface{}{commit, startHeight, endHeight, includeProof, includeMerkleProof}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("get_outputs", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetOutputs")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var outputs []api.OutputPrintable
	if err := json.Unmarshal(result.Ok, &outputs); err != nil {
		return nil, err
	}
	return &outputs, nil
}

// GetPMMRIndices retrieves the PMMR indices based on the provided block height(s).
func (foreign *NodeForeignAPI) GetPMMRIndices(startBlockHeight uint64, endHBlockHeight *uint64) (*api.OutputListing, error) {
	arrayParams := [2]interface{}{startBlockHeight, endHBlockHeight}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("get_pmmr_indices", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetPMMRIndices")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var outputListing api.OutputListing
	if err := json.Unmarshal(result.Ok, &outputListing); err != nil {
		return nil, err
	}
	return &outputListing, nil
}

// GetPoolSize returns the number of transaction in the transaction pool.
func (foreign *NodeForeignAPI) GetPoolSize() (*uint, error) {
	envl, err := foreign.client.Request("get_pool_size", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetPoolSize")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var poolSize uint
	if err := json.Unmarshal(result.Ok, &poolSize); err != nil {
		return nil, err
	}
	return &poolSize, nil
}

// GetStempoolSize returns the number of transaction in the stem transaction pool.
func (foreign *NodeForeignAPI) GetStempoolSize() (*uint, error) {
	envl, err := foreign.client.Request("get_stempool_size", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetStempoolSize")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var poolSize uint
	if err := json.Unmarshal(result.Ok, &poolSize); err != nil {
		return nil, err
	}
	return &poolSize, nil
}

// GetTip returns details about the state of the current fork tip.
func (foreign *NodeForeignAPI) GetTip() (*api.Tip, error) {
	envl, err := foreign.client.Request("get_tip", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetTip")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var tip api.Tip
	if err := json.Unmarshal(result.Ok, &tip); err != nil {
		return nil, err
	}
	return &tip, nil
}

// GetUnconfirmedTransactions returns the unconfirmed transactions in the transaction pool.
// Will not return transactions in the stempool.
func (foreign *NodeForeignAPI) GetUnconfirmedTransactions() (*[]pool.PoolEntry, error) {
	envl, err := foreign.client.Request("get_unconfirmed_transactions", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetUnconfirmedTransactions")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var poolEntries []pool.PoolEntry
	if err := json.Unmarshal(result.Ok, &poolEntries); err != nil {
		return nil, err
	}
	return &poolEntries, nil
}

// GetUnspentOutputs is an UTXO traversal. Retrieves last utxos since a start_index until a max.
func (foreign *NodeForeignAPI) GetUnspentOutputs(startIndex uint64, endIndex *uint64, max uint64, includeProof *bool) (*api.OutputListing, error) {
	arrayParams := [4]interface{}{startIndex, endIndex, max, includeProof}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return nil, err
	}
	envl, err := foreign.client.Request("get_unspent_outputs", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetUnspentOutputs")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var unspentOutputs api.OutputListing
	if err := json.Unmarshal(result.Ok, &unspentOutputs); err != nil {
		return nil, err
	}
	return &unspentOutputs, nil
}

// GetVersion returns details about the state of the current fork tip.
func (foreign *NodeForeignAPI) GetVersion() (*api.Version, error) {
	envl, err := foreign.client.Request("get_version", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetVersion")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var version api.Version
	if err := json.Unmarshal(result.Ok, &version); err != nil {
		return nil, err
	}
	return &version, nil
}

// PushTransaction pushes a new transaction to our local transaction pool.
func (foreign *NodeForeignAPI) PushTransaction(tx core.Transaction, fluff *bool) error {
	arrayParams := [1]interface{}{tx}
	paramsBytes, err := json.Marshal(arrayParams)
	if err != nil {
		return err
	}
	envl, err := foreign.client.Request("push_transaction", paramsBytes)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("NodeForeignAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during PushTransaction")
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
