SHELL := /bin/bash
GOPATH			:= $(shell go env GOPATH)
WORKING_DIR     := $(shell pwd)

GITHANDLE 		 := github.com/kislerdm
PROJECT          := $(GITHANDLE)/pulumi-neon
GO_SDK			 := pulumi-sdk-neon
PACK             := neon
PY_PKG_NAME		 := pulumi_neon

PROVIDER        := pulumi-resource-${PACK}

VERSION 		:=
VERSION_SET     ?= $(if $(VERSION),$(VERSION),$(shell pulumictl get version))
VERSION_PY      ?= $(if $(VERSION),$(VERSION),$(shell pulumictl get version --language python))

.PHONY: help
help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z0-9_-].+::.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

.provider:
	@ CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/provider.Version=${VERSION_SET}"

provider:: .provider gen_schema ## Builds provider.

gen_schema:: ## Generates schema.json.
	@ pulumi package get-schema bin/$(PROVIDER) > schema.json

read_version:: ## Reads plugin version from schema.
	@ jq '.version' schema.json

tests:: ## Runs unit tests.
	cd provider && go test -short -v -count=1 -cover -timeout 30m -parallel 1 ./...

acctests:: ## Runs acc tests.
	cd acc-test && go test -short -v -timeout 30m -parallel 1 ./...

lint:: ## Lints the provider's codebase.
	for DIR in "provider"; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

verify_version:: ## Checks that the schema version corresponds to release version.
	@ if [ "$(shell make read_version)" != "$(VERSION_SET)" ]; then echo inconsistent versions && exit 1; fi

sdk_go.local:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/go/go.* sdk-template/go/README.md ## Generates Go SDK.
	@ git submodule update --depth 1 --init --recursive --remote -f
	@ make sdk_go.ci

sdk_go:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/go/go.* sdk-template/go/README.md ## Generates Go SDK.
	@ rm -rf $(WORKING_DIR)/$(GO_SDK)/*
	@ cp -rf sdk-template/go/* $(WORKING_DIR)/$(GO_SDK)/ && cp -f LICENSE $(WORKING_DIR)/$(GO_SDK)/
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o $(WORKING_DIR)/$(GO_SDK) --language go
	@ cd $(WORKING_DIR)/$(GO_SDK) && mv go/$(GO_SDK)/* . && rm -r go
	@ cd $(WORKING_DIR)/$(GO_SDK) && go mod tidy

sdk_nodejs:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/nodejs/README.md ## Generates Node.js SDK.
	@ rm -rf sdk-nodejs
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-nodejs --language nodejs
	@ cd sdk-nodejs/nodejs/ && \
		npm install && \
		npm run build && \
		cp ../../sdk-template/nodejs/README.md bin/README.md && \
		cp ../../LICENSE package.json package-lock.json bin/ && \
		sed -i.bak 's/$${VERSION_SET}/$(VERSION_SET)/g' bin/package.json && \
		rm ./bin/package.json.bak
	@ mv sdk-nodejs/nodejs/bin/* sdk-nodejs/ && rm -r sdk-nodejs/nodejs

sdk_python:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/python/README.md ## Generates python SDK.
	@ rm -rf sdk-python
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-python --language python
	@ mv sdk-python/python/* sdk-python/ && rm -r sdk-python/python
	@ cp sdk-template/python/README.md sdk-python/README.md && \
 		cp LICENSE sdk-python/$(PY_PKG_NAME)/
	@ cd sdk-python && \
		python3 -m venv .venv && source .venv/bin/activate && pip install setuptools 2>&1 > /dev/null && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(VERSION_PY)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION_PY)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist 2>/dev/null

sdk_dotnet:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/dotnet/README.md ## Generates .Net SDK.
	@ rm -rf sdk-dotnet
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-dotnet --language dotnet
	@ cp fig/logo.png sdk-dotnet/dotnet/ && mv sdk-dotnet/dotnet/* sdk-dotnet && rm -r sdk-dotnet/dotnet
	@ cp sdk-template/dotnet/README.md sdk-dotnet/README.md
	@ cp sdk-template/dotnet/README.md sdk-dotnet/Config/README.md
	@ cd sdk-dotnet/ && \
		echo "${VERSION_SET}" >version.txt && \
		dotnet build /p:Version=${VERSION_SET}

sdk_java:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/java/README.md ## Generates Java SDK.
	@ rm -rf sdk-java
	@ pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) -o sdk-java --language java
	@ cd sdk-java && mv java/* . && rm -r java
	@ cp -r sdk-template/java/* sdk-java/
