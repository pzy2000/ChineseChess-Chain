package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	nodeId1      = "node1"
	nodeId2      = "node2"
	nodeName1    = "NodeName1"
	nodeName2    = "NodeName2"
	nodeAddress1 = "12345"
	nodeAddress2 = "22345"
)

func insertNodeTest() ([]*db.Node, error) {
	insertList := []*db.Node{
		{
			NodeId:   nodeId1,
			NodeName: nodeName1,
			OrgId:    orgId1,
			Address:  nodeAddress1,
		},
		{
			NodeId:   nodeId2,
			NodeName: nodeName2,
			OrgId:    orgId2,
			Address:  nodeAddress2,
		},
	}
	err := BatchInsertNode(ChainID, insertList)
	return insertList, err
}

func TestBatchInsertNode(t *testing.T) {
	insertList := []*db.Node{
		{
			NodeId:   nodeId1,
			NodeName: nodeName1,
			OrgId:    orgId1,
			Address:  nodeAddress1,
		},
		{
			NodeId:   nodeId2,
			NodeName: nodeName2,
			OrgId:    orgId2,
			Address:  nodeAddress2,
		},
	}

	type args struct {
		chainId  string
		nodeList []*db.Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				nodeList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchInsertNode(tt.args.chainId, tt.args.nodeList); (err != nil) != tt.wantErr {
				t.Errorf("BatchInsertNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNodeInById(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}
	type args struct {
		chainId string
		nodeIds []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Node
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				nodeIds: []string{
					nodeId1,
				},
			},
			want: map[string]*db.Node{
				nodeId1: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeInById(tt.args.chainId, tt.args.nodeIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeInById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Node{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetNodeInById() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetNodeList(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		nodeName string
		orgId    string
		nodeId   string
		offset   int
		limit    int
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Node
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				offset:  0,
				limit:   10,
			},
			want:    insertList,
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetNodeList(tt.args.chainId, tt.args.nodeName, tt.args.orgId, tt.args.nodeId, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Node{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetNodeList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if got1 != tt.want1 {
				t.Errorf("GetNodeList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetNodeNum(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		role    string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeNum(tt.args.chainId, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNodeNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNodeNumByOrg(t *testing.T) {
	insertList, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		orgId   string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    int64(len(insertList)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodeNumByOrg(tt.args.chainId, tt.args.orgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodeNumByOrg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNodeNumByOrg() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNodesRef(t *testing.T) {
	_, err := insertNodeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []string{
				nodeId1,
				nodeId2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNodesRef(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodesRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNodesRef() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestUpdateNode(t *testing.T) {
// 	_, err := insertNodeTest()
// 	if err != nil {
// 		return
// 	}

// 	type args struct {
// 		chainId string
// 		nodeIds []string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "test: case 1",
// 			args: args{
// 				chainId: ChainID,
// 				nodeIds: []string{
// 					nodeId1,
// 					nodeId2,
// 				},
// 			},

// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := UpdateNode(tt.args.chainId, tt.args.nodeIds); (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateNode() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
