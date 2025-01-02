# Ensure GOPATH/bin is in the PATH
export PATH=$PATH:$(go env GOPATH)/bin

# proto dir
mkdir -p proto

## dingofs
### dingofs-proto/cli2.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    ./third-party/dingofs-proto/cli2.proto

### dingofs-proto/common.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    ./third-party/dingofs-proto/common.proto

### dingofs-proto/copyset.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    ./third-party/dingofs-proto/copyset.proto

### dingofs-proto/heartbeat.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    --go_opt=Mproto/heartbeat.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/heartbeat \
    ./third-party/dingofs-proto/heartbeat.proto

### dingofs-proto/mds.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    --go_opt=Mtopology.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/topology \
    ./third-party/dingofs-proto/mds.proto

### dingofs-proto/metaserver.proto 
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    ./third-party/dingofs-proto/metaserver.proto

### dingofs-proto/schedule.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    ./third-party/dingofs-proto/schedule.proto

### dingofs-proto/space.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    ./third-party/dingofs-proto/space.proto

### dingofs-proto/topology.proto
protoc --go_out=proto --proto_path=./third-party/dingofs-proto \
    --go_opt=Mcommon.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/common \
    --go_opt=Mheartbeat.proto=github.com/dingodb/dingofs-tools/proto/dingofs/proto/heartbeat \
    ./third-party/dingofs-proto/topology.proto

# grpc
## fs
protoc --go-grpc_out=proto --proto_path=./third-party/dingofs-proto ./third-party/dingofs-proto/*.proto