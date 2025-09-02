module github.com/people257/poor-guy-shop/oss-infra

go 1.24.4

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250717165733-d22d418d82d8.1
	github.com/aliyun/aliyun-oss-go-sdk v3.0.2+incompatible
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.6.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/file v1.2.0
	github.com/knadh/koanf/v2 v2.2.2
	github.com/people257/poor-guy-shop/common/db v0.0.0-20250820165901-4f14d03768c9
	github.com/people257/poor-guy-shop/common/gateway v0.0.0-20250805161543-d606e53d9e9b
	github.com/people257/poor-guy-shop/common/server v0.0.0-20250805161543-d606e53d9e9b
	github.com/qiniu/go-sdk/v7 v7.25.4
	github.com/redis/go-redis/extra/redisotel/v9 v9.12.1
	github.com/redis/go-redis/v9 v9.12.1
	github.com/spf13/viper v1.20.1
	go.opentelemetry.io/otel v1.37.0
	go.uber.org/zap v1.27.0
	golang.org/x/sync v0.16.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
	gorm.io/driver/postgres v1.6.0
	gorm.io/gen v0.3.27
	gorm.io/gorm v1.30.1
	gorm.io/plugin/dbresolver v1.6.2
	gorm.io/plugin/opentelemetry v0.1.16
	moul.io/zapgorm2 v1.3.0

)

require (
	buf.build/go/protovalidate v0.14.0 // indirect
	cel.dev/expr v0.24.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/ClickHouse/ch-go v0.61.5 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.30.0 // indirect
	github.com/alex-ant/gomath v0.0.0-20160516115720-89013a210a82 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gammazero/toposort v0.1.1 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.7.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/form v3.1.4+incompatible // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.3.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/google/cel-go v0.25.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.2 // indirect
	github.com/hashicorp/consul/api v1.32.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/labstack/echo/v4 v4.13.4 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/paulmach/orb v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/people257/poor-guy-shop/common/resolver v0.0.0-20250804165728-bde3c4e748dc // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/redis/go-redis/extra/rediscmd/v9 v9.12.1 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.12.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.62.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/runtime v0.62.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/log v0.13.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.13.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.3 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.2.4 // indirect
	gorm.io/driver/clickhouse v0.7.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/hints v1.1.0 // indirect
	modernc.org/fileutil v1.0.0 // indirect
)

replace (
	github.com/people257/poor-guy-shop/common/db => ../common/db
	github.com/people257/poor-guy-shop/common/gateway => ../common/gateway
	github.com/people257/poor-guy-shop/common/server => ../common/server
)
