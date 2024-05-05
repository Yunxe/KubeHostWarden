# KubeHostWarden

> An advanced version of an APM-Server project, which also serves as graduation project.

## Intro
A system which automatically collects metrics from host machine and saves them to Influxdb.

## Run
### Build images
1. ```make build-ops-image```
2. ```make build-host-image```

### Load images to your cluster
1. ```kind load docker-image opscenter:latest```
2. ```kind load docker-image host:latest```

### Deploy
1. ```kubectl apply -f backend/deploy/opscenter/```