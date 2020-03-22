This Rails app depends on the following protobuf libraries:
- go/svc/patent/proto/service/service.proto.

To produce Ruby library from the proto file, run the following command from the project root:

```shell
grpc_tools_ruby_protoc --ruby_out=lib/proto/svc/patent/ --grpc_out=lib/proto/svc/patent/ --proto_path=../../../go/svc/patent/proto/service service.proto
```
