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

	return &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        volOptions.PVName,
			Annotations: map[string]string{},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: volOptions.PersistentVolumeReclaimPolicy,
			AccessModes:                   volOptions.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): volOptions.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{},
		},
	}, nil
}

// Delete a bucket from obs
func (p *Provisioner) Delete(pv *v1.PersistentVolume) error {
	return nil
}
