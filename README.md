# T9

This is my reference implementation of a [T9](https://en.wikipedia.org/wiki/T9_\(predictive_text\)) server and client. There is particular emphasis on modularity with clear boundaries, reliability, and performance.

## Building

### Prerequisites

The build process requires the following tools installed:

* [Make](https://www.gnu.org/software/make/)
* [OpenSSL](https://www.openssl.org/)
* The [Go](https://golang.org/) toolchain
* Go [Dep](https://github.com/golang/dep)
* [Go-Bindata](https://github.com/jteeuwen/go-bindata)
* [Gomock](https://github.com/golang/mock)
* [Node.js](https://nodejs.org/en/) and [NPM](https://www.npmjs.com/)
* [Protobuf compiler](https://github.com/google/protobuf/blob/master/README.md#protocol-compiler-installation) and [Protobuf Go Runtime](https://github.com/golang/protobuf)
* [GRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
* Optionally, [Docker](https://www.docker.com/) and [Docker-compose](https://docs.docker.com/compose/)

### Building and installing locally

Once the prerequisites are installed, simply typing the below will build and then install on the local system:

```shell
$ make
```

### Building for deployment on Linux

Alternatively, you may wish to build a binary and then scp it to a server for deployment. You can build such a binary for deployment to an AMD64/Linux server with the command:

```shell
$ make linux
```

### Building for deployment with Docker

If you have installed *Docker* and *Docker-Compose* (see **Prerequisites**), then you may build a Docker container. Simply type:

```shell
make docker
```

### Cleaning the workspace

If any of the build steps fail, then it is advisable to clean up the workspace before attempting another build. Type:

```shell
$ make clean
```

## Running

On your local machine after local installation, run the server by typing:

```shell
$ t9 server
```

While the server is running you may either:
* Direct your browser to the address specified
* Run the CLI client by typing:

```shell
$ t9 client
```

## TODO

* Support lets-encrypt?
* Add cache-control/expires headers
* Unit tests with mocks
* Update frontend to use vanilla js instead of jquery and semantic
* Support Brotli (once there is native, pure-go support)