# Neon Pulumi Provider 

-----

<div align="center">
    ⭐ The project needs your support! Please leave a star and become a GitHub sponsor! ⭐
</div>

-----

Pulumi provider to manage the [Neon](https://neon.tech/) Postgres projects.

## Contents

- [How to contribute](#how-to-contribute)
    * [Prerequisites](#prerequisites)
    * [Build & test the Neon provider](#build---test-the-neon-provider)
        + [Build the provider and install the plugin](#build-the-provider-and-install-the-plugin)
        + [Test against the example](#test-against-the-example)
- [References](#references)

## How to contribute 

The Neon Pulumi provider is a community effort. Please do not hesitate to raise the Github issue to request new functionality,
or to report regression, or misbehaviour of the provider.

### Prerequisites

You will need to ensure the following tools are installed and present in your `$PATH`:

* [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
* [Go 1.21](https://golang.org/dl/) or 1.latest
* [NodeJS](https://nodejs.org/en/) 14.x.  We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
* [Yarn](https://yarnpkg.com/)
* [TypeScript](https://www.typescriptlang.org/)
* [Python](https://www.python.org/downloads/) (called as `python3`).  For recent versions of MacOS, the system-installed version is fine.
* [.NET](https://dotnet.microsoft.com/download)
* java and gradle

### Build & test the Neon provider

Run to see all available commands:

```commandline
make help
```

1. Run `make build install` to build and install the provider.
2. Run `make gen_examples` to generate the example programs in `examples/` off of the source `examples/yaml` example program.
3. Run `make up` to run the example program in `examples/yaml`.
4. Run `make down` to tear down the example program.

Note that you could execute the commands in the Github CodeSpaces environment using this repository.

#### Build the provider and install the plugin

```bash
make build install
```
   
This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

#### Test against the example
   
```bash
$ cd examples/simple
$ yarn link @pulumi/neon
$ yarn install
$ pulumi stack init test
$ pulumi up
```

#### Additional Details

This repository depends on the [pulumi-go-provider](https://github.com/pulumi/pulumi-go-provider) 
and [Neon Go SDK](https://github.com/kislerdm/neon-sdk-go) libraries. 

## References

* [Neon API](https://api-docs.neon.tech/reference/getting-started-with-neon-api)
* [Neon Go SDK](https://pkg.go.dev/github.com/kislerdm/neon-sdk-go)
* [Pulumi Command provider](https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go)
* [Pulumi Go Provider repository](https://github.com/pulumi/pulumi-go-provider)
