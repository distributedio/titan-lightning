module github.com/nioshield/titan-lightning

go 1.14

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/cheggaaa/pb/v3 v3.0.6 // indirect
	github.com/cockroachdb/pebble v0.0.0-20210217155127-444296cfa2bb // indirect
	github.com/distributedio/configo v0.0.0-20200107073829-efd79b027816
	github.com/docker/go-units v0.4.0
	github.com/fsouza/fake-gcs-server v1.22.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/joho/sqltocsv v0.0.0-20210208114054-cb2c3a95fb99 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/montanaflynn/stats v0.6.4 // indirect
	github.com/pingcap/br v5.0.0-rc.0.20201223100334-c344d1edf20c+incompatible // indirect
	github.com/pingcap/kvproto v0.0.0-20210204074845-dd36cf2e1c6b // indirect
	github.com/pingcap/tidb v1.1.0-beta.0.20210105101819-f55e8f2bf835 // indirect
	github.com/pingcap/tidb-lightning v4.0.10+incompatible
	github.com/pingcap/tipb v0.0.0-20210204051656-2870a0852037 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/xitongsys/parquet-go v1.6.0 // indirect
	github.com/xitongsys/parquet-go-source v0.0.0-20201108113611-f372b7d813be // indirect
	modernc.org/mathutil v1.2.2 // indirect
)

replace google.golang.org/grpc v1.35.0 => google.golang.org/grpc v1.27.0
