// Copyright 2020 BlockCypher
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
package pow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCuckooParams(t *testing.T) {
	cp := new(CuckooParams)
	params := cp.new(19, 42)
	assert.Equal(t, params.edgeBits, uint8(19))
	assert.Equal(t, params.proofSize, 42)
	assert.Equal(t, params.numEdges, uint64(524288))
	assert.Equal(t, params.siphashKeys, [4]uint64{0, 0, 0, 0})
}
