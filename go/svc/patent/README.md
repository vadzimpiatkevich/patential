# Patent Service

The Patent service is the service responsible for managing **patents/portfolios**.

## Tools and languages

To work and run the service the following tools are needed:

- Go
- Bazel

### Build

To build the Patent service run the following command:

```shell
bazel build --sandbox_debug -- //go/svc/patent/...
```

### Run

The service can be run using Bazel with the following command:

```shell
$ bazel run //go/svc/patent/cmd --patents-db-host=...
```
