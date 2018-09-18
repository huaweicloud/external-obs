# external-obs
[![Go Report Card](https://goreportcard.com/badge/github.com/huaweicloud/external-obs)](https://goreportcard.com/badge/github.com/huaweicloud/external-obs)
[![Build Status](https://travis-ci.org/huaweicloud/external-obs.svg?branch=master)](https://travis-ci.org/huaweicloud/external-obs)
[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/external-obs/blob/master/LICENSE)

Object Storage Service (OBS) is a stable, secure, efficient,
and easy-to-use cloud storage service on huawei clouds.
With Representational State Transfer Application Programming Interfaces (REST APIs),
OBS is able to store unstructured data of any amount and form at 99.999999999% reliability.

This repository houses external obs provisioner and flexvolume for OpenShift.


## Getting Started on OpenShift

external-obs should be deployed in the OpenShift Master after OpenShift is deployed successfully. Please firstly run the following command to download this repository,
```
git clone https://github.com/huaweicloud/external-obs
```

### Deploy obs-provisioner

In default, the Cloud Tenant informations are stored in the file ```/etc/origin/cloudprovider/openstack.conf``` of OpenShift Master. If your OpenShift Master does not contain the file ```/etc/origin/cloudprovider/openstack.conf```, please modify the statefulset.yaml,
```
vi external-obs/deploy/obs-provisioner/openshift/statefulset.yaml
```
and replace ```/etc/origin/cloudprovider/openstack.conf``` with your Cloud Config file in the line 73 of statefulset.yaml and replace the path ```/etc/origin``` with your Cloud Config directory in the line 80 of statefulset.yaml,

if you want to increase the log level, please add the following two lines after the line 73 of statefulset.yaml.

```
            - name: OS_DEBUG
              value: true
```

finally you can run the following command.
```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:obs-provisioner
oc create -f external-obs/deploy/obs-provisioner/openshift/statefulset.yaml
```

### Deploy obs-flexvolume

In default, the flexvolume plugins are stored in the folder ```/usr/libexec/kubernetes/kubelet-plugins/volume/exec``` of OpenShift Cluster. If your OpenShift Cluster stores the flexvolume plugins in the other folder, please modify the daemonset.yaml,
```
vi external-obs/deploy/obs-flexvolume/openshift/daemonset.yaml
```
and replace ```/usr/libexec/kubernetes/kubelet-plugins/volume/exec``` with your flexvolume plugins folder in the line 74 of daemonset.yaml, finally you can run the following command.
```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:obs-flexvolume
oc create -f external-obs/deploy/obs-flexvolume/openshift/daemonset.yaml
```
Actually after the daemonset is running, it means the obs-flexvolume has already been deployed in the OpenShift Cluster. If you do not want the daemonset, you can delete it by the following command.
```
oc delete -f external-obs/deploy/obs-flexvolume/openshift/daemonset.yaml
```

### Usage

Before you start to use the example, please modify the example.yaml,
```
vi external-obs/examples/openshift/example.yaml
```
and replace ```ak``` and  ```sk``` with your cloud account access key and secret key in the line 7 and 8 of example.yaml,
 finally you can run the following example.
```
oc create -f external-obs/examples/openshift/example.yaml
```

## License

See the [LICENSE](LICENSE) file for details.
