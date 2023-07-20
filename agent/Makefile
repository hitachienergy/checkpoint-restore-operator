VERSION ?= 0.0.1
IMAGE_TAG_BASE ?= $(IMAGE_REGISTRY)/cr-operator/agent
IMG ?= $(IMAGE_TAG_BASE):v$(VERSION)
ARCH ?= amd64

.PHONY: generate-proto
generate-proto:
	protoc -I=. --go_out=. --go-grpc_out=. api/*.proto

.PHONY: generate
generate: generate-proto


.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet -tags "containers_image_openpgp" ./...

.PHONY: test
test: generate fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: generate fmt vet ## Build manager binary.
	go build -o bin/agent -tags "containers_image_openpgp" main.go

.PHONY: run
run: generate fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build
docker-build: test ## Build docker image with the manager.
	docker buildx build --platform linux/$(ARCH) --build-arg ARCH=$(ARCH) -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

.PHONY: deploy
deploy: ## Push docker image with the manager.
	kubectl apply -f agent.daemonset.yaml