package sync

import (
	"chainmaker_web/src/db"
	"testing"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestGetSenderAndPayerUser(t *testing.T) {
	type args struct {
		chainId  string
		hashType string
		txInfo   *common.Transaction
	}
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")

	tests := []struct {
		name    string
		args    args
		want    *db.SenderPayerUser
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				chainId:  "chain1",
				hashType: crypto.CRYPTO_ALGO_SM3,
				txInfo:   txInfo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetSenderAndPayerUser(tt.args.chainId, tt.args.hashType, tt.args.txInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSenderAndPayerUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
