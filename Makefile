SHELL := /bin/bash
GOPATH			:= $(shell go env GOPATH)
WORKING_DIR     := $(shell pwd)

GITHANDLE 		 := github.com/kislerdm
PROJECT          := $(GITHANDLE)/pulumi-neon
GO_SDK			 := pulumi-sdk-neon
PACK             := neon
NODE_MODULE_NAME := @neon
NUGET_PKG_NAME   := neon

PROVIDER        := pulumi-resource-${PACK}
VERSION         ?= 0.0.1+$(shell git rev-parse --short HEAD)

.PHONY: help
help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z0-9_-].+::.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

.provider:
	@ CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/provider.Version=${VERSION}"

provider:: .provider gen_schema ## Builds provider.

gen_schema:: ## Generates schema.json.
	@ pulumi package get-schema bin/$(PROVIDER) > schema.json

tests:: ## Runs unit tests.
	cd provider && go test -short -v -count=1 -cover -timeout 30m -parallel 1 ./...

acctests:: ## Runs acc tests.
	cd acc-test && go test -short -v -timeout 30m -parallel 1 ./...

lint:: ## Lints the provider's codebase.
	for DIR in "provider"; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

go_sdk:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/go/go.* sdk-template/go/README.md ## Generates Go SDK.
	@ rm -rf $(GO_SDK)/*
	@ cp -r sdk-template/go/* $(GO_SDK)/ && cp LICENSE $(GO_SDK)/
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o $$(GO_SDK) --language go
	@ cd $(GO_SDK) && mv go/* . && rm -r go
	@ cd $(GO_SDK) && go mod tidy

nodejs_sdk:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json ## Generates Node.js SDK.
	@ rm -rf sdk-nodejs
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-nodejs --language nodejs
	@ cd sdk-nodejs/nodejs/ && \
		npm install && \
		npm run build && \
		cp ../../sdk-readme/README-nodejs.md bin/README.md && \
		cp ../../LICENSE package.json package-lock.json bin/ && \
		sed -i.bak 's/$${VERSION}/$(VERSION)/g' bin/package.json && \
		rm ./bin/package.json.bak

python_sdk:: $(WORKING_DIR)/bin/$(PROVIDER) ## Generates python SDK.
	@ rm -rf sdk-python
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-python --language python
	@ cp sdk-readme/README-python.md sdk-python/python/README.md && cp LICENSE sdk-python/python/
	@ cd sdk-python/python/ && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(VERSION)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist

dotnet_sdk:: $(WORKING_DIR)/bin/$(PROVIDER) ## Generates .Net SDK.
	@ rm -rf sdk-dotnet
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-dotnet --language dotnet
	@ cd sdk-dotnet/dotnet/&& \
		echo "${VERSION}" >version.txt && \
		dotnet build /p:Version=${VERSION}

java_sdk:: $(WORKING_DIR)/bin/$(PROVIDER) ## Generates Go SDK.
	@ rm -rf sdk-java
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-java --language java
