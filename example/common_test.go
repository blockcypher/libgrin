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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodTypeToString(t *testing.T) {
	httpPayout := HTTPPayoutMethod
	stringHTTPPayout := httpPayout.methodTypeToString()
	assert.Equal(t, "http", stringHTTPPayout)
	filePayout := FilePayoutMethod
	stringFilePayout := filePayout.methodTypeToString()
	assert.Equal(t, "file", stringFilePayout)
	invalidPayout := InvalidPayoutMethod
	stringInvalidPayout := invalidPayout.methodTypeToString()
	assert.Equal(t, "http", stringInvalidPayout)
}

func TestStringToMethodType(t *testing.T) {
	httpPayout, err := stringToMethodType("http")
	assert.NoError(t, err)
	assert.Equal(t, HTTPPayoutMethod, httpPayout)
	filePayout, err := stringToMethodType("file")
	assert.NoError(t, err)
	assert.Equal(t, FilePayoutMethod, filePayout)
	invalidPayout, err := stringToMethodType("tor")
	assert.Error(t, err)
	assert.Equal(t, InvalidPayoutMethod, invalidPayout)
}
