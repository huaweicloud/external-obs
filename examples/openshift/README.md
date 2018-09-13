## Usage external-obs in openshift

1. Create a storage class named ```obs-storage-class```.

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-obs/master/examples/openshift/sc.yaml
```

2. Create a obs pvc named ```obs-pvc```.

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-obs/master/examples/openshift/pvc.yaml
```

3. Create a nginx pod with obs pvc.

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-obs/master/examples/openshift/pod.yaml
```

If you want to create all of the above resources, you could run:

```
oc create -f https://raw.githubusercontent.com/huaweicloud/external-obs/master/examples/openshift/example.yaml
```
