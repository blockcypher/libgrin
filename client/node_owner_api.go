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

package client

import (
	"encoding/json"
	"errors"

	"github.com/blockcypher/libgrin/v5/api"
	"github.com/blockcypher/libgrin/v5/p2p"
	log "github.com/sirupsen/logrus"
)

// NodeOwnerAPI represents the node owner API (v2)
type NodeOwnerAPI struct {
	client RPCHTTPClient
}

// NewNodeOwnerAPI creates a new node owner API
func NewNodeOwnerAPI(url string) *NodeOwnerAPI {
	return &NodeOwnerAPI{client: RPCHTTPClient{URL: url}}
}

// GetStatus returns various information about the node, the network
// and the current sync status.
func (owner *NodeOwnerAPI) GetStatus() (*api.Status, error) {
	envl, err := owner.client.Request("get_status", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetStatus")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var status api.Status
	if err := json.Unmarshal(result.Ok, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// ValidateChain triggers a validation of the chain state.
func (owner *NodeOwnerAPI) ValidateChain() error {
	envl, err := owner.client.Request("validate_chain", nil)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during ValidateChain")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	return nil
}

// CompactChain triggers a compaction of the chain state to regain storage space.
func (owner *NodeOwnerAPI) CompactChain() error {
	envl, err := owner.client.Request("compact_chain", nil)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during CompactChain")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	return nil
}

// GetPeers retrieves information about stored peers.
// If None is provided, will list all stored peers.
func (owner *NodeOwnerAPI) GetPeers(peerAddr *string) (*[]p2p.PeerData, error) {
	paramsBytes, err := json.Marshal(peerAddr)
	if err != nil {
		return nil, err
	}
	envl, err := owner.client.Request("get_peers", paramsBytes)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetPeers")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var peersData []p2p.PeerData
	if err := json.Unmarshal(result.Ok, &peersData); err != nil {
		return nil, err
	}
	return &peersData, nil
}

// GetConnectedPeers retrieve information about stored peers. If None is provided,
// will list all stored peers.
func (owner *NodeOwnerAPI) GetConnectedPeers() (*[]p2p.PeerInfoDisplay, error) {
	envl, err := owner.client.Request("get_connected_peers", nil)
	if err != nil {
		return nil, err
	}
	if envl == nil {
		return nil, errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during GetConnectedPeers")
		return nil, errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return nil, err
	}
	if result.Err != nil {
		return nil, errors.New(string(result.Err))
	}
	var peers []p2p.PeerInfoDisplay
	if err := json.Unmarshal(result.Ok, &peers); err != nil {
		return nil, err
	}
	return &peers, nil
}

// BanPeer bans a specific peer.
func (owner *NodeOwnerAPI) BanPeer(peerAddr *string) error {
	paramsBytes, err := json.Marshal(peerAddr)
	if err != nil {
		return err
	}
	envl, err := owner.client.Request("ban_peer", paramsBytes)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during BanPeer")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	var peersData []p2p.PeerData
	if err := json.Unmarshal(result.Ok, &peersData); err != nil {
		return err
	}
	return nil
}

// UnbanPeer unbans a specific peer.
func (owner *NodeOwnerAPI) UnbanPeer(peerAddr *string) error {
	paramsBytes, err := json.Marshal(peerAddr)
	if err != nil {
		return err
	}
	envl, err := owner.client.Request("unban_peer", paramsBytes)
	if err != nil {
		return err
	}
	if envl == nil {
		return errors.New("NodeOwnerAPI: Empty RPC Response from grin")
	}
	if envl.Error != nil {
		log.WithFields(log.Fields{
			"code":    envl.Error.Code,
			"message": envl.Error.Message,
		}).Error("NodeOwnerAPI: RPC Error during UnbanPeer")
		return errors.New(string(envl.Error.Code) + "" + envl.Error.Message)
	}
	var result Result
	if err = json.Unmarshal(envl.Result, &result); err != nil {
		return err
	}
	if result.Err != nil {
		return errors.New(string(result.Err))
	}
	return nil
}
