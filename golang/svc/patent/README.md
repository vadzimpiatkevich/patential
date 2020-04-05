# Patent Service

The Patent service is the service responsible for managing **patents**.

## Tools and languages

To work, run and deploy the service the following tools are needed:

- Go: https://golang.org/dl/
- Bazel: https://docs.bazel.build/versions/master/install-os-x.html
- Skaffold: https://skaffold.dev/docs/install/
- Kustomize: https://github.com/kubernetes-sigs/kustomize/blob/master/docs/INSTALL.md

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


### Deploy

The service can be built & deployed using Skaffol. Use `skaffold run` inside the
`k8s` directory to build and deploy the service once, similar to a CI/CD
pipeline:

```shell
$ skaffold run
```
