backend := $(shell pwd)/backend

.PHONY: ops dev
ops dev:
	cd $(backend) && \
	go build -o build/opscenter cmd/opscenter/opscenter.go && \
	./build/opscenter -env ./.env -e "dev"

.PHONY: ops
ops:
	cd $(backend) && \
	go build -o build/opscenter cmd/opscenter/opscenter.go && \
	./build/opscenter -env ./.env