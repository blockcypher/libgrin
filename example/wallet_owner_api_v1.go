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

package example

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/blockcypher/libgrin/libwallet"
	log "github.com/sirupsen/logrus"
)

// IssueSendTx issues a send transaction
func IssueSendTx(ip string, method PayoutMethodType, amount uint64) (*libwallet.Slate, error) {
	// Wallet Owner IP is binded to localhost
	url := "http://127.0.0.1:3420/v1/wallet/owner/issue_send_tx"
	dest := "file"
	if method == HTTPPayoutMethod {
		// Add http only on empty prefix
		if !strings.HasPrefix(ip, "https://") && !strings.HasPrefix(ip, "http://") {
			dest = "http://" + ip
		} else {
			dest = ip
		}
	}
	log.WithFields(log.Fields{
		"method": method,
		"amount": amount,
		"dest":   dest,
	}).Info("GRINAPI: Sending amount")
	sendParams := libwallet.SendTXArgs{
		Amount:                    amount,
		MinimumConfirmations:      1,
		Method:                    method.methodTypeToString(),
		Dest:                      dest,
		MaxOutputs:                500,
		NumChangeOutputs:          1,
		SelectionStrategyIsUseAll: false,
	}

	response, err := PostSendJSON(url, sendParams)
	if err != nil {
		return nil, err
	}
	responseString := string(response)
	if strings.Contains(strings.ToLower(responseString), "error") {
		log.WithFields(log.Fields{
			"response": response,
		}).Error("GRINAPI: wallet error response")
		return nil, errors.New(responseString)
	}
	if strings.Contains(strings.ToLower(responseString), "not enough funds") {
		log.WithFields(log.Fields{
			"response": response,
		}).Error("GRINAPI: wallet error response")
		return nil, errors.New(responseString)
	}
	// Unmarshalling to slate
	var slate libwallet.Slate
	if err := json.Unmarshal(response, &slate); err != nil {
		// If we could not unmarshal here the transaction is still sent so do not return error
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GRINAPI: Could not unmarshal slate")
		return nil, err
	}
	return &slate, nil
}

// CancelTxByID cancels a transaction by id
func CancelTxByID(id uint) error {
	log.WithFields(log.Fields{
		"id": id,
	}).Info("GRINAPI: Cancelling transaction")
	// Wallet Owner IP is binded to localhost
	url := "http://127.0.0.1:3420/v1/wallet/owner/cancel_tx?id=" + string(id)
	_, err := PostSendJSON(url, nil)
	if err != nil {
		return err
	}
	return nil
}

// CancelTxByTxID cancels a transaction by txid
func CancelTxByTxID(txID string) error {
	log.WithFields(log.Fields{
		"txID": txID,
	}).Info("GRINAPI: Cancelling transaction")
	// Wallet Owner IP is binded to localhost
	url := "http://127.0.0.1:3420/v1/wallet/owner/cancel_tx?tx_id=" + txID
	_, err := PostSendJSON(url, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetAmountCurrentlySpendable retrieves the amount currently spendable
func GetAmountCurrentlySpendable(refresh bool) (uint64, error) {
	log.Info("GRINAPI: Getting amount currently spendable")
	// Wallet Owner IP is binded to localhost
	var url string
	if refresh {
		url = "http://127.0.0.1:3420/v1/wallet/owner/retrieve_summary_info?refresh"
	} else {
		url = "http://127.0.0.1:3420/v1/wallet/owner/retrieve_summary_info"
	}

	var resInterface interface{}
	if err := GetJSONLongTimeout(url, &resInterface); err != nil {
		return 0, err
	}
	walletInfoInterface, ok := resInterface.([]interface{})
	if !ok {
		log.Error("GRINAPI: Could not cast wallet info to interface")
		return 0, fmt.Errorf("could not cast wallet info to interface")
	}
	if len(walletInfoInterface) != 2 {
		log.Error("GRINAPI: Wallet info is not the proper length")
		return 0, fmt.Errorf("wallet info is not the proper length")
	}
	walletInfoMap, ok := walletInfoInterface[1].(map[string]interface{})
	if !ok {
		log.Error("GRINAPI: Could not cast wallet info to map")
		return 0, fmt.Errorf("could not cast wallet info to map")
	}

	amountCurrentlySpendableFloat := walletInfoMap["amount_currently_spendable"].(float64)
	if !ok {
		log.Error("GRINAPI: Could not cast amount to float64")
		return 0, fmt.Errorf("could not cast amount to float64")
	}
	amountCurrentlySpendable := uint64(amountCurrentlySpendableFloat)
	log.WithFields(log.Fields{
		"amount": amountCurrentlySpendable,
	}).Info("GRINAPI: Wallet balance")

	return amountCurrentlySpendable, nil
}

// PostTransaction Post a final slate to the Grin Node
func PostTransaction(slate *libwallet.Slate) error {
	log.WithFields(log.Fields{
		"txid":   slate.ID,
		"amount": slate.Amount,
	}).Info("GRINAPI: Posting transaction")
	// Wallet Owner IP is binded to localhost
	url := "http://127.0.0.1:3420/v1/wallet/owner/post_tx?fluff"
	response, err := PostSendJSONResponse(url, &slate)
	if err != nil {
		return err
	}
	// Check if the post went well
	if response.StatusCode != 200 {
		log.WithFields(log.Fields{
			"txid":        slate.ID,
			"amount":      slate.Amount,
			"status code": response.StatusCode,
		}).Error("GRINAPI: Error posting transaction")
		return fmt.Errorf("error posting transaction")
	}
	return nil
}

// FinalizeTransaction finalize a transaction
func FinalizeTransaction(slate *libwallet.Slate) (*libwallet.Slate, error) {
	log.WithFields(log.Fields{
		"txid":   slate.ID,
		"amount": slate.Amount,
	}).Info("GRINAPI: Finalizing transaction")
	// Wallet Owner IP is binded to localhost
	url := "http://127.0.0.1:3420/v1/wallet/owner/finalize_tx"
	response, err := PostSendJSON(url, &slate)
	if err != nil {
		return nil, err
	}
	var finalizedSlate libwallet.Slate
	if err := json.Unmarshal(response, &finalizedSlate); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GRINAPI: Could not unmarshal slate")
		return nil, err
	}
	return &finalizedSlate, nil
}

// FinalizeTransactionString same as above but with a json string
func FinalizeTransactionString(slateString string) (*libwallet.Slate, error) {
	var slate libwallet.Slate
	if err := json.Unmarshal([]byte(slateString), &slate); err != nil {
		// If we could not unmarshal here the transaction is still sent so do not return error
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GRINAPI: Could not unmarshal slate")
		return nil, err
	}
	log.WithFields(log.Fields{
		"txid":   slate.ID,
		"amount": slate.Amount,
	}).Info("GRINAPI: Finalizing transaction")
	// Wallet Owner IP is binded to localhost
	url := "http://127.0.0.1:3420/v1/wallet/owner/finalize_tx"
	response, err := PostSendJSON(url, &slate)
	if err != nil {
		return nil, err
	}
	var finalizedSlate libwallet.Slate
	if err := json.Unmarshal(response, &finalizedSlate); err != nil {
		// If we could not unmarshal here the transaction is still sent so do not return error
		log.WithFields(log.Fields{
			"error": err,
		}).Error("GRINAPI: Could not unmarshal slate. Probably a wallet error.")
		return nil, err
	}
	return &finalizedSlate, nil
}
