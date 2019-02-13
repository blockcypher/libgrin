// Copyright 2018 BlockCypher
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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/sethgrid/pester"
	log "github.com/sirupsen/logrus"
)

// PostSendJSON do a post request and populate a struct
func PostSendJSON(url string, params interface{}) ([]byte, error) {
	client := pester.New()
	// We don't to retry here
	client.MaxRetries = 0
	client.Timeout = 60 * time.Second

	jsonValue, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	r, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, nil
}

// PostSendJSONResponse do a post request and populate a struct
func PostSendJSONResponse(url string, params interface{}) (*http.Response, error) {
	client := pester.New()
	client.Timeout = 10 * time.Second
	jsonValue, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	r, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetJSON do a get request and populate a struct
func GetJSON(url string, target interface{}) error {
	client := pester.New()
	client.Timeout = 10 * time.Second
	client.KeepLog = true

	r, err := client.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"pester": client.LogString(),
			"error":  err,
			"url":    url,
		}).Error("HTTP: Pester log")
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return fmt.Errorf("error during getJSON. Status code: %d", r.StatusCode)
	}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&target); err != nil {
		log.WithFields(log.Fields{
			"status code": r.StatusCode,
			"error":       err,
			"url":         url,
		}).Error("HTTP: Error during decode")
		return err
	}
	return nil
}

// GetJSONLongTimeout do a get request and populate a struct with long timeout
func GetJSONLongTimeout(url string, target interface{}) error {
	client := pester.New()
	client.Timeout = 60 * time.Second
	client.KeepLog = true

	r, err := client.Get(url)
	if err != nil {
		log.WithFields(log.Fields{
			"pester": client.LogString(),
			"error":  err,
			"url":    url,
		}).Error("HTTP: Pester log")
		return err
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if r.StatusCode != 200 {
		return fmt.Errorf("error during getJSON. Status code: %d", r.StatusCode)
	}
	if err := dec.Decode(&target); err != nil {
		log.WithFields(log.Fields{
			"status code": r.StatusCode,
			"error":       err,
			"url":         url,
		}).Error("HTTP: Error during decode")
		return err
	}
	return nil
}

// CheckPort check if a port is open
func CheckPort(host string) bool {
	conn, err := net.DialTimeout("tcp", host, time.Duration(60)*time.Second)
	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		// Timeout
		return false
	}
	if err != nil {
		// Different error
		log.WithFields(log.Fields{
			"error": err,
			"host":  host,
		}).Error("HTTP: Check port error")
		return false
	}
	// Connection open closing it
	defer conn.Close()
	return true
}
