/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flexvolume

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/huaweicloud/external-obs/pkg/provisioner/obs"
)

// FlexVolumeDriver defines
type FlexVolumeDriver struct {
	uuid string
}

// driverOp defines
type driverOp func(*FlexVolumeDriver, []string) (map[string]interface{}, error)

// cmdInfo defines
type cmdInfo struct {
	numArgs int
	run     driverOp
}

// commands defines
var commands = map[string]cmdInfo{
	"init": {
		0, func(d *FlexVolumeDriver, args []string) (map[string]interface{}, error) {
			return d.init()
		},
	},
	"mount": {
		2, func(d *FlexVolumeDriver, args []string) (map[string]interface{}, error) {
			return d.mount(args[0], args[1])
		},
	},
	"unmount": {
		1, func(d *FlexVolumeDriver, args []string) (map[string]interface{}, error) {
			return d.unmount(args[0])
		},
	},
}

// NewFlexVolumeDriver returns a flex volume driver
func NewFlexVolumeDriver(uuid string) *FlexVolumeDriver {
	return &FlexVolumeDriver{
		uuid: uuid,
	}
}

// Run is entrance for driver
func (d *FlexVolumeDriver) Run(args []string) string {
	return formatResult(d.doRun(args))
}

// doRun for driver
func (d *FlexVolumeDriver) doRun(args []string) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, errors.New("no arguments passed to flexvolume driver")
	}
	nArgs := len(args) - 1
	op := args[0]
	if cmdInfo, found := commands[op]; found {
		if cmdInfo.numArgs == nArgs {
			return cmdInfo.run(d, args[1:])
		} else {
			return nil, fmt.Errorf("unexpected number of args %d (expected %d) for operation %q", nArgs, cmdInfo.numArgs, op)
		}
	} else {
		return map[string]interface{}{
			"status": "Not supported",
		}, nil
	}
}

// init: <driver executable> init
func (d *FlexVolumeDriver) init() (map[string]interface{}, error) {
	glog.V(5).Info("flexvolume init() is called")

	// "{\"status\": \"Success\", \"capabilities\": {\"attach\": false}}"
	return map[string]interface{}{
		"capabilities": map[string]bool{
			"attach": false,
		},
	}, nil
}

// mount: <driver executable> mount <mount dir> <json options>
func (d *FlexVolumeDriver) mount(targetMountDir, jsonOptions string) (map[string]interface{}, error) {
	glog.V(5).Infof("flexvolume targetMountDir: %s jsonOptions: %s", targetMountDir, jsonOptions)

	// execIsMounted
	res, err := execIsMounted(targetMountDir)
	if err != nil {
		return nil, err
	}
	if res == "1" {
		glog.V(5).Infof("flexvolume targetMountDir: %s has already been mounted", targetMountDir)
		return nil, nil
	}

	// Unmarshal
	var volOptions map[string]interface{}
	json.Unmarshal([]byte(jsonOptions), &volOptions)

	// OBSAccessKey
	var strOBSAccessKey string
	if OBSAccessKey, ok := volOptions[obs.OBSAccessKey]; ok {
		strOBSAccessKey = OBSAccessKey.(string)
	}
	if strOBSAccessKey == "" {
		return nil, errors.New("OBSAccessKey is empty")
	}

	// OBSSecretKey
	var strOBSSecretKey string
	if OBSSecretKey, ok := volOptions[obs.OBSSecretKey]; ok {
		strOBSSecretKey = OBSSecretKey.(string)
	}
	if strOBSSecretKey == "" {
		return nil, errors.New("OBSSecretKey is empty")
	}

	// OBSBucket
	var strOBSBucket string
	if OBSBucket, ok := volOptions[obs.OBSBucket]; ok {
		strOBSBucket = OBSBucket.(string)
	}
	if strOBSBucket == "" {
		return nil, errors.New("OBSBucket is empty")
	}

	// OBSEndpoint
	var strOBSEndpoint string
	if OBSEndpoint, ok := volOptions[obs.OBSEndpoint]; ok {
		strOBSEndpoint = OBSEndpoint.(string)
	}
	if strOBSEndpoint == "" {
		return nil, errors.New("OBSEndpoint is empty")
	}

	// execMount
	err = execMount(targetMountDir, strOBSAccessKey, strOBSSecretKey, strOBSBucket, strOBSEndpoint)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// unmount: <driver executable> unmount <mount dir>
func (d *FlexVolumeDriver) unmount(targetMountDir string) (map[string]interface{}, error) {
	glog.V(5).Infof("flexvolume targetMountDir: %s", targetMountDir)

	// check the target directory
	if _, err := os.Stat(targetMountDir); os.IsNotExist(err) {
		glog.V(5).Infof("flexvolume targetMountDir: %s has does not exist", targetMountDir)
		return nil, nil
	}

	// execIsMounted
	res, err := execIsMounted(targetMountDir)
	if err != nil {
		return nil, err
	}
	if res == "0" {
		glog.V(5).Infof("flexvolume targetMountDir: %s has already been unmounted", targetMountDir)
		return nil, nil
	}

	// execUnmount
	err = execUnmount(targetMountDir)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
