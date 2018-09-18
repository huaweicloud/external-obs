## Deploy obs-flexvolume in openshift

```
oc adm policy add-scc-to-user privileged system:serviceaccount:default:obs-flexvolume
oc create -f https://raw.githubusercontent.com/huaweicloud/external-obs/master/deploy/obs-flexvolume/openshift/daemonset.yaml
```
