// Copyright 2020 BlockCypher
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

package libwallet

// TorConfig is the Tor configuration
type TorConfig struct {
	// Whether to start tor listener on listener startup (default true)
	UseTorListener bool `json:"use_tor_listener"`
	// Just the address of the socks proxy for now
	SocksProxyAddr string `json:"socks_proxy_addr"`
	// Send configuration directory
	SendConfigDir string `json:"send_config_dir"`
}
