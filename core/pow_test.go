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

package core

import (
	"encoding/binary"
	"testing"

	"github.com/blockcypher/libgrin/core/pow"

	"github.com/blockcypher/libgrin/core/consensus"
	"github.com/stretchr/testify/assert"
)

func TestVerifySize(t *testing.T) {
	prePoW := []uint8{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 109, 239}
	bh := new(BlockHeader)

	bh.PoW.Nonce = 28143
	bh.PoW.Proof.EdgeBits = 15
	bh.PoW.Proof.Nonces = []uint64{749, 873, 927, 1637, 2687, 3668, 4346, 5192, 5787, 6055, 6270, 7064, 7140, 7474, 7805, 9017, 9095, 9492, 10634, 11708, 11785, 11799, 12362, 12498, 12667, 13680, 13941, 15360, 17955, 18519, 18691, 20589, 22113, 23605, 24538, 24871, 24945, 25137, 27372, 29195, 31787, 32687}
	assert.Nil(t, VerifySize(consensus.UserTesting, prePoW, bh))
}

func TestVerifySize2(t *testing.T) {
	prePoW := []uint8{0, 1, 0, 0, 0, 0, 0, 0, 166, 143, 0, 0, 0, 0, 92, 97, 221, 123, 65, 75, 2, 126, 250,
		225, 45, 158, 210, 248, 55, 105, 130, 90, 94, 118, 135, 252, 69, 37, 52, 50, 194, 170,
		145, 153, 2, 30, 128, 217, 37, 221, 147, 20, 62, 61, 126, 215, 67, 191, 62, 42, 186,
		32, 5, 198, 186, 40, 161, 236, 126, 109, 184, 26, 161, 190, 73, 222, 239, 118, 83, 183,
		71, 73, 132, 11, 64, 125, 146, 139, 225, 244, 47, 141, 199, 90, 76, 56, 244, 136, 121,
		185, 23, 216, 41, 133, 1, 241, 186, 2, 77, 240, 57, 63, 126, 162, 125, 227, 34, 164,
		26, 107, 220, 53, 77, 101, 171, 75, 116, 63, 83, 159, 171, 197, 12, 45, 135, 48, 30,
		114, 93, 101, 112, 37, 24, 147, 140, 41, 119, 150, 116, 16, 37, 126, 102, 215, 165,
		123, 55, 207, 43, 81, 54, 220, 0, 255, 147, 139, 28, 4, 108, 212, 207, 194, 201, 226,
		112, 192, 69, 189, 145, 238, 225, 99, 88, 238, 153, 38, 147, 44, 113, 155, 215, 33, 80,
		6, 55, 35, 208, 211, 100, 18, 179, 165, 128, 210, 171, 164, 141, 246, 193, 166, 0, 0,
		0, 0, 0, 2, 245, 206, 0, 0, 0, 0, 0, 2, 7, 66, 0, 0, 0, 3, 199, 132, 143, 238, 0, 0, 0,
		13, 16, 89, 244, 58, 146, 89, 127, 83,
	}
	bh := new(BlockHeader)
	bh.PoW.Nonce = 1178241309934714707
	bh.PoW.Proof.EdgeBits = 29
	bh.PoW.Proof.Nonces = []uint64{
		18852094, 18878486, 39783881, 59379092, 62326621, 71455167, 131832576, 143026722,
		143954436, 155338092, 199111429, 207884782, 211343283, 226025553, 233881058, 237856564,
		244323712, 246236308, 253743368, 258760447, 259000289, 262474233, 268999312, 303276522,
		348709059, 371226190, 380435344, 381559211, 382415438, 385006790, 385328950, 389303551,
		424071479, 431735335, 462433478, 476234373, 512600249, 513776715, 514612369, 518186065,
		526136923, 533118850,
	}
	assert.Nil(t, VerifySize(consensus.Mainnet, prePoW, bh))
}

func TestVerifySizeWithoutNonceInPrePow(t *testing.T) {
	prePoW := []uint8{0, 1, 0, 0, 0, 0, 0, 1, 136, 103, 0, 0, 0, 0, 92, 161, 240, 23, 11, 54, 6, 137, 119, 213, 181, 62, 140, 201, 185, 216, 68, 65, 165, 93, 55, 90, 52, 98, 81, 27, 185, 236, 201, 210, 4, 219, 92, 131, 246, 22, 117, 163, 209, 158, 107, 69, 158, 111, 33, 82, 240, 128, 250, 114, 209, 178, 160, 128, 70, 201, 118, 164, 106, 137, 199, 18, 183, 251, 204, 208, 238, 254, 214, 235, 67, 221, 26, 22, 175, 249, 124, 65, 195, 23, 20, 169, 140, 45, 187, 140, 193, 71, 6, 74, 67, 57, 149, 241, 253, 76, 12, 213, 80, 53, 21, 206, 37, 226, 255, 56, 91, 252, 249, 48, 224, 169, 190, 99, 246, 195, 217, 170, 3, 68, 109, 51, 103, 161, 245, 241, 183, 172, 58, 59, 229, 193, 43, 189, 56, 176, 129, 173, 222, 37, 108, 81, 185, 123, 249, 200, 223, 97, 63, 205, 72, 41, 212, 53, 155, 224, 4, 27, 150, 143, 18, 45, 160, 27, 157, 128, 30, 242, 145, 74, 189, 175, 122, 40, 146, 87, 30, 120, 254, 146, 229, 150, 37, 1, 142, 166, 185, 170, 27, 176, 25, 174, 122, 85, 159, 58, 0, 0, 0, 0, 0, 4, 222, 130, 0, 0, 0, 0, 0, 3, 211, 103, 0, 0, 0, 3, 208, 114, 212, 188, 0, 0, 0, 13}
	bh := new(BlockHeader)
	bh.PoW.Nonce = 16079481998891884557
	// Add nonce to prepow
	nonceBytes := make([]uint8, 8)
	binary.BigEndian.PutUint64(nonceBytes, bh.PoW.Nonce)
	prePoW = append(prePoW, nonceBytes...)
	bh.PoW.Proof.EdgeBits = 29
	bh.PoW.Proof.Nonces = []uint64{
		4950556, 10444042, 26994871, 63816933, 64006601, 70454862, 74408437, 101859857, 103156578, 103619764, 110918645, 112676394, 156469828, 164995210, 177571941, 197003830, 206258400, 232973126, 235492427, 243875402, 250871506, 261431148, 294643091, 315606197, 320713204, 328097841, 331983190, 340029134, 341429798, 349593608, 352254617, 363452582, 376534642, 385998553, 399426703, 399588750, 417560407, 418344217, 464144305, 478639713, 500541067, 511159362}
	assert.Nil(t, VerifySize(consensus.Mainnet, prePoW, bh))
}

// Check that we create the appropriate PoW context
func TestMainnetContext(t *testing.T) {
	var zero []uint64

	// One block before hf
	ctx := createPoWContext(consensus.Mainnet, consensus.YearHeight/2-1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckarooContext{}, ctx)
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight/2-1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// Hard fork height
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight/2, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroodContext{}, ctx)
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight/2, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// After hard fork
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight/2+1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroodContext{}, ctx)
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight/2+1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// One block before second hf
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight-1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroodContext{}, ctx)
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight-1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// Second hard fork height
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroomContext{}, ctx)
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// After second hard fork
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight+1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroomContext{}, ctx)
	ctx = createPoWContext(consensus.Mainnet, consensus.YearHeight+1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)
}

func TestFloonetContext(t *testing.T) {
	var zero []uint64

	// One block before first hf
	ctx := createPoWContext(consensus.Floonet, consensus.FloonetFirstHardFork-1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckarooContext{}, ctx)
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetFirstHardFork-1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// First hard fork height
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetFirstHardFork, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroodContext{}, ctx)
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetFirstHardFork, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// After first hard fork
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetFirstHardFork+1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroodContext{}, ctx)
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetFirstHardFork+1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// One block before second hf
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetSecondHardFork-1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroodContext{}, ctx)
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetSecondHardFork-1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// Second hard fork height
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetSecondHardFork, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroomContext{}, ctx)
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetSecondHardFork, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)

	// After second hard fork
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetSecondHardFork+1, 29, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckaroomContext{}, ctx)
	ctx = createPoWContext(consensus.Floonet, consensus.FloonetSecondHardFork+1, 31, 42, zero, maxSols)
	assert.IsType(t, &pow.CuckatooContext{}, ctx)
}
