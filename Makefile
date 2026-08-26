REGISTRY := registry.gitlab.com
IMAGE_BASE := constacode/menunderfireapp
TAG := $(shell git rev-parse --short HEAD)

.PHONY: build-frontend build-backend build-all push-frontend push-backend push-all

build-frontend:
	docker build \
		-t $(REGISTRY)/$(IMAGE_BASE)/frontend:$(TAG) \
		-t $(REGISTRY)/$(IMAGE_BASE)/frontend:latest \
		-f frontend/Dockerfile \
		frontend/

build-backend:
	docker build \
		-t $(REGISTRY)/$(IMAGE_BASE)/backend:$(TAG) \
		-t $(REGISTRY)/$(IMAGE_BASE)/backend:latest \
		-f backend_go/Dockerfile \
		backend_go/

build-all: build-frontend build-backend

push-frontend: build-frontend
	docker push $(REGISTRY)/$(IMAGE_BASE)/frontend:$(TAG)
	docker push $(REGISTRY)/$(IMAGE_BASE)/frontend:latest

push-backend: build-backend
	docker push $(REGISTRY)/$(IMAGE_BASE)/backend:$(TAG)
	docker push $(REGISTRY)/$(IMAGE_BASE)/backend:latest

push-all: push-frontend push-backend
