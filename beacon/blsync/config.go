// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package blsync

import (
	"github.com/Kevionte/go-sovereign/beacon/types"
	"github.com/Kevionte/go-sovereign/cmd/utils"
	"github.com/Kevionte/go-sovereign/common"
	"github.com/Kevionte/go-sovereign/common/hexutil"
	"github.com/urfave/cli/v2"
)

// lightClientConfig contains beacon light client configuration
type lightClientConfig struct {
	*types.ChainConfig
	Checkpoint common.Hash
}

var (
	MainnetConfig = lightClientConfig{
		ChainConfig: (&types.ChainConfig{
			GenesisValidatorsRoot: common.HexToHash("0xe7928a5608c8748de824906249149e1279fd118a1b05970caff2594770c6c772"),
			GenesisTime:           1724336895,
		}).
			AddFork("GENESIS", 0, []byte{0, 0, 0, 0}).
			AddFork("ALTAIR", 74240, []byte{1, 0, 0, 0}).
			AddFork("BELLATRIX", 144896, []byte{2, 0, 0, 0}).
			AddFork("CAPELLA", 194048, []byte{3, 0, 0, 0}).
			AddFork("DENEB", 269568, []byte{4, 0, 0, 0}),
		Checkpoint: common.HexToHash("0xe7928a5608c8748de824906249149e1279fd118a1b05970caff2594770c6c772"),
	}

	GoerliConfig = lightClientConfig{
		ChainConfig: (&types.ChainConfig{
			GenesisValidatorsRoot: common.HexToHash("0x043db0d9a83813551ee2f33450d23797757d430911a9320530ad8a0eabc43efb"),
			GenesisTime:           1614588812,
		}).
			AddFork("GENESIS", 0, []byte{0, 0, 16, 32}).
			AddFork("ALTAIR", 36660, []byte{1, 0, 16, 32}).
			AddFork("BELLATRIX", 112260, []byte{2, 0, 16, 32}).
			AddFork("CAPELLA", 162304, []byte{3, 0, 16, 32}).
			AddFork("DENEB", 231680, []byte{4, 0, 16, 32}),
		Checkpoint: common.HexToHash("0x53a0f4f0a378e2c4ae0a9ee97407eb69d0d737d8d8cd0a5fb1093f42f7b81c49"),
	}
)

func makeChainConfig(ctx *cli.Context) lightClientConfig {
	var config lightClientConfig
	customConfig := ctx.IsSet(utils.BeaconConfigFlag.Name)
	utils.CheckExclusive(ctx, utils.MainnetFlag, utils.GoerliFlag, utils.BeaconConfigFlag)
	switch {
	case ctx.Bool(utils.MainnetFlag.Name):
		config = MainnetConfig
	case ctx.Bool(utils.GoerliFlag.Name):
		config = GoerliConfig
	default:
		if !customConfig {
			config = MainnetConfig
		}
	}
	// Genesis root and time should always be specified together with custom chain config
	if customConfig {
		if !ctx.IsSet(utils.BeaconGenesisRootFlag.Name) {
			utils.Fatalf("Custom beacon chain config is specified but genesis root is missing")
		}
		if !ctx.IsSet(utils.BeaconGenesisTimeFlag.Name) {
			utils.Fatalf("Custom beacon chain config is specified but genesis time is missing")
		}
		if !ctx.IsSet(utils.BeaconCheckpointFlag.Name) {
			utils.Fatalf("Custom beacon chain config is specified but checkpoint is missing")
		}
		config.ChainConfig = &types.ChainConfig{
			GenesisTime: ctx.Uint64(utils.BeaconGenesisTimeFlag.Name),
		}
		if c, err := hexutil.Decode(ctx.String(utils.BeaconGenesisRootFlag.Name)); err == nil && len(c) <= 32 {
			copy(config.GenesisValidatorsRoot[:len(c)], c)
		} else {
			utils.Fatalf("Invalid hex string", "beacon.genesis.gvroot", ctx.String(utils.BeaconGenesisRootFlag.Name), "error", err)
		}
		if err := config.ChainConfig.LoadForks(ctx.String(utils.BeaconConfigFlag.Name)); err != nil {
			utils.Fatalf("Could not load beacon chain config file", "file name", ctx.String(utils.BeaconConfigFlag.Name), "error", err)
		}
	} else {
		if ctx.IsSet(utils.BeaconGenesisRootFlag.Name) {
			utils.Fatalf("Genesis root is specified but custom beacon chain config is missing")
		}
		if ctx.IsSet(utils.BeaconGenesisTimeFlag.Name) {
			utils.Fatalf("Genesis time is specified but custom beacon chain config is missing")
		}
	}
	// Checkpoint is required with custom chain config and is optional with pre-defined config
	if ctx.IsSet(utils.BeaconCheckpointFlag.Name) {
		if c, err := hexutil.Decode(ctx.String(utils.BeaconCheckpointFlag.Name)); err == nil && len(c) <= 32 {
			copy(config.Checkpoint[:len(c)], c)
		} else {
			utils.Fatalf("Invalid hex string", "beacon.checkpoint", ctx.String(utils.BeaconCheckpointFlag.Name), "error", err)
		}
	}
	return config
}
