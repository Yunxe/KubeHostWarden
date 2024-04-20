backend := $(shell pwd)/backend

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