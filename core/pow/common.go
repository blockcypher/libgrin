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
	"encoding/binary"

	"golang.org/x/crypto/blake2b"
)

func blake2BHash256(byteToHash []byte) []byte {
	bits := make([]byte, 32)
	hash := blake2b.Sum256(byteToHash)
	copy(bits, hash[:32])
	return bits
}

func setHeaderNonce(header []uint8, nonce *uint32) [4]uint64 {
	if nonce != nil {
		header = header[:len(header)-4]
		buf := make([]uint8, 4)
		binary.LittleEndian.PutUint32(buf, *nonce)
		header = append(header, buf...)
		return createSiphashKeys(header)
	}
	return createSiphashKeys(header)
}

func createSiphashKeys(header []uint8) [4]uint64 {
	h := blake2BHash256(header)
	var s [4]uint64
	s[0] = binary.LittleEndian.Uint64(h[0:8])
	s[1] = binary.LittleEndian.Uint64(h[8:16])
	s[2] = binary.LittleEndian.Uint64(h[16:24])
	s[3] = binary.LittleEndian.Uint64(h[24:32])
	return s
}

// CuckooParams is a utility struct to calculate commonly used Cuckoo parameters
// calculated from header, nonce, edge_bits, etc.
type CuckooParams struct {
	edgeBits    uint8
	proofSize   int
	numEdges    uint64
	siphashKeys [4]uint64
	edgeMask    uint64
}

// Instantiates new params and calculate edge mask, etc
func (c *CuckooParams) new(edgeBits uint8, proofSize int) CuckooParams {
	numEdges := uint64(1) << edgeBits
	edgeMask := numEdges - 1
	var siphashKeys [4]uint64
	return CuckooParams{edgeBits, proofSize, numEdges, siphashKeys, edgeMask}
}

// Reset the main keys used for siphash from the header and nonce
func (c *CuckooParams) resetHeaderNonce(header []uint8, nonce *uint32) {
	c.siphashKeys = setHeaderNonce(header, nonce)
}

// Return siphash masked for type
func (c *CuckooParams) sipnode(edge, uorv uint64, shift bool) uint64 {
	hashUint64 := SipHash24(c.siphashKeys, 2*edge+uorv, 21)
	masked := hashUint64 & c.edgeMask
	if shift {
		masked <<= 1
		masked |= uorv
	}
	return masked
}
