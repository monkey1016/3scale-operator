# Notes for Red Hat Summit 3Scale

## Building Operator Image
This requires an account on quay.io
1. `docker login quay.io` and enter credentials
2. `make build IMAGE=quay.io/kjanania/3scale-bof-summit-2020 VERSION=summit-2.7`
3. `make push IMAGE=quay.io/kjanania/3scale-bof-summit-2020 VERSION=summit-2.7`
   (Careful, could take 15 or so minutes for it to become available)


## Launching a Cluster
1. Assign 6 CPUs and at least 18GB or RAM to CodeReady Containers
2. `crc setup`
3. `crc start`
4. `oc new-project 3scale`
4. `oc create secret generic redhat-registry --from-file=.dockerconfigjson=$HOME/.docker/config.json --type=kubernetes.io/dockerconfigjson`
5. `oc secrets link default redhat-registry --for=pull`
6. ```for i in `ls deploy/crds/*_crd.yaml`; do oc create -f $i ; done```
7. `oc create -f deploy/service_account.yaml`
8. `oc create -f deploy/cluster_role.yaml`
9. `oc create -f deploy/cluster_role_binding.yaml`
10. `sed -i 's|REPLACE_IMAGE|quay.io/kjanania/3scale-bof-summit-2020:summit-2.7|g' deploy/operator.yaml`
11. `oc create -f deploy/operator.yaml`
12. `oc create -f summit/api-manager.yaml`
13. Log in to https://master.apps-crc.testing/
14. Use `system-seed` secret with `MASTER_USER` and `MASTER_PASSWORD` as credentials

If you can log in, it was successful

## Creating an API
```for i in `ls summit/demo-app/*.yaml`; do oc create -f $i; done```


## Cleaning Up
```for i in `ls -r summit/demo-app/*.yaml`; do oc delete -f $i; done```
