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

// Define these here, as they should be developer-set, not really tweakable /
//by users

// AutomatedTestingMinEdgeBits is the automated testing edgebits
const AutomatedTestingMinEdgeBits uint8 = 9

// AutomatedTestingProofSize is the automated testing proof size
const AutomatedTestingProofSize int = 4

// UserTestingMinEdgeBits is the user testing edgebits
const UserTestingMinEdgeBits uint8 = 15

// UserTestingProofSize is theuser testing proof size
const UserTestingProofSize int = 42

// AutomatedTestingCoinbaseMaturity is the automated testing coinbase maturity
const AutomatedTestingCoinbaseMaturity uint64 = 3

// UserTestingCoinbaseMaturity is the user testing coinbase maturity
const UserTestingCoinbaseMaturity uint64 = 3

// TestingCutThroughHorizon is the testing cut through horizon in blocks
const TestingCutThroughHorizon uint32 = 70

// TestingStateSyncThreshold is the testing state sync threshold in blocks
const TestingStateSyncThreshold uint32 = 20

// TestingInitialGraphWeight is the testing initial graph weight
const TestingInitialGraphWeight uint32 = 1

// TestingInitialDifficulty is the testing initial block difficulty
const TestingInitialDifficulty uint64 = 1

// TestingMaxBlockWeight is the testing maxblockweight (artificially low, just enough to support a few txs).
const TestingMaxBlockWeight int = 150

// ChainType is the type of chain a server can run with, dictates the genesis block and / and
//mining parameters used.
type ChainType int

const (
	// AutomatedTesting is for CI testing
	AutomatedTesting ChainType = iota
	// UserTesting is for User testing
	UserTesting
	// Floonet is the protocol testing network
	Floonet
	// Mainnet is the main production network
	Mainnet
)

func (c ChainType) shortname() string {
	var shortName string
	switch c {
	case AutomatedTesting:
		shortName = "auto"
	case UserTesting:
		shortName = "user"
	case Floonet:
		shortName = "floo"
	case Mainnet:
		shortName = "main"
	}
	return shortName
}

// The minimum acceptable edge_bits
func minEdgeBits(chainType ChainType) uint8 {
	switch chainType {
	case AutomatedTesting:
		return AutomatedTestingMinEdgeBits
	case UserTesting:
		return UserTestingMinEdgeBits
	default:
		return DefaultMinEdgeBits
	}
}

// Reference edge_bits used to compute factor on higher Cuck(at)oo graph sizes,
// while the min_edge_bits can be changed on a soft fork, changing /
//base_edge_bits is a hard fork.
func baseEdgeBits(chainType ChainType) uint8 {
	switch chainType {
	case AutomatedTesting:
		return AutomatedTestingMinEdgeBits
	case UserTesting:
		return UserTestingMinEdgeBits
	default:
		return BaseEdgeBits
	}
}

// ChainTypeProofSize return the proof size based on the proofsize
func ChainTypeProofSize(chainType ChainType) int {
	switch chainType {
	case AutomatedTesting:
		return AutomatedTestingProofSize
	case UserTesting:
		return UserTestingProofSize
	default:
		return ProofSize
	}
}

// Coinbase maturity for coinbases to be spent
func coinbaseMaturity(chainType ChainType) uint64 {
	switch chainType {
	case AutomatedTesting:
		return AutomatedTestingCoinbaseMaturity
	case UserTesting:
		return UserTestingCoinbaseMaturity
	default:
		return CoinbaseMaturity
	}
}

// Initial mining difficulty
func initialBlockDifficulty(chainType ChainType) uint64 {
	switch chainType {
	case AutomatedTesting:
		return TestingInitialDifficulty
	case UserTesting:
		return TestingInitialDifficulty
	case Floonet:
		return InitialDifficulty
	case Mainnet:
		return InitialDifficulty
	default:
		return InitialDifficulty
	}
}

// Initial mining secondary scale
func initialGraphWeight(chainType ChainType) uint32 {
	switch chainType {
	case AutomatedTesting:
		return TestingInitialGraphWeight
	case UserTesting:
		return TestingInitialGraphWeight
	case Floonet:
		return uint32(GraphWeight(chainType, 0, SecondPoWEdgeBits))
	case Mainnet:
		return uint32(GraphWeight(chainType, 0, SecondPoWEdgeBits))
	default:
		return uint32(GraphWeight(chainType, 0, SecondPoWEdgeBits))
	}
}

// Maximum allowed block weight.
func maxBlockWeight(chainType ChainType) int {
	switch chainType {
	case AutomatedTesting:
		return TestingMaxBlockWeight
	case UserTesting:
		return TestingMaxBlockWeight
	case Floonet:
		return MaxBlockWeight
	case Mainnet:
		return MaxBlockWeight
	default:
		return MaxBlockWeight
	}
}

// Horizon at which we can cut-through and do full local pruning
func cutThroughHorizon(chainType ChainType) uint32 {
	switch chainType {
	case AutomatedTesting:
		return TestingCutThroughHorizon
	case UserTesting:
		return TestingCutThroughHorizon
	default:
		return CutThroughHorizon
	}
}

// Threshold at which we can request a txhashset (and full blocks from)
func stateSyncThreshold(chainType ChainType) uint32 {
	switch chainType {
	case AutomatedTesting:
		return TestingStateSyncThreshold
	case UserTesting:
		return TestingStateSyncThreshold
	default:
		return StateSyncThreshold
	}
}

// Are we in automated testing mode?
func isAutomatedTestingMode(chainType ChainType) bool {
	return AutomatedTesting == chainType
}

// Are we in user testing mode?
func isUserTestingMode(chainType ChainType) bool {
	return UserTesting == chainType
}

// Are we in production mode? / Production defined as a live public network,
//testnet[n] or mainnet.
func isProductionMode(chainType ChainType) bool {
	return Floonet == chainType || Mainnet == chainType
}

// Are we in floonet? / Note: We do not have a corresponding is_mainnet() as we
//want any tests to be as close / as possible to "mainnet" configuration as
//possible. / We want to avoid missing any mainnet only code paths.
func isFloonet(chainType ChainType) bool {
	return Floonet == chainType
}

// Are we for real?
func isMainnet(chainType ChainType) bool {
	return Mainnet == chainType
}

// Helper function to get a nonce known to create a valid POW on / the genesis
//block, to prevent it taking ages. Should be fine for now / as the genesis
//block POW solution turns out to be the same for every new / block chain at the
//moment
func getGenesisNonce(chainType ChainType) uint64 {
	switch chainType {
	case AutomatedTesting:
		// won't make a difference
		return 0
	case UserTesting:
		// Magic nonce for current genesis block at cuckatoo15
		return 27944
	case Floonet:
		// Placeholder, obviously not the right value
		return 0
	case Mainnet:
		// Placeholder, obviously not the right value
		return 0
	default:
		return 0
	}
}

// Short name representing the current chain type ("floo", "main", etc.)
func chainShortname(chainType ChainType) string {
	return chainType.shortname()
}
