// Copyright 2019 BlockCypher
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

package keychain

// IdentifierSize is the size of an identifier in bytes
const IdentifierSize uint = 17

// Identifier is a  semi-opaque structure (just bytes) to track keys
// within the Keychain.
type Identifier [IdentifierSize]uint8
