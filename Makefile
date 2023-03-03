export CGO_ENABLED=0
export GO111MODULE=on

LDFLAGS := -w -s

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/kosli-dev/cli/internal/version.version=${BINARY_VERSION}
endif

VERSION_METADATA = unreleased
# Clear the "unreleased" string in BuildMetadata
ifneq ($(GIT_TAG),)
	VERSION_METADATA =
endif

LDFLAGS += -X github.com/kosli-dev/cli/internal/version.metadata=${VERSION_METADATA}
LDFLAGS += -X github.com/kosli-dev/cli/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X github.com/kosli-dev/cli/internal/version.gitTreeState=${GIT_DIRTY}
LDFLAGS += -extldflags "-static"

ldflags:
	@echo $(LDFLAGS)

fmt: ## Reformat package sources
	@go fmt ./...

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.50-alpine golangci-lint run --timeout=5m  -v ./...

vet: fmt
	@go vet ./...

deps: ## Install depdendencies. Runs `go get` internally.
	@GOFLAGS="" go mod download
	@GOFLAGS="" go mod tidy

build: deps vet ## Build the binary
	@go build -o kosli -ldflags '$(LDFLAGS)' ./cmd/kosli/

check_dirty:
	@git diff-index --quiet HEAD --  || echo "Cannot test release with dirty git repo"
	@git diff-index --quiet HEAD -- 

add_test_tag:
	@git tag -d v0.0.99 2> /dev/null || true
	@git tag v0.0.99

build_release: check_dirty add_test_tag
	rm -rf dist/
	goreleaser release --skip-publish --debug
	@git tag -d v0.0.99 2> /dev/null || true

ensure_network:
	docker network inspect cli_net > /dev/null || docker network create --driver bridge cli_net

test_integration_setup:
	./bin/docker_login_aws.sh staging
	@docker-compose down || true
	@docker pull 772819027869.dkr.ecr.eu-central-1.amazonaws.com/merkely:latest || true
	@docker-compose up -d
	./mongo/ip_wait.sh localhost:8001
	@docker exec cli_kosli_server /demo/create_test_users.py
	@go install gotest.tools/gotestsum@latest


test_integration: deps vet ensure_network test_integration_setup ## Run tests except the too slow ones
	@export KOSLI_TESTS=true && ~/go/bin/gotestsum -- --short -p=8 -coverprofile=cover.out ./...
	@go tool cover -func=cover.out | grep total:
	@go tool cover -html=cover.out


test_integration_full: deps vet ensure_network test_integration_setup ## Run all tests
	@export KOSLI_TESTS=true && ~/go/bin/gotestsum -- -p=8 -coverprofile=cover.out ./...
	@go tool cover -func=cover.out


test_integration_no_setup: 
	@~/go/bin/gotestsum -- --short -p=8 -coverprofile=cover.out ./...
	@go tool cover -html=cover.out


test_integration_single: test_integration_setup
	@export KOSLI_TESTS=true && ~/go/bin/gotestsum -- -p=4 ./... -run "${TARGET}"


test_docs: deps vet ensure_network test_integration_setup
	./bin/test_docs_cmds.sh docs.kosli.com/content/use_cases/simulating_a_devops_system/_index.md


docker: deps vet lint
	@docker build -t kosli-cli .

docs: helm-docs build
	@rm -f docs.kosli.com/content/client_reference/kosli*
	@export DOCS=true && ./kosli docs --dir docs.kosli.com/content/client_reference

licenses:
	@rm -rf licenses || true
	@go install github.com/google/go-licenses@latest
	@go-licenses save ./... --save_path="licenses/" || true
	$(eval DATA := $(shell go-licenses csv ./...))
	@echo $(DATA) | tr " " "\n" > licenses/licenses.csv

generate-json-metadata:
	echo '{"currentversion": "local"}' > docs.kosli.com/assets/metadata.json

hugo: docs generate-json-metadata
	cd docs.kosli.com && hugo server --minify --buildDrafts --port=1515

hugo-local: generate-json-metadata
	cd docs.kosli.com && hugo server --minify --buildDrafts --port=1515

helm-lint: 
	@cd charts/k8s-reporter && helm lint .

helm-docs: helm-lint
	@cd charts/k8s-reporter &&  docker run --rm --volume "$(PWD):/helm-docs" jnorwood/helm-docs:latest --template-files README.md.gotmpl,_templates.gotmpl --output-file README.md
	@cd charts/k8s-reporter &&  docker run --rm --volume "$(PWD):/helm-docs" jnorwood/helm-docs:latest --template-files README.md.gotmpl,_templates.gotmpl --output-file ../../docs.kosli.com/content/helm/_index.md

release:
	@git remote update
	@git status -uno | grep --silent "Your branch is up to date" || (echo "ERROR: your branch is NOT up to date with remote" && return 1)
	git tag -a $(tag) -m"$(tag)"
	git push origin $(tag)

# check-links:
# 	@docker run -v ${PWD}:/tmp:ro --rm -i --entrypoint '' ghcr.io/tcort/markdown-link-check:stable /bin/sh -c 'find /tmp/docs.kosli.com/content -name \*.md -print0 | xargs -0 -n1 markdown-link-check -q -c /tmp/link-checker-config.json'

check-links: 
	@cd docs.kosli.com && hugo --minify
	@docker run -v ${PWD}:/test --rm wjdp/htmltest -c .htmltest.yml -l 1
