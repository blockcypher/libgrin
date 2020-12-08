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

package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphWeight(t *testing.T) {
	// initial weights
	assert.Equal(t, GraphWeight(Mainnet, 1, 31), uint64(256*31))
	assert.Equal(t, GraphWeight(Mainnet, 1, 32), uint64(512*32))
	assert.Equal(t, GraphWeight(Mainnet, 1, 33), uint64(1024*33))

	// one year in, 31 starts going down, the rest stays the same
	assert.Equal(t, GraphWeight(Mainnet, YearHeight, 31), uint64(256*30))
	assert.Equal(t, GraphWeight(Mainnet, YearHeight, 32), uint64(512*32))
	assert.Equal(t, GraphWeight(Mainnet, YearHeight, 33), uint64(1024*33))

	// 31 loses one factor per week
	assert.Equal(t, GraphWeight(Mainnet, YearHeight+WeekHeight, 31), uint64(256*29))
	assert.Equal(t, GraphWeight(Mainnet, YearHeight+2*WeekHeight, 31), uint64(256*28))
	assert.Equal(t, GraphWeight(Mainnet, YearHeight+32*WeekHeight, 31), uint64(0))

	// 2 years in, 31 still at 0, 32 starts decreasing
	assert.Equal(t, GraphWeight(Mainnet, 2*YearHeight, 31), uint64(0))
	assert.Equal(t, GraphWeight(Mainnet, 2*YearHeight, 32), uint64(512*32))
	assert.Equal(t, GraphWeight(Mainnet, 2*YearHeight, 33), uint64(1024*33))

	// 32 phaseout on hold
	assert.Equal(t, GraphWeight(Mainnet, 2*YearHeight+WeekHeight, 32), uint64(512*32))
	assert.Equal(t, GraphWeight(Mainnet, 2*YearHeight+WeekHeight, 31), uint64(0))
	assert.Equal(t, GraphWeight(Mainnet, 2*YearHeight+30*WeekHeight, 32), uint64(512*32))
	assert.Equal(t,
		GraphWeight(Mainnet, 2*YearHeight+31*WeekHeight, 32), uint64(512*32))

	// 3 years in, nothing changes
	assert.Equal(t, GraphWeight(Mainnet, 3*YearHeight, 31), uint64(0))
	assert.Equal(t, GraphWeight(Mainnet, 3*YearHeight, 32), uint64(512*32))
	assert.Equal(t, GraphWeight(Mainnet, 3*YearHeight, 33), uint64(1024*33))

	// 4 years in, still on hold
	assert.Equal(t, GraphWeight(Mainnet, 4*YearHeight, 31), uint64(0))
	assert.Equal(t, GraphWeight(Mainnet, 4*YearHeight, 32), uint64(512*32))
	assert.Equal(t, GraphWeight(Mainnet, 4*YearHeight, 33), uint64(1024*33))
}

func TestSecondaryPoWRatio(t *testing.T) {
	// Tests for mainnet chain type.
	assert.Equal(t, SecondaryPoWRatio(1), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(89), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(90), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(91), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(179), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(180), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(181), uint64(90))

	oneWeek := uint64(60 * 24 * 7)
	assert.Equal(t, SecondaryPoWRatio(oneWeek-1), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(oneWeek), uint64(90))
	assert.Equal(t, SecondaryPoWRatio(oneWeek+1), uint64(90))

	twoWeeks := oneWeek * 2
	assert.Equal(t, SecondaryPoWRatio(twoWeeks-1), uint64(89))
	assert.Equal(t, SecondaryPoWRatio(twoWeeks), uint64(89))
	assert.Equal(t, SecondaryPoWRatio(twoWeeks+1), uint64(89))

	t4ForkHeight := uint64(64000)
	assert.Equal(t, SecondaryPoWRatio(t4ForkHeight-1), uint64(85))
	assert.Equal(t, SecondaryPoWRatio(t4ForkHeight), uint64(85))
	assert.Equal(t, SecondaryPoWRatio(t4ForkHeight+1), uint64(85))

	oneYear := oneWeek * 52
	assert.Equal(t, SecondaryPoWRatio(oneYear), uint64(45))

	ninetyOneWeeks := oneWeek * 91
	assert.Equal(t, SecondaryPoWRatio(ninetyOneWeeks-1), uint64(12))
	assert.Equal(t, SecondaryPoWRatio(ninetyOneWeeks), uint64(12))
	assert.Equal(t, SecondaryPoWRatio(ninetyOneWeeks+1), uint64(12))

	twoYears := oneYear * 2
	assert.Equal(t, SecondaryPoWRatio(twoYears-1), uint64(1))
	assert.Equal(t, SecondaryPoWRatio(twoYears), uint64(0))
	assert.Equal(t, SecondaryPoWRatio(twoYears+1), uint64(0))
}

func TestHardForks(t *testing.T) {
	// Tests for Mainnet
	{
		assert.True(t, ValidHeaderVersion(Mainnet, 0, 1))
		assert.True(t, ValidHeaderVersion(Mainnet, 10, 1))

		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval-1, 1))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval, 2))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval+1, 2))

		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*2-1, 2))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*2, 3))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*2+1, 3))

		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*3-1, 3))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*3, 4))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*3+1, 4))

		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*4-1, 4))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*4, 5))
		assert.True(t, ValidHeaderVersion(Mainnet, HardForkInterval*4+1, 5))
	}
	// Tests for Testnet
	{
		assert.True(t, ValidHeaderVersion(Testnet, 0, 1))
		assert.True(t, ValidHeaderVersion(Testnet, 10, 1))

		assert.True(t, ValidHeaderVersion(Testnet, TestnetFirstHardFork-1, 1))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetFirstHardFork, 2))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetFirstHardFork+1, 2))

		assert.True(t, ValidHeaderVersion(Testnet, TestnetSecondHardFork-1, 2))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetSecondHardFork, 3))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetSecondHardFork+1, 3))

		assert.True(t, ValidHeaderVersion(Testnet, TestnetThirdHardFork-1, 3))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetThirdHardFork, 4))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetThirdHardFork+1, 4))

		assert.True(t, ValidHeaderVersion(Testnet, TestnetFourthHardFork-1, 4))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetFourthHardFork, 5))
		assert.True(t, ValidHeaderVersion(Testnet, TestnetFourthHardFork+1, 5))
	}
}
