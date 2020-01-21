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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/blockcypher/libgrin/api"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var blocks []api.BlockPrintable

// startTestGrinAPIServer starts the test Grin API server
func startTestGrinAPIServer(addr string) *http.Server {
	router := mux.NewRouter()
	addBlock()
	router.HandleFunc("/v1/blocks/{hash}", getBlock).Methods("GET")
	router.HandleFunc("/v1/status", getStatus).Methods("GET")
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}

// GetBlock displays a single block
func getBlock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range blocks {
		if item.Header.Hash == params["hash"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	height, _ := strconv.ParseUint(params["hash"], 10, 64)
	for _, item := range blocks {
		if item.Header.Height == height {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&api.BlockPrintable{})
}

// GetBlock displays a single block
func getStatus(w http.ResponseWriter, r *http.Request) {
	statusJSON, e := ioutil.ReadFile("test_data/status.json")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Fatal("File error")
		os.Exit(1)
	}
	var status api.Status
	if err := json.Unmarshal(statusJSON, &status); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error json unmarshal")
	}
	json.NewEncoder(w).Encode(status)
}

func addBlock() {
	jsonBlock16111, e := ioutil.ReadFile("test_data/block16111.json")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Fatal("File error")
		os.Exit(1)
	}
	var block16111 api.BlockPrintable
	if err := json.Unmarshal(jsonBlock16111, &block16111); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error json unmarshal")
	}
	blocks = append(blocks, block16111)
	jsonBlock16112, e := ioutil.ReadFile("test_data/block16112.json")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Fatal("File error")
		os.Exit(1)
	}
	var block16112 api.BlockPrintable
	if err := json.Unmarshal(jsonBlock16112, &block16112); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error json unmarshal")
	}
	blocks = append(blocks, block16112)
}

func nextAPI(increment int) (grinAPI, string) {
	var startPort = 23413
	portInt := startPort + increment
	port := strconv.Itoa(portInt)
	addr := "127.0.0.1:" + port
	return grinAPI{GrinServerAPI: addr}, addr
}

// The API used here

func TestGetBlockReward(t *testing.T) {
	grinAPI, addr := nextAPI(1)
	srv := startTestGrinAPIServer(addr)
	var blockHash = "0822cd711993d0f9a3ffdb4e755defdd4a40aa25ce72f8053fa330247a36f687"
	blockReward, err := grinAPI.GetBlockReward(blockHash)
	assert.NoError(t, err)
	var expectedBlockReward uint64 = 60013000000
	assert.Equal(t, expectedBlockReward, blockReward)
	srv.Shutdown(context.TODO())
}

func TestGetBlockRewardMissing(t *testing.T) {
	grinAPI, addr := nextAPI(2)
	srv := startTestGrinAPIServer(addr)
	var blockHash = "0822cd711993d0f9a3ffdb4e755defdd4a50aa25ce72f8053fa330247a36f687"
	blockReward, err := grinAPI.GetBlockReward(blockHash)
	assert.Error(t, err)
	assert.Equal(t, uint64(0), blockReward)
	srv.Shutdown(context.TODO())
}

func TestGetBlockByHash(t *testing.T) {
	grinAPI, addr := nextAPI(3)
	srv := startTestGrinAPIServer(addr)
	var blockHash = "0822cd711993d0f9a3ffdb4e755defdd4a40aa25ce72f8053fa330247a36f687"
	block, err := grinAPI.GetBlockByHash(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, blockHash, block.Header.Hash)
	srv.Shutdown(context.TODO())
}

func TestGetBlockByHashMissing(t *testing.T) {
	grinAPI, addr := nextAPI(4)
	srv := startTestGrinAPIServer(addr)
	var blockHash = "0822cd711993d0f9a3ffdb4e755defd84a40aa25ce72f8053fa330247a36f687"
	block, err := grinAPI.GetBlockByHash(blockHash)
	assert.Error(t, err)
	assert.Nil(t, block)
	srv.Shutdown(context.TODO())
}

func TestGetBlockByHashUnreachable(t *testing.T) {
	grinAPI := grinAPI{}
	var blockHash = "0822cd711993d0f9a3ffdb4e755defd84a40aa25ce72f8053fa330247a36f687"
	block, err := grinAPI.GetBlockByHash(blockHash)
	assert.Error(t, err)
	assert.Nil(t, block)
}

func TestGetBlockByHeight(t *testing.T) {
	grinAPI, addr := nextAPI(5)
	srv := startTestGrinAPIServer(addr)
	var blockHash = "0822cd711993d0f9a3ffdb4e755defdd4a40aa25ce72f8053fa330247a36f687"
	block, err := grinAPI.GetBlockByHeight(16111)
	assert.NoError(t, err)
	assert.Equal(t, blockHash, block.Header.Hash)
	srv.Shutdown(context.TODO())
}

func TestGetBlockByHeightMissing(t *testing.T) {
	grinAPI, addr := nextAPI(6)
	srv := startTestGrinAPIServer(addr)
	block, err := grinAPI.GetBlockByHeight(1619)
	assert.Error(t, err)
	assert.Nil(t, block)
	srv.Shutdown(context.TODO())
}

func TestGetBlockUnreachable(t *testing.T) {
	grinAPI := grinAPI{}
	block, err := grinAPI.GetBlockByHeight(1619)
	assert.Error(t, err)
	assert.Nil(t, block)
}

func TestGetStatus(t *testing.T) {
	grinAPI, addr := nextAPI(7)
	srv := startTestGrinAPIServer(addr)
	status, err := grinAPI.GetStatus()
	assert.NoError(t, err)
	assert.Equal(t, 1, status.ProtocolVersion)
	assert.Equal(t, "MW/Grin 0.4.0", status.UserAgent)
	var blockHash16112 = "044753a1faf8c6b5e01ba4a244cdc9f2d83587ce19b02e9d53a9aa9623708e37"
	assert.Equal(t, blockHash16112, status.Tip.LastBlockPushed)
	var blockHash16111 = "0822cd711993d0f9a3ffdb4e755defdd4a40aa25ce72f8053fa330247a36f687"
	assert.Equal(t, blockHash16111, status.Tip.PrevBlockToLast)
	srv.Shutdown(context.TODO())
}

func TestGetStatusUnreachable(t *testing.T) {
	grinAPI := grinAPI{}
	status, err := grinAPI.GetStatus()
	assert.Error(t, err)
	assert.Nil(t, status)
}

func TestGetTargetDifficultyAndHashratesMissingLastBlock(t *testing.T) {
	grinAPI, addr := nextAPI(9)
	srv := startTestGrinAPIServer(addr)
	status, err := grinAPI.GetStatus()
	assert.NoError(t, err)
	// update with fake block hash
	status.Tip.LastBlockPushed = "0822cd711993d0f9a3ffdb4e755defd84a40aa25ce72f8053fa330247a36f687"
	td, _, h, err := grinAPI.GetTargetDifficultyAndHashrates(status)
	assert.Error(t, err)
	expectedTD := uint64(0)
	assert.Equal(t, expectedTD, td)
	expectedHashrate := 0.0
	assert.Equal(t, expectedHashrate, h)
	srv.Shutdown(context.TODO())
}

func TestGetTargetDifficultyAndHashratesMissingPreviousBlock(t *testing.T) {
	grinAPI, addr := nextAPI(10)
	srv := startTestGrinAPIServer(addr)
	status, err := grinAPI.GetStatus()
	assert.NoError(t, err)
	status.Tip.PrevBlockToLast = "0822cd711993d0f9a3ffdb4e755defd84a40aa25ce72f8053fa330247a36f687"
	td, _, h, err := grinAPI.GetTargetDifficultyAndHashrates(status)
	assert.Error(t, err)
	expectedTD := uint64(0)
	assert.Equal(t, expectedTD, td)
	expectedHashrate := 0.0
	assert.Equal(t, expectedHashrate, h)
	srv.Shutdown(context.TODO())
}

func TestGetTargetDifficultyAndHashratesUnreachable(t *testing.T) {
	grinAPI := grinAPI{}
	status, err := grinAPI.GetStatus()
	assert.Error(t, err)
	td, _, h, err := grinAPI.GetTargetDifficultyAndHashrates(status)
	assert.Error(t, err)
	expectedTD := uint64(0)
	assert.Equal(t, expectedTD, td)
	expectedHashrate := 0.0
	assert.Equal(t, expectedHashrate, h)
}

func TestGetTargetDifficultyAndHashratesUnreachableAfter(t *testing.T) {
	grinAPI := grinAPI{}
	status := api.Status{ProtocolVersion: 1}
	td, _, h, err := grinAPI.GetTargetDifficultyAndHashrates(&status)
	assert.Error(t, err)
	expectedTD := uint64(0)
	assert.Equal(t, expectedTD, td)
	expectedHashrate := 0.0
	assert.Equal(t, expectedHashrate, h)
}
