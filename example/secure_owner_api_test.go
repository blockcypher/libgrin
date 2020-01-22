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
	"testing"
)

//"github.com/blockcypher/libgrin/example"
//"github.com/blockcypher/libgrin/libwallet"
//"github.com/stretchr/testify/assert"

func TestReal(t *testing.T) {
	// commenting this since this can't be done on CI for now
	/*url := "http://127.0.0.1:3420/v3/owner"
	ownerAPI := example.NewSecureOwnerAPI(url)
	if err := ownerAPI.Init(); err != nil {
		assert.Error(t, err)
	}
	if err := ownerAPI.Open(nil, ""); err != nil {
		assert.Error(t, err)
	}
	nodeHeight, err := ownerAPI.NodeHeight()
	if err != nil {
		assert.Error(t, err)
	}
	assert.NotNil(t, nodeHeight)
	torConfig := libwallet.TorConfig{
		UseTorListener: true,
		SocksProxyAddr: "127.0.0.1:59050",
		SendConfigDir:  ".",
	}
	if err := ownerAPI.SetTorConfig(torConfig); err != nil {
		assert.Error(t, err)
	}
	if err := ownerAPI.Close(nil); err != nil {
		assert.Error(t, err)
	}*/
}
