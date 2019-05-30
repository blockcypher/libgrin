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

import "fmt"

// ShareType is the possible type of shares
type ShareType int

const (
	// ValidShare represents a valid share
	ValidShare ShareType = iota
	// StaleShare represents a slate share
	StaleShare
	// RejectedShare represents a rejected share
	RejectedShare
	// RejectedLowDiffShare represents a rejected share because of low difficulty
	RejectedLowDiffShare
	// RejectedDuplicateShare represents a rejected share because duplicate
	RejectedDuplicateShare
	// ValidHighDiffShare represents a valid share which is above the internal target diff
	ValidHighDiffShare
)

// PayoutMethodType The payout method type
type PayoutMethodType int

const (
	// HTTPPayoutMethod is the http payout method
	HTTPPayoutMethod PayoutMethodType = iota
	// FilePayoutMethod is the file payout method
	FilePayoutMethod
	// InvalidPayoutMethod is an invalid payout method
	InvalidPayoutMethod
)

func (p PayoutMethodType) methodTypeToString() string {
	switch p {
	case HTTPPayoutMethod:
		return "http"
	case FilePayoutMethod:
		return "file"
	default:
		return "http"
	}
}

func stringToMethodType(method string) (PayoutMethodType, error) {
	switch method {
	case "http":
		return HTTPPayoutMethod, nil
	case "file":
		return FilePayoutMethod, nil
	default:
		return InvalidPayoutMethod, fmt.Errorf("invalid payout method %s", method)
	}
}
