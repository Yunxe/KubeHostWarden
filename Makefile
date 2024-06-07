backend := $(shell pwd)/backend

.PHONY: build-ops-image
build-ops-image:
	docker build -t opscenter:latest -f backend/deploy/opscenter/dockerfile ./backend

.PHONY: load-host-image
build-host-image:
	docker build -t host:latest -f backend/deploy/host/dockerfile ./backend
	kind load docker-image host:latest

.PHONY: ops
ops:
	bash update_env.sh
	cd $(backend) && \
	go build -o build/opscenter cmd/opscenter/opscenter.go && \
	./build/opscenter -env ./.env

.PHONY: host
host:
	cd $(backend) && \
	go build -o build/host cmd/host/host.go && \
	./build/host -env ./.env

.PHONY: host-local
host-local:
	bash update_env.sh
	cd $(backend) && \
	go build -o build/host cmd/host/host.go && \
	./build/host -env ./.env