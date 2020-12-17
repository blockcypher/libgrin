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

package client_test

import (
	"testing"
)

func TestNodeForeignAPI(t *testing.T) {
	// commenting this since this can't be done on CI for now
	/*
		url := "http://127.0.0.1:3413/v2/foreign"
		nodeOwnerAPI := client.NewNodeForeignAPI(url)
		// GetBlock
		{
			var height uint64 = 1006625
			block, err := nodeOwnerAPI.GetBlock(&height, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, block)
		}
		// GetHeader
		{
			var height uint64 = 1006625
			header, err := nodeOwnerAPI.GetHeader(&height, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, header)
		}
		// GetKernel
		{
			excess := "09b39b14c7dd1199622345ed005fd45cd3378296884f4cac6b94f5de5f0c2d53e6"
			kernel, err := nodeOwnerAPI.GetKernel(excess, nil, nil)
			assert.NoError(t, err)
			assert.NotNil(t, kernel)
		}
		// GetOutputs
		{
			commit := "08bebd4a9a0b217955c9385d8a3bb656e6c855a0a667ac03f041ea9534fff8ff64"
			commitArray := []string{commit}
			trueBool := true
			outputs, err := nodeOwnerAPI.GetOutputs(&commitArray, nil, nil, &trueBool, &trueBool)
			assert.NoError(t, err)
			assert.NotNil(t, outputs)
		}
		// GetPMMRIndices
		{
			var endBlockHeight uint64 = 100
			outputs, err := nodeOwnerAPI.GetPMMRIndices(0, &endBlockHeight)
			assert.NoError(t, err)
			assert.NotNil(t, outputs)
		}
		// GetPoolSize
		{
			poolSize, err := nodeOwnerAPI.GetPoolSize()
			assert.NoError(t, err)
			assert.NotNil(t, poolSize)
		}
		// GetStempoolSize
		{
			stempoolSize, err := nodeOwnerAPI.GetStempoolSize()
			assert.NoError(t, err)
			assert.NotNil(t, stempoolSize)
		}
		// GetTip
		{
			tip, err := nodeOwnerAPI.GetTip()
			assert.NoError(t, err)
			assert.NotNil(t, tip)
		}
		// GetUnconfirmedTransactions
		{
			unconfirmedTxs, err := nodeOwnerAPI.GetUnconfirmedTransactions()
			assert.NoError(t, err)
			assert.NotNil(t, unconfirmedTxs)
		}
		// GetUnspentOutputs
		{
			trueBool := true
			outputListing, err := nodeOwnerAPI.GetUnspentOutputs(1, nil, 2, &trueBool)
			assert.NoError(t, err)
			assert.NotNil(t, outputListing)
		}
		// GetVersion
		{
			version, err := nodeOwnerAPI.GetVersion()
			assert.NoError(t, err)
			assert.NotNil(t, version)
		}
	*/
}
