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

package example_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/blockcypher/libgrin/example"
	"github.com/stretchr/testify/assert"
)

func TestEnvelope(t *testing.T) {
	requestBody, err := json.Marshal(example.Envelope{
		Method: "test",
		Params: nil,
	})
	assert.Nil(t, err)

	var envl example.Envelope
	err = json.Unmarshal(requestBody, &envl)
	fmt.Println(envl)
	assert.Nil(t, err)
	assert.Equal(t, example.JSONRPCID("1"), envl.ID)
	assert.Equal(t, "test", envl.Method)
	assert.Equal(t, example.JSONRPCV2Version("2.0"), envl.Version)
	assert.Nil(t, envl.Params)
}

/*
Skipping this test for the CI
func TestRetrieveSummaryInfoRaw(t *testing.T) {
	url := "http://127.0.0.1:3420/v2/owner"
	client := wallet.RPCHTTPClient{URL: url}
	params := []interface{}{true, 1}
	paramsBytes, err := json.Marshal(params)
	envl, err := client.Request("retrieve_summary_info", paramsBytes)
	assert.Nil(t, err)
	assert.Nil(t, envl.Error)
	var result wallet.Result
	err = json.Unmarshal(envl.Result, &result)
	assert.Equal(t, wallet.JSONRPCID("1"), envl.ID)
	assert.Equal(t, wallet.JSONRPCV2Version("2.0"), envl.Version)

	var okArray []json.RawMessage
	err = json.Unmarshal(result.Ok, &okArray)
	assert.Nil(t, err)
	var refreshed bool
	err = json.Unmarshal(okArray[0], &refreshed)
	assert.Nil(t, err)
	var walletInfo libwallet.WalletInfo
	err = json.Unmarshal(okArray[1], &walletInfo)
	assert.Nil(t, err)
}

*/
