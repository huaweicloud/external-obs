.PHONY: all build obs-provisioner docker clean

all:build

build:obs-provisioner

package:
	mkdir -p ./bin/

obs-provisioner:package
	go build -o ./bin/obs-provisioner ./cmd/obs-provisioner

docker:obs-provisioner
	cp ./bin/obs-provisioner ./cmd/obs-provisioner
	docker build cmd/obs-provisioner -t quay.io/huaweicloud/obs-provisioner:latest

clean:
	rm -rf ./bin/
	rm -rf ./cmd/obs-provisioner/obs-provisioner
