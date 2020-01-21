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

// Parameters to the siphash block algorithm. Used by Cuckaroo but can be seen
// as a generic way to derive a hash within a block of them.
const sipHashBlockBits uint64 = 6
const sipHashBlockSize uint64 = 1 << sipHashBlockBits
const sipHashBlockMask uint64 = sipHashBlockSize - 1

// SipHashBlock builds a block of siphash values by repeatedly hashing from the
// nonce truncated to its closest block start, up to the end of the block.
// Returns the resulting hash at the nonce's position.
func SipHashBlock(v [4]uint64, nonce uint64, rotE uint8, xorAll bool) uint64 {
	// beginning of the block of hashes
	nonce0 := nonce & ^sipHashBlockMask
	nonceI := nonce & sipHashBlockMask
	nonceHash := make([]uint64, sipHashBlockSize)
	// repeated hashing over the whole block
	s := new(sipHash24)
	siphash := s.new(v)
	var i uint64
	for i = 0; i < sipHashBlockSize; i++ {
		siphash.hash(nonce0+i, rotE)
		nonceHash[i] = siphash.digest()
	}
	// xor the hash at nonce_i < SIPHASH_BLOCK_MASK with some or all later hashes to force hashing the whole block
	var xor uint64 = nonceHash[nonceI]
	var xorFrom uint64
	if xorAll || nonceI == sipHashBlockMask {
		xorFrom = nonceI + 1
	} else {
		xorFrom = sipHashBlockMask
	}

	for i := xorFrom; i < sipHashBlockSize; i++ {
		xor ^= nonceHash[i]
	}
	return xor
}

// SipHash24 is an utility function to compute a single siphash 2-4 based on a seed and a nonce.
func SipHash24(v [4]uint64, nonce uint64, rotE uint8) uint64 {
	s := new(sipHash24)
	siphash := s.new(v)
	siphash.hash(nonce, rotE)
	return siphash.digest()
}

type sipHash24 struct {
	v0, v1, v2, v3 uint64
}

func (s *sipHash24) new(v [4]uint64) sipHash24 {
	return sipHash24{v[0], v[1], v[2], v[3]}
}

// One siphash24 hashing, consisting of 2 and then 4 rounds
func (s *sipHash24) hash(nonce uint64, rotE uint8) {
	s.v3 ^= nonce
	s.round(rotE)
	s.round(rotE)

	s.v0 ^= nonce
	s.v2 ^= 0xff

	for i := 0; i < 4; i++ {
		s.round(rotE)
	}
}

// Resulting hash digest
func (s *sipHash24) digest() uint64 {
	return (s.v0 ^ s.v1) ^ (s.v2 ^ s.v3)
}

func (s *sipHash24) round(rotE uint8) {
	s.v0 = s.v0 + s.v1
	s.v2 = s.v2 + s.v3
	s.v1 = rotl(s.v1, 13)
	s.v3 = rotl(s.v3, 16)
	s.v1 ^= s.v0
	s.v3 ^= s.v2
	s.v0 = rotl(s.v0, 32)
	s.v2 = s.v2 + s.v1
	s.v0 = s.v0 + s.v3
	s.v1 = rotl(s.v1, 17)
	s.v3 = rotl(s.v3, rotE)
	s.v1 ^= s.v2
	s.v3 ^= s.v0
	s.v2 = rotl(s.v2, 32)
}

func rotl(val uint64, shift uint8) uint64 {
	num := (val << shift) | (val >> (64 - shift))
	return num
}
