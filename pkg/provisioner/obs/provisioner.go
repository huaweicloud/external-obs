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
	"github.com/huaweicloud/external-obs/pkg/provisioner/config"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// Provisioner implements controller.Provisioner interface
type Provisioner struct {
	clientset   clientset.Interface
	cloudconfig config.CloudCredentials
}

// NewProvisioner creates a new instance of obs provisioner
func NewProvisioner(c clientset.Interface, cc config.CloudCredentials) *Provisioner {

	// return provisioner instance
	return &Provisioner{
		clientset:   c,
		cloudconfig: cc,
	}
}

// Provision a bucket in obs
func (p *Provisioner) Provision(volOptions controller.VolumeOptions) (*v1.PersistentVolume, error) {

	// selector check
	glog.Infof("Provision volOptions: %v", volOptions)
	if volOptions.PVC.Spec.Selector != nil {
		return nil, fmt.Errorf("Claim Selector is not supported")
	}

	// get ak
	ak := volOptions.Parameters[OBSParametersAccessKey]
	if ak == "" {
		return nil, fmt.Errorf("%s is not set", OBSParametersAccessKey)
	}

	// get sk
	sk := volOptions.Parameters[OBSParametersSecretKey]
	if sk == "" {
		return nil, fmt.Errorf("%s is not set", OBSParametersSecretKey)
	}

	// init obs client
	glog.Info("Init obs client...")
	client, err := p.cloudconfig.OBSClient(ak, sk)
	if err != nil {
		return nil, fmt.Errorf("Failed to create obs client: %v", err)
	}

	// close obs client
	if client != nil {
		defer client.Close()
	}

	// create bucket
	glog.Info("Create bucket begin...")
	bucket, err := CreateBucket(client, &volOptions, p)
	if err != nil {
		return nil, fmt.Errorf("Failed to create bucket: %v", err)
	}

	// Example: https://{BucketName}.{Endpoint}
	endpoint := client.GetEndpoint()
	glog.Infof("Provision endpoint: %s", endpoint)

	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: volOptions.PVName,
			Annotations: map[string]string{
				OBSAnnotationID: *bucket,
				OBSAnnotationAK: ak,
				OBSAnnotationSK: sk,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: volOptions.PersistentVolumeReclaimPolicy,
			AccessModes:                   volOptions.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): volOptions.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				FlexVolume: &v1.FlexVolumeSource{
					Driver: OBSFlexVolume,
					Options: map[string]string{
						OBSBucket:    *bucket,
						OBSAccessKey: ak,
						OBSSecretKey: sk,
						OBSEndpoint:  endpoint,
					},
					ReadOnly: false,
				},
			},
		},
	}, nil
}

// Delete a bucket from obs
func (p *Provisioner) Delete(pv *v1.PersistentVolume) error {

	// get ak
	ak, ok := pv.ObjectMeta.Annotations[OBSAnnotationAK]
	if (!ok) || (ak == "") {
		return fmt.Errorf("Failed to get ak %v", pv)
	}

	// get sk
	sk, ok := pv.ObjectMeta.Annotations[OBSAnnotationSK]
	if (!ok) || (sk == "") {
		return fmt.Errorf("Failed to get sk %v", pv)
	}

	// init obs client
	glog.Info("Init obs client...")
	client, err := p.cloudconfig.OBSClient(ak, sk)
	if err != nil {
		return fmt.Errorf("Failed to create obs client: %v", err)
	}

	// close obs client
	if client != nil {
		defer client.Close()
	}

	// get bucket
	bucket, ok := pv.ObjectMeta.Annotations[OBSAnnotationID]
	if (!ok) || (bucket == "") {
		return fmt.Errorf("Failed to get bucket %v", pv)
	}

	// delete objects in bucket
	glog.Infof("Delete objects in bucket: %s", bucket)
	err = DeleteObjects(client, bucket)
	if err != nil {
		return fmt.Errorf("failed to delete objects in bucket: %v", err)
	}

	// delete bucket
	glog.Infof("Delete bucket: %s", bucket)
	err = DeleteBucket(client, bucket)
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %v", err)
	}

	return nil
}
