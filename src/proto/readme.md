### Generate go code from Protocol Buffers.

_Command_: `protoc --go_out=src/grpc --go-grpc_out=src/grpc src/proto/[folder_name]/[file_name].proto`

### Why do we use Protocol Buffers?

- Protocol Buffers is faster than HTTP API calls because of the efficiency in serialization and deserialization compared to JSON.
- It generates the client for us, saving time.
