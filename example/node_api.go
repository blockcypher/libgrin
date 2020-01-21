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

package example

import (
	"fmt"

	"github.com/blockcypher/libgrin/api"
	"github.com/blockcypher/libgrin/core/consensus"
	log "github.com/sirupsen/logrus"
)

// API struct
type grinAPI struct {
	GrinServerAPI string
}

// GetBlockReward queries the node to get the block reward with fees
func (grinAPI *grinAPI) GetBlockReward(blockHash string) (uint64, error) {
	block, err := grinAPI.GetBlockByHash(blockHash)
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
func (grinAPI *grinAPI) GetBlockByHash(blockHash string) (*api.BlockPrintable, error) {
	var block api.BlockPrintable
	url := "http://" + grinAPI.GrinServerAPI + "/v1/blocks/" + blockHash
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

// GetBlockByHash returns a block using the height
func (grinAPI *grinAPI) GetBlockByHeight(height uint64) (*api.BlockPrintable, error) {
	var block api.BlockPrintable
	url := fmt.Sprintf("http://%s/v1/blocks/%d", grinAPI.GrinServerAPI, height)
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

// GetBlockByHash returns a block using the hash
func (grinAPI *grinAPI) GetStatus() (*api.Status, error) {
	var status api.Status
	url := "http://" + grinAPI.GrinServerAPI + "/v1/status"
	if err := getJSON(url, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (grinAPI *grinAPI) GetTargetDifficultyAndHashrates(status *api.Status) (uint64, float64, float64, error) {
	if status == nil {
		return 0, 0, 0, fmt.Errorf("status is nil")
	}
	lastBlock, err := grinAPI.GetBlockByHash(status.Tip.LastBlockPushed)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("API: Cannot get last block from Grin Node API")
		return 0, 0, 0, fmt.Errorf("API: Cannot get last block from Grin Node API: %s", err)
	}
	previousBlock, err := grinAPI.GetBlockByHash(status.Tip.PrevBlockToLast)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("API: Cannot get previous block from Grin Node API")
		return 0, 0, 0, fmt.Errorf("API: Cannot get previous block from Grin Node API: %s", err)
	}
	targetDifficulty := lastBlock.Header.TotalDifficulty - previousBlock.Header.TotalDifficulty

	primaryHashrate, secondaryHashrate, err := grinAPI.computePrimarySecondaryNetworkHashrates(lastBlock.Header.Height)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("API: Cannot compute hashrates")
		return 0, 0, 0, fmt.Errorf("API: Cannot compute hashrates: %s", err)
	}
	return targetDifficulty, primaryHashrate, secondaryHashrate, nil
}

// computePrimarySecondaryNetworkHashrates is a shameless copy of https://github.com/grin-pool/grin-pool/blob/f3e21651b26f759b5303e9f7f702a55d55f63ca4/grin-py/grinlib/grinstats.py#L93
func (grinAPI *grinAPI) computePrimarySecondaryNetworkHashrates(endHeight uint64) (float64, float64, error) {
	var blocks []api.BlockPrintable
	var startHeight uint64
	if endHeight < consensus.DifficultyAdjustWindow {
		startHeight = 0
	}
	startHeight = endHeight - consensus.DifficultyAdjustWindow
	for height := startHeight; height <= endHeight; height++ {
		block, err := grinAPI.GetBlockByHeight(height)
		if err != nil {
			return 0, 0, err
		}
		blocks = append(blocks, *block)
	}

	hashrates, err := estimateAllHashrates(blocks)
	if err != nil {
		return 0, 0, err
	}
	var primaryHashrate float64
	var secondaryHashrate float64
	if hashrate, ok := hashrates[consensus.DefaultMinEdgeBits]; ok {
		primaryHashrate = hashrate
	}
	if hashrate, ok := hashrates[consensus.SecondPoWEdgeBits]; ok {
		secondaryHashrate = hashrate
	}
	return primaryHashrate, secondaryHashrate, nil
}

func estimateAllHashrates(blocks []api.BlockPrintable) (map[uint8]float64, error) {
	// Calculate the gps for each graph size in the recent blocks list
	// Based on jaspervdm code - https://github.com/jaspervdm/grin_mining_sim

	// Avoid crashing on empty slice
	if len(blocks) == 0 {
		return make(map[uint8]float64), nil
	}

	height := blocks[len(blocks)-1].Header.Height
	// Get the difficulty of the most recent block
	difficulty := blocks[len(blocks)-1].Header.TotalDifficulty - blocks[len(blocks)-2].Header.TotalDifficulty
	// Get secondary_scaling value for the most recent block
	secondaryScaling := blocks[len(blocks)-1].Header.SecondaryScaling
	// Count the total number of each solution size in the window
	counts := make(map[uint8]int)
	var countPrimary int
	for _, block := range blocks {
		counts[block.Header.EdgeBits]++
		if block.Header.EdgeBits != consensus.SecondPoWEdgeBits {
			countPrimary++
		}
	}
	// ratios
	q := float64(consensus.SecondaryPoWRatio(height)) / 100.0
	r := 1.0 - q

	// Calculate the GPS
	gpsMap := make(map[uint8]float64)
	for edgeBits, count := range counts {
		var gps float64
		if edgeBits == consensus.SecondPoWEdgeBits {
			gps = float64(42*float64(difficulty)*q) / float64(secondaryScaling) / 60
		} else {
			countRatio := float64(count) / float64(countPrimary)
			gps = 42 * float64(difficulty) * r * countRatio / float64(consensus.GraphWeight(consensus.Mainnet, 0, edgeBits)) / 60
		}
		gpsMap[edgeBits] = gps
	}
	return gpsMap, nil
}
