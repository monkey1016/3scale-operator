#!/usr/bin/env bash
set -x

oc new-project 3scale-operator-default

oc -n 3scale-operator-default create secret generic redhat-registry --from-file=.dockerconfigjson=$HOME/.docker/config.json --type=kubernetes.io/dockerconfigjson
oc -n 3scale-operator-default secrets link default redhat-registry --for=pull

for i in `ls deploy/crds/*_crd.yaml`; do oc create -f $i ; done

oc -n 3scale-operator-default create -f deploy/service_account.yaml
oc create -n 3scale-operator-default -f deploy/role.yaml
oc create -n 3scale-operator-default -f deploy/role_binding.yaml

oc -n 3scale-operator-default create -f deploy/operator.yaml
oc create -f summit/api-manager.yaml

# Recreate the storage as ReadWriteOnce since the cluster doesn't support ReadWriteMany
oc delete -f summit/system-storage.yaml
oc create -f summit/system-storage.yaml