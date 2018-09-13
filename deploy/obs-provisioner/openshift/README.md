## Deploy obs-provisioner in openshift

```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:obs-provisioner
oc create -f https://raw.githubusercontent.com/huaweicloud/external-obs/master/deploy/obs-provisioner/openshift/statefulset.yaml
```
