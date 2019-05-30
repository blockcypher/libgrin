// Copyright 2019 BlockCypher
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

package wrappers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/blockcypher/libgrin/libwallet"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// startTestGrinAPIServer starts the test Grin API server
func startTestGrinOwnerAPIServer(addr string) {
	router := mux.NewRouter()
	router.HandleFunc("/v1/wallet/owner/retrieve_summary_info", handleGetSummaryInfo).Methods("GET")
	log.Fatal(http.ListenAndServe(addr, router))
}

// GetBlock displays a single block
func handleGetSummaryInfo(w http.ResponseWriter, r *http.Request) {
	walletinfoJSON, e := ioutil.ReadFile("test_data/summary_info.json")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Fatal("File error")
		os.Exit(1)
	}
	var walletInfo interface{}
	if err := json.Unmarshal(walletinfoJSON, &walletInfo); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error json unmarshal")
	}
	json.NewEncoder(w).Encode(walletInfo)
}

func TestGetAmountCurrentlySpendable(t *testing.T) {
	go startTestGrinOwnerAPIServer("127.0.0.1:3420")
	amountSpendable, err := GetAmountCurrentlySpendable(true)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5100000000000), amountSpendable)
	amountSpendable, err = GetAmountCurrentlySpendable(false)
	assert.NoError(t, err)
	assert.Equal(t, uint64(5100000000000), amountSpendable)
}

func TestUnmarshallingSlate(t *testing.T) {
	slateJSON, err := ioutil.ReadFile("test_data/slate.json")
	assert.NoError(t, err)
	var slate libwallet.Slate
	err = json.Unmarshal(slateJSON, &slate)
	assert.NoError(t, err)
}

func TestUnmarshallingPartialSlate(t *testing.T) {
	slateJSON, err := ioutil.ReadFile("test_data/partial_slate1.json")
	assert.NoError(t, err)
	var slate libwallet.Slate
	err = json.Unmarshal(slateJSON, &slate)
	assert.NoError(t, err)
}
