package db

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"os"
	"testing"

	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	//初始化配置
	// 初始化数据库配置
	redisCfg1, err := InitRedisContainer()
	if err != nil {
		return
	}
	_, err = InitMySQLContainer()
	if err != nil {
		return
	}

	cache.InitRedis(redisCfg1)
	//db.InitDbConn(dbCfg)
	// 运行其他测试
	os.Exit(m.Run())
}

func TestConnectDatabase(t *testing.T) {
	// 初始化数据库配置
	dbCfg, err := InitMySQLContainer()
	if err != nil || dbCfg == nil {
		return
	}
	type args struct {
		dbConfig    *config.DBConf
		useDataBase bool
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.DB
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				dbConfig:    dbCfg,
				useDataBase: true,
			},
			wantErr: true,
		},
		{
			name: "test: case 2",
			args: args{
				dbConfig:    dbCfg,
				useDataBase: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = ConnectDatabase(tt.args.dbConfig, tt.args.useDataBase)
		})
	}
}

func TestCreateDatabase(t *testing.T) {
	type args struct {
		db         *gorm.DB
		database   string
		dbProvider string
	}
	tests := []struct {
		name string
		args args
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateDatabase(tt.args.db, tt.args.database, tt.args.dbProvider)
		})
	}
}

//func TestInitClickHouseTable(t *testing.T) {
//	dbCfg, err := InitClickHouseContainer()
//	if err != nil || dbCfg == nil {
//		return
//	}
//
//	type args struct {
//		chainList []*config.ChainInfo
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "test: case 1",
//			args: args{
//				chainList: []*config.ChainInfo{
//					{
//						ChainId: "chain1",
//					},
//					{
//						ChainId: "chain2",
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			InitClickHouseTable(tt.args.chainList)
//		})
//	}
//}

//func TestInitDbConn(t *testing.T) {
//	clickhouseCfg, _ := InitClickHouseContainer()
//	mysqlCfgInfo, _ := InitMySQLContainer()
//
//	if clickhouseCfg == nil || mysqlCfgInfo == nil {
//		return
//	}
//
//	type args struct {
//		dbConfig *config.DBConf
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "test : case 1",
//			args: args{
//				dbConfig: mysqlCfgInfo,
//			},
//		},
//		{
//			name: "test : case 2",
//			args: args{
//				dbConfig: clickhouseCfg,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			InitDbConn(tt.args.dbConfig)
//		})
//	}
//}

//func TestInitMysqlTable(t *testing.T) {
//	mysqlCfgInfo, err := InitMySQLContainer()
//	if mysqlCfgInfo == nil || err != nil {
//		return
//	}
//	type args struct {
//		chainList []*config.ChainInfo
//	}
//	tests := []struct {
//		name string
//		args args
//	}{
//		{
//			name: "test: case 1",
//			args: args{
//				chainList: []*config.ChainInfo{
//					{
//						ChainId: "chain1",
//					},
//					{
//						ChainId: "chain2",
//					},
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			InitMysqlTable(tt.args.chainList)
//		})
//	}
//}
