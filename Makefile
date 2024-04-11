default: run

run:
	go run backend/cmd/opscenter/opscenter.go

.PHONY: host
host:
	go run backend/cmd/host/host.go