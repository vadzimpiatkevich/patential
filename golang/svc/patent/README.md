# Patent Service

The Patent service is the service responsible for managing **patents**.

## Tools and languages

To work and run the service the following tools are needed:

- Go
- Bazel

### Build

To build the Patent service run the following command:

```shell
bazel build --sandbox_debug -- //golang/svc/patent/...
```

### Run

The service can be run using Bazel with the following command:

```shell
$ bazel run //golang/svc/patent/cmd --db-host=...
```

You can invoke any gRPC using grpcurl (https://github.com/fullstorydev/grpcurl). As an example:

```shell
grpcurl -d '{ "pagination": { "limit": 1 } }' -plaintext localhost:9000 patent.Service/ListPatents
```
