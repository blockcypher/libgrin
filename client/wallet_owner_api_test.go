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

func TestWalletOwnerAPI(t *testing.T) {
	// commenting this since this can't be done on CI for now

	/*
		url := "http://127.0.0.1:3420/v3/owner"
		ownerAPI := client.NewSecureOwnerAPI(url)
		if err := ownerAPI.Init(); err != nil {
			assert.NoError(t, err)
		}
		if err := ownerAPI.Open(nil, ""); err != nil {
			assert.NoError(t, err)
		}
		// NodeHeight
		{
			nodeHeight, err := ownerAPI.NodeHeight()
			assert.NoError(t, err)
			assert.NotNil(t, nodeHeight)
		}
		// SetTorConfig
		{
			torConfig := libwallet.TorConfig{
				UseTorListener: true,
				SocksProxyAddr: "127.0.0.1:59050",
				SendConfigDir:  "/Users/quentin/.grin/floo/tor/sender",
			}
			if err := ownerAPI.SetTorConfig(torConfig); err != nil {
				assert.Error(t, err)
			}
		}
		// Get a transaction by txid and get output
		{
			txSlateID := uuid.MustParse("2bd40747-366f-42b0-a798-dab09133d648")
			_, txLog, err := ownerAPI.RetrieveTxs(true, nil, &txSlateID)
			assert.NotEmpty(t, txLog)
			assert.NoError(t, err)
			txLogOwned := *txLog
			_, outputsRetrieved, err := ownerAPI.RetrieveOutputs(false, false, &txLogOwned[0].ID)
			assert.NoError(t, err)
			assert.NotEmpty(t, outputsRetrieved)
			slate, err := ownerAPI.GetStoredTx(nil, &txSlateID)
			assert.NoError(t, err)
			for _, comm := range *slate.Coms {
				if comm.P == nil {
					fmt.Println("found input", comm.C)
				}
			}
			spew.Dump(slate)
			assert.NotEmpty(t, slate)
			assert.True(t, false)
		}
		// GetSlatepackAddress
		{
			slatepackAddress, err := ownerAPI.GetSlatepackAddress(0)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.NotNil(t, slatepackAddress)
		}
		// GetSlatepackSecretKey
		{
			slatepackSecretKey, err := ownerAPI.GetSlatepackSecretKey(0)
			if err != nil {
				assert.NoError(t, err)
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
				Off: "d202964900000000d302964900000000d402964900000000d502964900000000",
				Amt: 6000000000,
				Fee: 7000000,
				Sigs: []slateversions.ParticipantDataV4{
					{
						Xs:    "030152d2d72e2dba7c6086ad49a219d9ff0dfe0fd993dcaea22e058c210033ce93",
						Nonce: "031b84c5567b126440995d3ed5aaba0565d71e1834604819ff9c17f5e9d5dd078f",
					},
				},
			}

			var senderIndex uint32 = 0
			recipients := make([]string, 0)

			slatepackMessage, err := ownerAPI.CreateSlatepackMessage(0, slateV4, &senderIndex, recipients)
			if err != nil {
				assert.NoError(t, err)
			}
			assert.NotNil(t, slatepackMessage)
		}
		// SlateFromSlatepackMessage
		{
			message := "BEGINSLATEPACK. 8GQrdcwdLKJD28F 3a9siP7ZhZgAh7w BR2EiZHza5WMWmZ Cc8zBUemrrYRjhq j3VBwA8vYnvXXKU BDmQBN2yKgmR8mX UzvXHezfznA61d7 qFZYChhz94vd8Ew NEPLz7jmcVN2C3w wrfHbeiLubYozP2 uhLouFiYRrbe3fQ 4uhWGfT3sQYXScT dAeo29EaZJpfauh j8VL5jsxST2SPHq nzXFC2w9yYVjt7D ju7GSgHEp5aHz9R xstGbHjbsb4JQod kYLuELta1ohUwDD pvjhyJmsbLcsPei k5AQhZsJ8RJGBtY bou6cU7tZeFJvor 4LB9CBfFB3pmVWD vSLd5RPS75dcnHP nbXD8mSDZ8hJS2Q A9wgvppWzuWztJ2 dLUU8f9tLJgsRBw YZAs71HiVeg7. ENDSLATEPACK."

			slateV4, err := ownerAPI.SlateFromSlatepackMessage(message, []uint32{0})
			if err != nil {
				assert.NoError(t, err)
			}
			assert.NotNil(t, slateV4)
		}
		// DecodeSlatepackMessage
		{
			message := "BEGINSLATEPACK. t9EcGgrKr1GFCQB SK2jPCxME6Hgpqx bntpQm3zKFycoPY nW4UeoL4KQ7ExNK At6EQsvpz6MjUs8 6WG8KHEbMfqufJQ ZJTw2gkcdJmJjiJ f29oGgYqqXDZox4 ujPSjrtoxCN4h3e i1sZ8dYsm3dPeXL 7VQLsYNjAefciqj ZJXPm4Pqd7VDdd4 okGBGBu3YJvYzT6 arAxeCEx66us31h AJLcDweFwyWBkW5 J1DLiYAjt5ftFTo CjpfW9KjiLq2LM5 jepXWEHJPSDAYVK 4macDZUhRbJiG6E hrQcPrJBVC716mb Hw5E1PFrE6on5wq oEmrS4j9vaB5nw8 Z9ZyXvPc2LN7tER yt6pSHZeY9EpYdY zv4bthzfRfF8ePT TMeMpV2gpgyRXQa CPD2TR. ENDSLATEPACK."

			slatepack, err := ownerAPI.DecodeSlatepackMessage(message, []uint32{0})
			if err != nil {
				assert.NoError(t, err)
			}
			assert.NotNil(t, slatepack)
		}
		// Direct send to slatepack, also works with http
		{
			// Send the payment proof only for Tor Address
			dest := "tgrin1m7vmg0z5tu373haxtglape4hv45lw86xscq93ur7jj5rmpk5sw7qf2xf35"
			/*paymentProofRecipientAddress, err := slatepack.NewSlatepackAddressFromString(dest)
			if err != nil {
				assert.NoError(t, err)
			}

			sendArgs := libwallet.InitTxSendArgs{
				Dest:    dest,
				PostTx:  true,
				Fluff:   true,
				SkipTor: false,
			}
			args := libwallet.InitTxArgs{
				SrcAcctName:                  nil,
				Amount:                       1000000000,
				MinimumConfirmations:         2,
				MaxOutputs:                   500,
				NumChangeOutputs:             1,
				SelectionStrategyIsUseAll:    false,
				TargetSlateVersion:           nil,
				TTLBlocks:                    nil,
				PaymentProofRecipientAddress: nil,
				EstimateOnly:                 nil,
				SendArgs:                     &sendArgs,
			}
			slateV4, err := ownerAPI.InitSendTx(args)
			assert.NoError(t, err)
			assert.Equal(t, slateV4.Sta, slateversions.Standard3SlateState)
		}
		// Async send to file
		{
			args := libwallet.InitTxArgs{
				SrcAcctName:                  nil,
				Amount:                       1000000000,
				MinimumConfirmations:         2,
				MaxOutputs:                   500,
				NumChangeOutputs:             1,
				SelectionStrategyIsUseAll:    false,
				TargetSlateVersion:           nil,
				TTLBlocks:                    nil,
				PaymentProofRecipientAddress: nil,
				EstimateOnly:                 nil,
				SendArgs:                     nil,
			}
			slateV4, err := ownerAPI.InitSendTx(args)
			assert.NoError(t, err)
			assert.Equal(t, slateV4.Sta, slateversions.Standard1SlateState)
			spew.Dump(slateV4)
			assert.True(t, false)
			// create slatepack message
			var senderIndex uint32 = 0
			recipients := make([]string, 0)

			slatepackMessage, err := ownerAPI.CreateSlatepackMessage(0, *slateV4, &senderIndex, recipients)
			if err != nil {
				assert.NoError(t, err)
			}
			spew.Dump(slatepackMessage)
			assert.NotNil(t, slatepackMessage)
			// lock outputs
			if err := ownerAPI.TxLockOutputs(*slateV4); err != nil {
				assert.NoError(t, err)
			}
		}

		// finalize tx given a slatepack
		{
			message := "BEGINSLATEPACK. afMATt9TKZ4SYEU dXsGHPprrMr6yQb swGDFACw6VNWysN H128knQQ6EBGCsT jQY9KaSqk4A49Pp sWzsp3DRfveRHZ6 dFcm7MCv5mDeP56 q26eQ7WpLNpn9VS tL9hERwY4D6pgZS BuuqFdvhR5GCYLw w1ttb6jFQffr92a EWaDCqXhVzUnsj5 a8DQFZrJCavDBV6 vJsxtSw31NArLCi ddj4tNQgLy7GkTU oUWU9XabeMNkXpH kM89RUpWic12kaV cNZHt2r1XFsfFBf z3HYSm9ePd4mptA 2CFdQHKm4aCf4dt SN2FmPSLT8YM2w2 FDiPoirwYkZn1d4 txDSLoz3CgoS3d1 kqhrovhsDfn28Gj Urbdu9f48wDoZbf LDrkg6p6xu8cAFa CuWyNaL2TSpW469 F2LX6QAT3rFS1Aa DDLPswGdF9E1cFR CMZHjkmSMoEauys 1zCEJYuBWtoDq9S TAyWtdjs5RNw2gi 1YJtmwJf29xh1HN QJMCoZjudGN179n CAQLjHHzS5Pt7ap cAyaBQG5sKsXDh4 r2eUDFuLg4Apn6J NzMrmUkjaGHYpj7 UZpLNY15ef4xWsa 38qJtDJCgtUH4Rq QQ7i62C8upMuS9S XzpJgJFDDfcGbta sTcCmBwEDfiit3h rLmhDi8ki916D4T joWNKbSbJzcZ4be vRAiLD2dv9pDZxp vfTL2hfrmoQ9M4E QMWfzzbCAz4Qwt1 ucTQqPPWpnNyaYM 1mah3qcgphnWUSL 6NZiRTzuVN2ysRV gzvt484wC4Am2iw M4ViwCMz9nNQqA5 EGjxJXV6bf2pvJE zo8s8KSnvMsCYQq DenNmNH4XWY4BN8 RtRCHuywj3Pr14c JeUiGuu1yFuoVQ9 GSoTpzEpw56j4Tb e7P3Qjrbrr6eGcS n7stPopBBm8nbdi wMmQa2vPvk4RX5o DgZaX6CqT729cKw J19MoBQoga6oA8w Q1T7ZpcpzJz6rAg GwmWPzXaWUXkJkH iCxs9XvVYtDQybc NV2KVrHgKfHHQnK QD3cU1KHSu4qhXM 3zdUbxXcSwAr4qP huPCfPgm8XiYRsV UwWDAfMoFh5aFi5 3VFGCK7Ly34k2CJ xRDvKAwjTeUARdd SzoexW1WAxDqMvo cifEy3BFTbQmvxS 7ER5uxDKfQM2ioZ sdZMLwCX12t7pMi jfGzqWzCZcejLFe 9ne3aLF4kAi38hq Bsc8C5mWBPewFxH h8hUttHGKwJ1am5 xXtgSW8obySQiBq nJpWxaTbjFyBLr7 bX1XpPh4WbUPGeq 584HSXyx4x37S22 t1NazBm4Tys9VHh nKLaTA2pDPJeVHK ZK5xocG4X2cVtsM g652gchnQUNWfw7 sK5uMRRpphKkdX8 wgTGSX5FSW5k4cr oW67ZbGKa6EC1tE zPiKpTFkjDJJ7VU vUHcS4P2B4KQhNX csVKu1fatpaMu3E EXEdTAxM8G6ySsZ dwr93LwGAqBH3zD 4QKseK2xbAuNuR8 J1wMVFYTFqSG3Fe fntwo81wD6mGyGE 1aXA2UvUyUnydLj 6r6fVjEvEFnDCVd srAJPLV3vpa5aYz HFCGyqxhFcepMdw imN3VnhryExT1qZ q1RT4gSJ6Ku2YHG W4Xzrifxbp8SLdY E4P8SXxotrGe5Hr hETY97LKUYe9tme zWAtJfPPeGjxhQk GgxP1U1ZVEjCpcr 4rq1nMGFVdG34BV A4DYMNFY4L8haZe Sr6PvWSWYPuPTmi 26MDcTG5. ENDSLATEPACK."
			// no need to decode because the slatepack was not encrypted
			receivedSlateV4, err := ownerAPI.SlateFromSlatepackMessage(message, []uint32{0})
			assert.NoError(t, err)

			assert.NotNil(t, receivedSlateV4)
			spew.Dump(receivedSlateV4)
			assert.Equal(t, receivedSlateV4.Sta, slateversions.Standard2SlateState)
			// finalize it
			slateV4, err := ownerAPI.FinalizeTx(*receivedSlateV4)
			assert.NoError(t, err)
			// post it
			ownerAPI.PostTx(*slateV4, true)

			assert.True(t, false)
		}

		if err := ownerAPI.Close(nil); err != nil {
			assert.Error(t, err)
		}
	*/
}
