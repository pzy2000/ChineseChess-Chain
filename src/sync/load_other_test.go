package sync

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"chainmaker.org/chainmaker/pb-go/v2/discovery"
)

func TestGetAdminUserByConfig(t *testing.T) {
	type args struct {
		chainConfig *pbconfig.ChainConfig
	}
	tests := []struct {
		name string
		args args
		want []*db.User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAdminUserByConfig(tt.args.chainConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAdminUserByConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChainNodeData(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name  string
		args  args
		want  []*discovery.Node
		want1 []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetChainNodeData(tt.args.sdkClient)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChainNodeData() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetChainNodeData() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPeriodicCheckSubChainStatus(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PeriodicCheckSubChainStatus(tt.args.sdkClient)
		})
	}
}

func TestPeriodicLoadStart(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PeriodicLoadStart(tt.args.sdkClient)
		})
	}
}

func Test_loadChainRefInfos(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadChainRefInfos(tt.args.sdkClient); (err != nil) != tt.wantErr {
				t.Errorf("loadChainRefInfos() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_loadChainUser(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadChainUser(tt.args.sdkClient); (err != nil) != tt.wantErr {
				t.Errorf("loadChainUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_loadNodeInfo(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadNodeInfo(tt.args.sdkClient)
		})
	}
}

func Test_loadOrgInfo(t *testing.T) {
	type args struct {
		sdkClient *SdkClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := loadOrgInfo(tt.args.sdkClient); (err != nil) != tt.wantErr {
				t.Errorf("loadOrgInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parseNodeInfo(t *testing.T) {
	type args struct {
		node *discovery.Node
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []string
		want2   []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := parseNodeInfo(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNodeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseNodeInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseNodeInfo() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("parseNodeInfo() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_parseNodeList(t *testing.T) {
	type args struct {
		nodeList []*discovery.Node
	}
	tests := []struct {
		name string
		args args
		want []*db.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNodeList(tt.args.nodeList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNodeList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePublicNodeList(t *testing.T) {
	type args struct {
		nodeList         []*discovery.Node
		consensusNodeIds map[string]int
	}
	tests := []struct {
		name string
		args args
		want []*db.Node
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePublicNodeList(tt.args.nodeList, tt.args.consensusNodeIds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePublicNodeList() = %v, want %v", got, tt.want)
			}
		})
	}
}
