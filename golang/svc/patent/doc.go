/*
Package patent implements an internal service responsible for managing patents.
The service doesn't imply having any authentication mechanism as long as it's
supposed to be run inside the k8s cluster.

Tools and languages

To work, run and deploy the service the following tools are needed:
- Go: https://golang.org/dl/
- Bazel: https://docs.bazel.build/versions/master/install-os-x.html
- Skaffold: https://skaffold.dev/docs/install/
- Kustomize: https://github.com/kubernetes-sigs/kustomize/blob/master/docs/INSTALL.md

Build

To build the Patent service run the following command:
$ bazel build --sandbox_debug -- //golang/svc/patent/...

Run

The service can be run using Bazel with the following command:
$ bazel run //golang/svc/patent/cmd --db-host=...

You can invoke any gRPC using grpcurl (https://github.com/fullstorydev/grpcurl).
As an example:
$ grpcurl -d '{ "pagination": { "limit": 1 } }' -plaintext localhost:9000 patent.Service/ListPatents

Deploy

The service can be built & deployed using Skaffol. Run Skaffold inside the `k8s`
directory to build and deploy the service once, similar to a CI/CD pipeline:
$ skaffold run
*/
package patent
