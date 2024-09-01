# RUSH

Rush attendance check

### Proto buffer

Run the following command to compile the proto.

```
protoc \
  --proto_path=grpc \
  --go_out=grpc \
  --go_opt=paths=source_relative \
  --go-grpc_out=grpc \
  --go-grpc_opt=paths=source_relative \
  rush.proto
```
