.PHONY: all build obs-provisioner clean

all:build

build:obs-provisioner

package:
	mkdir -p ./bin/

obs-provisioner:package
	go build -o ./bin/obs-provisioner ./cmd/obs-provisioner

clean:
	rm -rf ./bin/ 
