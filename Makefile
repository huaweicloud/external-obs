.PHONY: all build obs-provisioner obs-flexvolume docker-obs-provisioner docker-obs-flexvolume docker-obs-mountvolume docker clean

all:build

build:obs-provisioner obs-flexvolume

package:
	mkdir -p ./bin/

obs-provisioner:package
	go build -o ./bin/obs-provisioner ./cmd/obs-provisioner

obs-flexvolume:package
	go build -o ./bin/obs-flexvolume ./cmd/obs-flexvolume

docker-obs-provisioner:obs-provisioner
	cp ./bin/obs-provisioner ./cmd/obs-provisioner
	docker build cmd/obs-provisioner -t quay.io/huaweicloud/obs-provisioner:latest

docker-obs-flexvolume:
	cp ./bin/obs-flexvolume ./cmd/obs-flexvolume
	docker build cmd/obs-flexvolume -t quay.io/huaweicloud/obs-flexvolume:latest

docker-obs-mountvolume:
	docker build cmd/obs-mountvolume -t quay.io/huaweicloud/obs-mountvolume:latest

docker:docker-obs-provisioner docker-obs-flexvolume docker-obs-mountvolume

clean:
	rm -rf ./bin/
	rm -rf ./cmd/obs-provisioner/obs-provisioner
	rm -rf ./cmd/obs-provisioner/obs-flexvolume
