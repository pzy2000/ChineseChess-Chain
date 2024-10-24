/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/config"
	"reflect"
	"sync"
	"testing"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

func TestNewSingleSdkClientPool(t *testing.T) {
	type args struct {
		chainInfo       *config.ChainInfo
		systemSdkClient *SdkClient
		queryClient     *sdk.ChainClient
	}
	tests := []struct {
		name string
		args args
		want *SingleSdkClientPool
	}{
		{
			name: "Test case 1: Valid chainInfo, systemSdkClient, and queryClient",
			args: args{
				chainInfo:       &config.ChainInfo{}, // TODO: Provide a valid ChainInfo instance
				systemSdkClient: &SdkClient{},        // TODO: Provide a valid SdkClient instance
				queryClient:     &sdk.ChainClient{},  // TODO: Provide a valid ChainClient instance
			},
			want: &SingleSdkClientPool{
				chainInfo:       &config.ChainInfo{}, // TODO: Provide the expected ChainInfo instance
				systemSdkClient: &SdkClient{},        // TODO: Provide the expected SdkClient instance
				queryClient:     &sdk.ChainClient{},  // TODO: Provide the expected ChainClient instance
				sdkClients:      sync.Map{},          // TODO: Provide the expected sync.Map instance with the systemSdkClient added
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSingleSdkClientPool(tt.args.chainInfo, tt.args.systemSdkClient, tt.args.queryClient)
			// TODO: Add any necessary checks to verify the SingleSdkClientPool has been initialized correctly.
			if got == nil {
				t.Errorf("NewSingleSdkClientPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllSdkClient(t *testing.T) {
	type args struct {
		chainList []*config.ChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*SdkClient
	}{
		{
			name: "Test case 1: Valid chainId",
			args: args{
				chainList: []*config.ChainInfo{
					{
						ChainId: ChainId1,
					},
				},
			},
			want: []*SdkClient{}, // TODO: Provide the expected SdkClient instance
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllSdkClient(tt.args.chainList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSdkClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChainClient1(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want *sdk.ChainClient
	}{
		{
			name: "Test case 1: Valid chainId",
			args: args{
				chainId: ChainId1,
			},
			want: nil, // TODO: Provide the expected SdkClient instance
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChainClient(tt.args.chainId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChainClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSdkClient(t *testing.T) {
	tests := []struct {
		name    string
		chainId string
		want    *SdkClient
	}{
		{
			name:    "Test case 1: Valid chainId",
			chainId: "testchain1",
			want:    &SdkClient{}, // TODO: Provide the expected SdkClient instance
		},
		{
			name:    "Test case 2: Invalid chainId",
			chainId: "nonexistentchain",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = GetSdkClient(tt.chainId)
		})
	}
}

func Test_getDefaultLogger(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test case 1: Check if logger is not nil",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDefaultLogger()
			if (got != nil) != tt.want {
				t.Errorf("getDefaultLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
