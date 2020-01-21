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

func TestSipash24(t *testing.T) {
	assert.Equal(t, SipHash24([4]uint64{1, 2, 3, 4}, 10, 21), uint64(928382149599306901))
	assert.Equal(t, SipHash24([4]uint64{1, 2, 3, 4}, 111, 21), uint64(10524991083049122233))
	assert.Equal(t, SipHash24([4]uint64{9, 7, 6, 7}, 12, 21), uint64(1305683875471634734))
	assert.Equal(t, SipHash24([4]uint64{9, 7, 6, 7}, 10, 21), uint64(11589833042187638814))
}

func TestSiphashBlock(t *testing.T) {
	assert.Equal(t, SipHashBlock([4]uint64{1, 2, 3, 4}, 10, 21, false), uint64(1182162244994096396))
	assert.Equal(t, SipHashBlock([4]uint64{1, 2, 3, 4}, 123, 21, false), uint64(11303676240481718781))
	assert.Equal(t, SipHashBlock([4]uint64{9, 7, 6, 7}, 12, 21, false), uint64(4886136884237259030))
}
