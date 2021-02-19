module github.com/nioshield/titan-lightning

go 1.14

require (
	cloud.google.com/go/pubsub v1.3.1 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/arthurkiller/rollingwriter v1.1.2
	github.com/cheggaaa/pb/v3 v3.0.6 // indirect
	github.com/cockroachdb/pebble v0.0.0-20210217155127-444296cfa2bb // indirect
	github.com/distributedio/configo v0.0.0-20200107073829-efd79b027816
	github.com/distributedio/titan v0.6.1-0.20210207122117-7ae6bc731ae1
	github.com/docker/go-units v0.4.0
	github.com/fsouza/fake-gcs-server v1.19.0 // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/google/pprof v0.0.0-20201218002935-b9804c9f04c2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/joho/sqltocsv v0.0.0-20210208114054-cb2c3a95fb99 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/montanaflynn/stats v0.6.4 // indirect
	github.com/pingcap/br v5.0.0-rc.0.20201223100334-c344d1edf20c+incompatible // indirect
	github.com/pingcap/kvproto v0.0.0-20210204074845-dd36cf2e1c6b
	github.com/pingcap/tidb v1.1.0-beta.0.20210105101819-f55e8f2bf835 // indirect
	github.com/pingcap/tidb-lightning v4.0.10+incompatible
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.2.0 // indirect
	github.com/sirupsen/logrus v1.7.1 // indirect
	github.com/tipsio/tips v0.0.0-20190604032214-b4d2924f0a97 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible // indirect
	github.com/xitongsys/parquet-go v1.6.0 // indirect
	github.com/xitongsys/parquet-go-source v0.0.0-20201108113611-f372b7d813be // indirect
	go.opencensus.io v0.22.5 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210203152818-3206188e46ba // indirect
	google.golang.org/grpc v1.35.0 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	modernc.org/mathutil v1.2.2 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

replace google.golang.org/grpc v1.35.0 => google.golang.org/grpc v1.27.0

replace gopkg.in/stretchr/testify.v1 => github.com/stretchr/testify v1.2.2
