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

import "math"

// GrinBase is the grin base. A grin is divisible to 10^9, following the SI prefixes
const GrinBase uint64 = 1000000000

// MilliGrin is a thousand of a grin
const MilliGrin uint64 = GrinBase / 1000

// MicroGrin is a thousand of a milligrin
const MicroGrin uint64 = MilliGrin / 1000

// NanoGrin is the smallest unit, takes a billion to make a grin
const NanoGrin uint64 = 1

// BlockTimeSec is the  block interval, in seconds, the network will tune its next_target for. Note
// that we may reduce this value in the future as we get more data on mining
// with Cuckoo Cycle, networks improve and block propagation is optimized
// (adjusting the reward accordingly).
const BlockTimeSec uint64 = 60

// Reward is the block subsidy amount, one grin per second on average
const Reward uint64 = BlockTimeSec * GrinBase

// Actual block reward for a given total fee amount
func reward(fee uint64) uint64 {
	return saturatingAddUint64(Reward, fee)
}

// HourHeight is the nominal height for standard time intervals, hour is 60 blocks
const HourHeight uint64 = 3600 / BlockTimeSec

// DayHeight is 1440 blocks
const DayHeight uint64 = 24 * HourHeight

// WeekHeight is 10 080 blocks
const WeekHeight uint64 = 7 * DayHeight

// YearHeight is 524 160 blocks
const YearHeight uint64 = 52 * WeekHeight

// CoinbaseMaturity is the number of blocks before a coinbase matures and can be spent
const CoinbaseMaturity uint64 = DayHeight

// SecondaryPoWRatio is the ratio the secondary proof of work should take over the primary, as a
// function of block height (time). Starts at 90% losing a percent
// approximately every week. Represented as an integer between 0 and 100.
func SecondaryPoWRatio(height uint64) uint64 {
	return saturatingSubUint64(90, height/(2*YearHeight/90))
}

// ProofSize is the Cuckoo-cycle proof size (cycle length)
const ProofSize int = 42

// DefaultMinEdgeBits is the default Cuckatoo Cycle edge_bits, used for mining and validating.
const DefaultMinEdgeBits uint8 = 31

// SecondPoWEdgeBits is the Cuckaroo proof-of-work edge_bits, meant to be ASIC resistant.
const SecondPoWEdgeBits uint8 = 29

// BaseEdgeBits is the original reference edge_bits to compute difficulty factors for higher
// Cuckoo graph sizes, changing this would hard fork
const BaseEdgeBits uint8 = 24

// CutThroughHorizon is the Default number of blocks in the past when cross-block cut-through will start
// happening. Needs to be long enough to not overlap with a long reorg.
// Rational
// behind the value is the longest bitcoin fork was about 30 blocks, so 5h. We
// add an order of magnitude to be safe and round to 7x24h of blocks to make it
// easier to reason about.
const CutThroughHorizon uint32 = uint32(WeekHeight)

// StateSyncThreshold is the default number of blocks in the past to determine the height where we request
// a txhashset (and full blocks from). Needs to be long enough to not overlap with
// a long reorg.
// Rational behind the value is the longest bitcoin fork was about 30 blocks, so 5h.
// We add an order of magnitude to be safe and round to 2x24h of blocks to make it
// easier to reason about.
const StateSyncThreshold uint32 = 2 * uint32(DayHeight)

// BlockInputWeight is the weight of an input when counted against the max block weight capacity
const BlockInputWeight int = 1

// BlockOutputWeight is the weight of an output when counted against the max block weight capacity
const BlockOutputWeight int = 21

// BlockKernelWeight is the weight of a kernel when counted against the max block weight capacity
const BlockKernelWeight int = 3

// MaxBlockWeight is the total maximum block weight. At current sizes, this means a maximum
// theoretical size of:
// * `(674 + 33 + 1) * (40_000 / 21) = 1_348_571` for a block with only outputs
// * `(1 + 8 + 8 + 33 + 64) * (40_000 / 3) = 1_520_000` for a block with only kernels
// * `(1 + 33) * 40_000 = 1_360_000` for a block with only inputs
//
// Regardless of the relative numbers of inputs/outputs/kernels in a block the maximum
// block size is around 1.5MB
// For a block full of "average" txs (2 inputs, 2 outputs, 1 kernel) we have -
// `(1 * 2) + (21 * 2) + (3 * 1) = 47` (weight per tx)
// `40_000 / 47 = 851` (txs per block)
//
const MaxBlockWeight int = 40000

// HardForkInterval every 6 months.
const HardForkInterval uint64 = YearHeight / 2

// FloonetFirstHardFork is the Floonet first hard fork height, set to happen around 2019-06-23
const FloonetFirstHardFork uint64 = 185040

// FloonetSecondHardFork is the Floonet second hard fork height, set to happen around 2019-12-19
const FloonetSecondHardFork uint64 = 298080

// HeaderVersion compute possible block version at a given height, implements
// 6 months interval scheduled hard forks for the first 2 years.
func HeaderVersion(chainType ChainType, height uint64) uint16 {
	hfInterval := uint16(1 + height/HardForkInterval)
	switch chainType {
	case Floonet:
		if height < FloonetFirstHardFork {
			return 1
		} else if height < FloonetSecondHardFork {
			return 2
		} else if height < 3*HardForkInterval {
			return 3
		} else {
			return hfInterval
		}
	// everything else just like mainnet
	default:
		return hfInterval
	}
}

// ValidHeaderVersion check whether the block version is valid at a given height, implements
// 6 months interval scheduled hard forks for the first 2 years.
func ValidHeaderVersion(chainType ChainType, height uint64, version uint16) bool {
	return height < 3*HardForkInterval && version == HeaderVersion(chainType, height)
}

// DifficultyAdjustWindow is the number of blocks used to calculate difficulty adjustments
const DifficultyAdjustWindow uint64 = HourHeight

// BlockTimeWindow is the average time span of the difficulty adjustment window
const BlockTimeWindow uint64 = DifficultyAdjustWindow * BlockTimeSec

// ClampFactor is the clamp factor to use for difficulty adjustment
// Limit value to within this factor of goal
const ClampFactor uint64 = 2

// DifficultyDampFactor is the Dampening factor to use for difficulty adjustment
const DifficultyDampFactor uint64 = 3

// ARScaleDampFactor is the dampening factor to use for AR scale calculation.
const ARScaleDampFactor uint64 = 13

// GraphWeight is a weight of a graph as number of siphash bits defining the graph
// Must be made dependent on height to phase out smaller size over the years
// This can wait until end of 2019 at latest
func GraphWeight(chainType ChainType, height uint64, edgeBits uint8) uint64 {
	xprEdgeBits := uint64(edgeBits)
	expiryHeight := YearHeight
	if edgeBits == 31 && height >= expiryHeight {
		xprEdgeBits = saturatingSubUint64(xprEdgeBits, 1+(height-expiryHeight)/WeekHeight)
	}
	return (uint64(2) << uint64(edgeBits-baseEdgeBits(chainType))) * xprEdgeBits
}

// MinDifficulty is the minimum difficulty, enforced in diff retargetting
// avoids getting stuck when trying to increase difficulty subject to dampening
const MinDifficulty uint64 = DifficultyDampFactor

// MinArScale is the minimum scaling factor for AR pow, enforced in diff retargetting
// avoids getting stuck when trying to increase ar_scale subject to dampening
const MinArScale uint64 = ARScaleDampFactor

// UnitDifficulty is the unit difficulty, equal to graph_weight(SECOND_POW_EDGE_BITS)
const UnitDifficulty uint64 = (uint64(2) << (SecondPoWEdgeBits - BaseEdgeBits)) * uint64(SecondPoWEdgeBits)

// InitialDifficulty is the initial difficulty at launch. This should be over-estimated
// and difficulty should come down at launch rather than up
// Currently grossly over-estimated at 10% of current
// ethereum GPUs (assuming 1GPU can solve a block at diff 1 in one block interval)
const InitialDifficulty uint64 = 1000000 * UnitDifficulty

// HeaderInfo is a header info and contains the minimal header information required for the Difficulty calculation to
// take place
type HeaderInfo struct {
	// Timestamp of the header, 1 when not used (returned info)
	timestamp uint64
	// Network difficulty or next difficulty to use
	difficulty Difficulty
	// Network secondary PoW factor or factor to use
	secondaryScaling uint32
	// Whether the header is a secondary proof of work
	isSecondary bool
}

// HeaderInfoFromTsDiff is a constructor from a timestamp and difficulty, setting a default secondary
// PoW factor
func HeaderInfoFromTsDiff(chainType ChainType, timestamp uint64, difficulty Difficulty) HeaderInfo {
	return HeaderInfo{
		timestamp:        timestamp,
		difficulty:       difficulty,
		secondaryScaling: initialGraphWeight(chainType),
		isSecondary:      true,
	}
}

// HeaderInfoFromDiffScaling is a constructor from a difficulty and secondary factor, setting a default
// timestamp
func HeaderInfoFromDiffScaling(difficulty Difficulty, secondaryScaling uint32) HeaderInfo {
	return HeaderInfo{
		timestamp:        1,
		difficulty:       difficulty,
		secondaryScaling: secondaryScaling,
		isSecondary:      true,
	}
}

// Move value linearly toward a goal
func damp(actual, goal, dampFactor uint64) uint64 {
	return (actual + (dampFactor-1)*goal) / dampFactor
}

// limit value to be within some factor from a goal
func clamp(actual, goal, clampFactor uint64) uint64 {
	return uint64(math.Max(float64(goal)/float64(clampFactor),
		math.Min(float64(actual), float64(goal)*float64(clampFactor))))
}
