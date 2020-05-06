#!/usr/bin/env bash

oc new-project 3scale-operator-default

oc -n 3scale-operator-default create secret generic redhat-registry --from-file=.dockerconfigjson=$HOME/.docker/config.json --type=kubernetes.io/dockerconfigjson
oc -n 3scale-operator-default secrets link default redhat-registry --for=pull

for i in `ls deploy/crds/*3scale.hosted*_crd.yaml`; do oc create -f $i ; done

oc -n 3scale-operator-default create -f deploy/service_account.yaml
oc create -n 3scale-operator-default -f deploy/role.yaml
oc create -n 3scale-operator-default-f deploy/role_binding.yaml

oc -n 3scale-operator-default create -f deploy/operator.yaml
oc create -f summit/api-manager.yaml