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
VERSION_SDK 	:= $(shell jq '.version' schema.json | sed 's/"//g')

.PHONY: help
help: ## Prints help message.
	@ grep -h -E '^[a-zA-Z0-9_-].+::.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1m%-30s\033[0m %s\n", $$1, $$2}'

.provider:
	@ CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/provider.Version=${VERSION_SET}"

provider:: .provider gen_schema ## Builds provider.

gen_schema:: ## Generates schema.json.
	@ pulumi package get-schema bin/$(PROVIDER) > schema.json

read_version:: ## Reads plugin version from schema.
	@ jq '.version' schema.json | sed 's/"//g'

tests:: ## Runs unit tests.
	cd provider && go test -short -v -count=1 -cover -timeout 30m -parallel 1 ./...

acctests:: ## Runs acc tests.
	cd acc-test && go test -short -v -timeout 30m -parallel 1 ./...

lint:: ## Lints the provider's codebase.
	for DIR in "provider"; do \
		pushd $$DIR && golangci-lint run -c ../.golangci.yml --timeout 10m && popd ; \
	done

verify_version:: VERSION_PROVIDER := $(shell jq '.version' schema.json)
verify_version:: ## Checks that the schema version corresponds to release version.
	@ if [ "$(VERSION_PROVIDER)" != $(VERSION_SET) ]; then echo inconsistent version: schema.json - $(VERSION_PROVIDER), variable - $(VERSION_SET) && exit 1; fi

sdk_go.local:: $(WORKING_DIR)/bin/$(PROVIDER) schema.json sdk-template/go/go.* sdk-template/go/README.md ## Generates Go SDK.
	@ git submodule update --depth 1 --init --recursive --remote -f
	@ make sdk_go.ci

sdk_go:: schema.json sdk-template/go/go.* sdk-template/go/README.md ## Generates Go SDK.
	@ rm -rf $(WORKING_DIR)/$(GO_SDK)/*
	@ cp -rf sdk-template/go/* $(WORKING_DIR)/$(GO_SDK)/ && cp -f LICENSE $(WORKING_DIR)/$(GO_SDK)/
	@ pulumi package gen-sdk schema.json -o $(WORKING_DIR)/$(GO_SDK) --language go
	@ cd $(WORKING_DIR)/$(GO_SDK) && mv go/$(GO_SDK)/* . && rm -r go

sdk_go.publish:: VERSION_SET := $(shell jq '.version' $(WORKING_DIR)/$(GO_SDK)/pulumi-plugin.json | sed 's/"//g')
sdk_go.publish:: verify_version ## Publishes Go SDK to github.
	@ cd $(WORKING_DIR)/$(GO_SDK) && \
 		git add --all . && \
 		git commit -S -m "release v$(VERSION_SET)" && \
 		git push origin master -f && \
 		git tag v$(VERSION_SET) && \
 		git push origin v$(VERSION_SET)

sdk_nodejs:: schema.json sdk-template/nodejs/README.md ## Generates Node.js SDK.
	@ rm -rf sdk/nodejs
	@ pulumi package gen-sdk schema.json --language nodejs
	@ cd sdk/nodejs && \
		npm install && \
		npm run build && \
		cp ../../sdk-template/nodejs/README.md bin/README.md && \
		cp ../../LICENSE package.json package-lock.json bin/ && \
		sed -i '' 's/$${VERSION_SDK}/$(VERSION_SDK)/g' bin/package.json && \
		mv bin ../sdk-nodejs-temp && rm -rf * && mv ../sdk-nodejs-temp/* . && rm -r ../sdk-nodejs-temp

sdk_nodejs.publish:: VERSION_SET := $(shell jq '.version' sdk/nodejs/package.json)
sdk_nodejs.publish:: verify_version ## Publishes Node.js SDK to npm.
	@ if [ -z "${NPM_TOKEN}" ]; then echo "env vriable NPM_TOKEN must be set" && exit 1 ; fi
	@ cd sdk/nodejs && \
 		npm set "//registry.npmjs.org/:_authToken=$${NPM_TOKEN}" && \
 		if [ -z "${GITHUB_ACTION}" ]; then npm publish --access public; else npm publish --access public --provenance; fi

sdk_python:: schema.json sdk-template/python/README.md ## Generates python SDK.
	@ rm -rf sdk/python
	@ pulumi package gen-sdk schema.json --language python
	@ cp sdk-template/python/README.md sdk/python/README.md && cp LICENSE sdk-python/$(PY_PKG_NAME)/

sdk_python.build:: VERSION_SET := $(shell jq '.version' $(WORKING_DIR)/sdk/python/$(PY_PKG_NAME)/pulumi-plugin.json | sed 's/"//g')
sdk_python.build:: verify_version ## Builds python SDK dist.
	@ cd sdk/python && \
		python3 -m venv .venv && source .venv/bin/activate && pip install setuptools 1> /dev/null && \
		python3 setup.py clean --all 1>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(VERSION_SET)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION_SET)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist 1>/dev/null

sdk_python.publish:: $(WORKING_DIR)/sdk/python/bin/dist/*.tar.gz  ## Publishes python SDK to PyPi.
	@ if [ -z "${PYPI_TOKEN}" ]; then echo "PYPI_TOKEN env variable must be set"; exit 1; fi
	@ cd sdk/python && \
      	python3 -m venv .venv && source .venv/bin/activate && pip install twine && \
      	twine upload -u "__token__" -p "${PYPI_TOKEN}" $(WORKING_DIR)/sdk/python/bin/dist/* --skip-existing --verbose

sdk_dotnet:: schema.json sdk-template/dotnet/README.md ## Generates .Net SDK.
	@ rm -rf sdk-dotnet
	@ pulumi package gen-sdk schema.json -o sdk-dotnet --language dotnet
	@ cp fig/logo.png sdk-dotnet/dotnet/ && mv sdk-dotnet/dotnet/* sdk-dotnet && rm -r sdk-dotnet/dotnet
	@ echo $(VERSION_SDK) > sdk-dotnet/version.txt
	@ cp sdk-template/dotnet/README.md sdk-dotnet/README.md && \
		cd sdk-dotnet && sed -i '' -e 's|</Project>|  <PropertyGroup>\n    <PackageReadmeFile>README.md</PackageReadmeFile>\n  </PropertyGroup>\n  <ItemGroup>\n    <None Include="README.md" Pack="True" PackagePath="/" />\n  </ItemGroup>\n</Project>|g' *.csproj
	@ cd sdk-dotnet && dotnet build --verbosity minimal 1>/dev/null

sdk_java:: schema.json sdk-template/java/README.md sdk-template/java/build.gradle sdk-template/java/settings.gradle ## Generates Java SDK.
	@ rm -rf sdk-java
	@ pulumi package gen-sdk schema.json -o sdk-java --language java
	@ cd sdk-java && mv java/* . && rm -r java
	@ cp -r sdk-template/java/* sdk-java/
	@ # fix for https://github.com/pulumi/pulumi-java/issues/1534
	@ cd sdk-java/src && find . -type f -name '*.java' | xargs sed -i '' -e 's/import javax.annotation/import jakarta.annotation/g'
	@ cd sdk-java && gradle build
