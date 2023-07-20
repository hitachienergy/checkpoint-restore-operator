VERSION=latest-amd64
REGISTRY=${IMAGE_REGISTRY}
IMAGE_NAME=hitachienergy/the-count

FULL_IMAGE_NAME= $(REGISTRY)/$(IMAGE_NAME):$(VERSION)

build:
	docker build -t $(FULL_IMAGE_NAME) .

push:
	docker push $(FULL_IMAGE_NAME)

deploy:
	kubectl apply -f test-app-deployment.yaml

undeploy:
	kubectl delete -f test-app-deployment.yaml