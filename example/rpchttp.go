// Copyright 2019 BlockCypher
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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/sethgrid/pester"
	log "github.com/sirupsen/logrus"
)

var requestCounter uint64

// RPCHTTPClient is a JSON-RPC over HTTP Client
type RPCHTTPClient struct {
	URL string
}

// Envelope is the JSON-RPC envelope
type Envelope struct {
	ID      JSONRPCID        `json:"id"`
	Version JSONRPCV2Version `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *rpcError        `json:"error,omitempty"`
}

// JSONRPCID represents the JSON-RPC V2 id
// will automatically be serialized
type JSONRPCID string

// MarshalJSON implement the Marshaler interface on JSONRPCVersion
func (e JSONRPCID) MarshalJSON() ([]byte, error) {
	counter := strconv.FormatUint(atomic.AddUint64(&requestCounter, 1), 10)
	b, err := json.Marshal(counter)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// rpcError represents a stratum error message
type rpcError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// JSONRPCV2Version represents the JSON-RPC V2 version string
// will always be serialized to "2.0"
type JSONRPCV2Version string

// Result is golang equivalent of Rust result
type Result struct {
	Ok  json.RawMessage
	Err json.RawMessage
}

// MarshalJSON implement the Marshaler interface on JSONRPCVersion
func (e JSONRPCV2Version) MarshalJSON() ([]byte, error) {
	version := "2.0"
	b, err := json.Marshal(version)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Request do a RPC POST request with the server
func (c *RPCHTTPClient) Request(method string, params json.RawMessage) (*Envelope, error) {
	requestBody, err := json.Marshal(Envelope{
		Method: method,
		Params: params,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't marshal RPC request")
		return nil, err
	}
	body, err := post(c.URL, requestBody)
	var envl Envelope
	if err := json.Unmarshal(body, &envl); err != nil {
		return nil, err
	}
	return &envl, nil
}

func post(url string, requestBody []byte) ([]byte, error) {
	client := pester.New()
	// We don't to retry here
	client.MaxRetries = 0
	client.Timeout = 60 * time.Second

	r, err := client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, nil
}
