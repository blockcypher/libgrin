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

	"github.com/blockcypher/libgrin/v5/client"
	"github.com/stretchr/testify/assert"
)

func TestNodeOwnerAPI(t *testing.T) {
	// commenting this since this can't be done on CI for now

	url := "http://127.0.0.1:3413/v2/owner"
	nodeOwnerAPI := client.NewNodeOwnerAPI(url)
	// GetStatus
	{
		status, err := nodeOwnerAPI.GetStatus()
		assert.NoError(t, err)
		assert.NotNil(t, status)
	}
	// ValidateChain (will timeout)
	{
		err := nodeOwnerAPI.ValidateChain()
		assert.NoError(t, err)
	}
	// CompactChain
	{
		err := nodeOwnerAPI.CompactChain()
		assert.NoError(t, err)
	}
}
