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

package obs

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/huaweicloud/golangsdk/openstack/obs"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
)

// CreateBucket in OBS
func CreateBucket(client *obs.ObsClient, volOptions *controller.VolumeOptions, p *Provisioner) (*string, error) {

	// build share createOpts
	createOpts := &obs.CreateBucketInput{}

	// Setting Bucket
	createOpts.Bucket = "pvc-" + string(volOptions.PVC.GetUID())

	// Setting Default StorageClass
	sc := volOptions.Parameters[OBSParametersStorageClass]
	if sc == OBSParametersStorageClassStandard {
		createOpts.StorageClass = obs.StorageClassStandard
	} else if sc == OBSParametersStorageClassInfrequentAccess {
		createOpts.StorageClass = obs.StorageClassWarm
	} else if sc == OBSParametersStorageClassArchive {
		createOpts.StorageClass = obs.StorageClassCold
	} else {
		createOpts.StorageClass = obs.StorageClassStandard
	}

	// Setting ACL: private, public-read, public-read-write
	acl := volOptions.Parameters[OBSParametersBucketPolicy]
	if acl == "" {
		createOpts.ACL = obs.AclPrivate
	} else {
		createOpts.ACL = obs.AclType(acl)
	}

	// Setting Location
	createOpts.Location = p.cloudconfig.Global.Region

	// create bucket
	glog.Infof("Create bucket createOpts: %v", createOpts)
	obsResponse, err := client.CreateBucket(createOpts)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create bucket in OBS: %v", err)
	}

	glog.Infof("Create bucket response: %v", obsResponse)
	return &createOpts.Bucket, nil
}

// DeleteBucket in OBS
func DeleteBucket(client *obs.ObsClient, bucket string) error {

	obsResponse, err := client.DeleteBucket(bucket)
	if err != nil {
		return fmt.Errorf("Couldn't delete bucket in OBS: %v", err)
	}

	glog.Infof("Delete bucket response: %v", obsResponse)
	return nil
}
