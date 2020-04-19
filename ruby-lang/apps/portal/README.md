#### Proto dependencies management

This Rails app depends on the following protobuf libraries:
- golang/svc/patent/proto/service/service.proto.

To produce Ruby library from the proto file, run the following command from the project root:

```shell
grpc_tools_ruby_protoc --ruby_out=lib/proto/svc/patent/ --grpc_out=lib/proto/svc/patent/ --proto_path=../../../golang/svc/patent/proto/service service.proto
```

#### Deploy

The app can be built & deployed using Skaffol. Run Skaffold from the root
directory to build and deploy the service once, similar to a CI/CD pipeline:
```shell
$ skaffold run
```
