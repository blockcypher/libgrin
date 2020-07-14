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

package client

import (
	"fmt"

	"github.com/blockcypher/libgrin/api"
	"github.com/blockcypher/libgrin/core/consensus"
)

// NodeAPI struct
type NodeAPI struct {
	URL string
}

// GetBlockReward queries the node to get the block reward with fees
func (nodeAPI *NodeAPI) GetBlockReward(blockHash string) (uint64, error) {
	block, err := nodeAPI.GetBlockByHash(blockHash)
	if err != nil {
		return 0, err
	}

	// Compute the sum of the txkernels fees
	var totalFee uint64
	for _, v := range block.Kernels {
		totalFee += v.Fee
	}
	blockRewardWithFee := consensus.Reward + totalFee
	return blockRewardWithFee, nil
}

// GetBlockByHash returns a block using the hash
func (nodeAPI *NodeAPI) GetBlockByHash(blockHash string) (*api.BlockPrintable, error) {
	var block api.BlockPrintable
	url := "http://" + nodeAPI.URL + "/v1/blocks/" + blockHash
	if err := getJSON(url, &block); err != nil {
		return nil, err
	}
	// Since we create an empty struct above we should handle the case where
	// the decode does not work properly (e.g. 404)
	if block.Header.Hash == "" {
		// Error during getJSON
		return nil, fmt.Errorf("error during GetJSON")
	}
	return &block, nil
}

// GetBlockByHeight returns a block using the height
func (nodeAPI *NodeAPI) GetBlockByHeight(height uint64) (*api.BlockPrintable, error) {
	var block api.BlockPrintable
	url := fmt.Sprintf("http://%s/v1/blocks/%d", nodeAPI.URL, height)
	if err := getJSON(url, &block); err != nil {
		return nil, err
	}
	// Since we create an empty struct above we should handle the case where
	// the decode does not work properly (e.g. 404)
	if block.Header.Hash == "" {
		// Error during getJSON
		return nil, fmt.Errorf("error during GetJSON")
	}
	return &block, nil
}

// GetStatus returns the node status
func (nodeAPI *NodeAPI) GetStatus() (*api.Status, error) {
	var status api.Status
	url := "http://" + nodeAPI.URL + "/v1/status"
	if err := getJSON(url, &status); err != nil {
		return nil, err
	}
	return &status, nil
}
