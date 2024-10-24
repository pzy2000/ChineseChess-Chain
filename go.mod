module chainmaker_web

go 1.21

//toolchain go1.21.5

require (
	chainmaker.org/chainmaker/common/v2 v2.3.4
	chainmaker.org/chainmaker/contract-utils v1.0.5-0.20240823082826-c715f4113db5
	chainmaker.org/chainmaker/pb-go/v2 v2.3.5
	chainmaker.org/chainmaker/sdk-go/v2 v2.3.1
	chainmaker.org/chainmaker/tcip-go/v2 v2.3.2-0.20240320084326-98d1c3f10d83
	chainmaker.org/chainmaker/utils/v2 v2.3.5
	github.com/emirpasic/gods v1.12.0
	github.com/ethereum/go-ethereum v1.13.2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.4.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hokaccha/go-prettyjson v0.0.0-20201222001619-a42f9ac2ec8e
	github.com/panjf2000/ants/v2 v2.4.8
	github.com/prometheus/client_golang v1.14.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/viper v1.9.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.389
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms v1.0.389
	github.com/testcontainers/testcontainers-go v0.27.0
	go.uber.org/zap v1.26.0
	golang.org/x/net v0.27.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	gorm.io/driver/clickhouse v0.5.1
	gorm.io/driver/mysql v1.5.2
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.7-0.20240204074919-46816ad31dde
)

require (
	chainmaker.org/chainmaker/contract-sdk-go/v2 v2.3.8 // indirect
	chainmaker.org/chainmaker/protocol/v2 v2.3.5 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/ClickHouse/ch-go v0.60.0 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.15.0 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/Microsoft/hcsshim v0.11.4 // indirect
	github.com/Rican7/retry v0.1.0 // indirect
	github.com/andybalholm/brotli v1.0.6 // indirect
	github.com/aws/aws-sdk-go v1.43.21 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/containerd v1.7.11 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dgryski/go-metro v0.0.0-20200812162917-85c65e2d0165 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v24.0.7+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.3 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/linvon/cuckoo-filter v0.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/moby/patternmatcher v0.6.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc5 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/paulmach/orb v0.10.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pingcap/errors v0.11.5-0.20190809092503-95897b64e011 // indirect
	github.com/pingcap/log v0.0.0-20200511115504-543df19646ad // indirect
	github.com/pingcap/parser v0.0.0-20200623164729-3a18f1e5dceb // indirect
	github.com/pingcap/tidb v1.1.0-beta.0.20200630082100-328b6d0a955c // indirect
	github.com/pingcap/tipb v0.0.0-20210425040103-dc47a87b52aa // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/shirou/gopsutil/v3 v3.23.11 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tidwall/gjson v1.10.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tidwall/tinylru v1.1.0 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/exp v0.0.0-20230810033253-352e893a4cad // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/genproto v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

exclude github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1
