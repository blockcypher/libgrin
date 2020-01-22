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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
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

// MarshalJSON implement the Marshaler interface on JSONRPCID
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

// EncryptedData are the params or result to send/receive for the encrypted owner API
type EncryptedData struct {
	Nonce   string `json:"nonce"`
	BodyEnc string `json:"body_enc"`
}

// EncryptedRequest do an encrypted RPC POST request with the server
func (c *RPCHTTPClient) EncryptedRequest(method string, params json.RawMessage, sharedSecret []byte) (*Envelope, error) {
	toEncryptRequestBody, err := json.Marshal(Envelope{
		Method: method,
		Params: params,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't marshal RPC request body to encrypt")
		return nil, err
	}
	nonce := make([]byte, 12)
	rand.Read(nonce)
	encryptedRequestBody, err := encrypt(sharedSecret, nonce, toEncryptRequestBody)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't encrypt request body")
		return nil, err
	}

	encryptedParams, err := json.Marshal(EncryptedData{
		Nonce:   hex.EncodeToString(nonce),
		BodyEnc: base64.StdEncoding.EncodeToString(encryptedRequestBody),
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't marshal encrypted RPC params")
		return nil, err
	}
	requestBody, err := json.Marshal(Envelope{
		Method: "encrypted_request_v3",
		Params: encryptedParams,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Couldn't marshal RPC request body")
		return nil, err
	}
	body, err := post(c.URL, requestBody)
	var envl Envelope
	if err := json.Unmarshal(body, &envl); err != nil {
		return nil, err
	}
	if envl.Error != nil {
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var encryptedOk EncryptedData
	if err = json.Unmarshal(result.Ok, &encryptedOk); err != nil {
		return nil, err
	}
	nonceResponse, err := hex.DecodeString(encryptedOk.Nonce)
	if err != nil {
		return nil, err
	}
	encryptedBody, err := base64.StdEncoding.DecodeString(encryptedOk.BodyEnc)
	if err != nil {
		return nil, err
	}
	decryptedBody, err := decrypt(sharedSecret, nonceResponse, encryptedBody)
	var envlDecrypted Envelope
	if err := json.Unmarshal(decryptedBody, &envlDecrypted); err != nil {
		return nil, err
	}
	return &envlDecrypted, nil
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
		return nil, err
	}
	return responseData, nil
}

func encrypt(key, nonce, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nil
}

func decrypt(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
