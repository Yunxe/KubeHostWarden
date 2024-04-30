backend := $(shell pwd)/backend

.PHONY: build-ops-image
build-ops-image:
	docker build -t opscenter:latest -f backend/deploy/opscenter/dockerfile ./backend

.PHONY: build-host-image
build-host-image:
	docker build -t host:latest -f backend/deploy/host/dockerfile ./backend

.PHONY: ops
ops:
	cd $(backend) && \
	go build -o build/opscenter cmd/opscenter/opscenter.go && \
	./build/opscenter -env ./.env

.PHONY: host
host:
	cd $(backend) && \
	go build -o build/host cmd/host/host.go && \
	./build/host -env ./.env