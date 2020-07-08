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

//"github.com/blockcypher/libgrin/client"
//"github.com/blockcypher/libgrin/libwallet"
//"github.com/stretchr/testify/assert"

func TestReal(t *testing.T) {
	// commenting this since this can't be done on CI for now
	/*
		url := "http://127.0.0.1:3420/v3/owner"
		ownerAPI := example.NewSecureOwnerAPI(url)
		if err := ownerAPI.Init(); err != nil {
			assert.Error(t, err)
		}
		if err := ownerAPI.Open(nil, ""); err != nil {
			assert.Error(t, err)
		}
		// NodeHeight
		{
			nodeHeight, err := ownerAPI.NodeHeight()
			if err != nil {
				assert.Error(t, err)
			}
			assert.NotNil(t, nodeHeight)
		}
		// SetTorConfig
		{
			torConfig := libwallet.TorConfig{
				UseTorListener: true,
				SocksProxyAddr: "127.0.0.1:59050",
				SendConfigDir:  ".",
			}
			if err := ownerAPI.SetTorConfig(torConfig); err != nil {
				assert.Error(t, err)
			}
		}
		// GetSlatepackAddress
		{
			slatepackAddress, err := ownerAPI.GetSlatepackAddress(0)
			if err != nil {
				assert.Error(t, err)
			}
			assert.NotNil(t, slatepackAddress)
		}
		// GetSlatepackSecretKey
		{
			slatepackSecretKey, err := ownerAPI.GetSlatepackSecretKey(0)
			if err != nil {
				assert.Error(t, err)
			}
			assert.NotNil(t, slatepackSecretKey)
		}
		// CreateSlatepackMessage
		{
			id, err := uuid.Parse("0436430c-2b02-624c-2032-570501212b00")
			if err != nil {
				assert.Error(t, err)
			}
			slateV4 := slateversions.SlateV4{
				Ver: slateversions.VersionCompatInfoV4{
					Version:            4,
					BlockHeaderVersion: 2,
				},
				ID:  id,
				Sta: slateversions.Standard1SlateState,
				Off: "0gKWSQAAAADTApZJAAAAANQClkkAAAAA1QKWSQAAAAA",
				Amt: 6000000000,
				Fee: 8000000,
				Sigs: []slateversions.ParticipantDataV4{
					{
						Xs:    "Ajh4zoRXJ/Ok7HbKPz20s4otBdY2uMNjIQi4V/7WPJbe",
						Nonce: "AxuExVZ7EmRAmV0+1aq6BWXXHhg0YEgZ/5wX9enV3QeP",
					},
				},
			}

			var senderIndex uint32 = 0
			recipients := make([]string, 0)

			slatepackMessage, err := ownerAPI.CreateSlatepackMessage(0, slateV4, &senderIndex, recipients)
			if err != nil {
				assert.Error(t, err)
			}
			assert.NotNil(t, slatepackMessage)
		}
		// SlateFromSlatepackMessage
		{
			message := "BEGINSLATEPACK. 8GQrdcwdLKJD28F 3a9siP7ZhZgAh7w BR2EiZHza5WMWmZ Cc8zBUemrrYRjhq j3VBwA8vYnvXXKU BDmQBN2yKgmR8mX UzvXHezfznA61d7 qFZYChhz94vd8Ew NEPLz7jmcVN2C3w wrfHbeiLubYozP2 uhLouFiYRrbe3fQ 4uhWGfT3sQYXScT dAeo29EaZJpfauh j8VL5jsxST2SPHq nzXFC2w9yYVjt7D ju7GSgHEp5aHz9R xstGbHjbsb4JQod kYLuELta1ohUwDD pvjhyJmsbLcsPei k5AQhZsJ8RJGBtY bou6cU7tZeFJvor 4LB9CBfFB3pmVWD vSLd5RPS75dcnHP nbXD8mSDZ8hJS2Q A9wgvppWzuWztJ2 dLUU8f9tLJgsRBw YZAs71HiVeg7. ENDSLATEPACK."

			slateV4, err := ownerAPI.SlateFromSlatepackMessage(message, []uint32{0})
			if err != nil {
				assert.Error(t, err)
			}
			assert.NotNil(t, slateV4)
		}
		// DecodeSlatepackMessage
		{
			message := "BEGINSLATEPACK. t9EcGgrKr1GFCQB SK2jPCxME6Hgpqx bntpQm3zKFycoPY nW4UeoL4KQ7ExNK At6EQsvpz6MjUs8 6WG8KHEbMfqufJQ ZJTw2gkcdJmJjiJ f29oGgYqqXDZox4 ujPSjrtoxCN4h3e i1sZ8dYsm3dPeXL 7VQLsYNjAefciqj ZJXPm4Pqd7VDdd4 okGBGBu3YJvYzT6 arAxeCEx66us31h AJLcDweFwyWBkW5 J1DLiYAjt5ftFTo CjpfW9KjiLq2LM5 jepXWEHJPSDAYVK 4macDZUhRbJiG6E hrQcPrJBVC716mb Hw5E1PFrE6on5wq oEmrS4j9vaB5nw8 Z9ZyXvPc2LN7tER yt6pSHZeY9EpYdY zv4bthzfRfF8ePT TMeMpV2gpgyRXQa CPD2TR. ENDSLATEPACK."

			slatepack, err := ownerAPI.DecodeSlatepackMessage(message, []uint32{0})
			if err != nil {
				assert.Error(t, err)
			}
			assert.NotNil(t, slatepack)
		}

		if err := ownerAPI.Close(nil); err != nil {
			assert.Error(t, err)
		}
	*/
}
