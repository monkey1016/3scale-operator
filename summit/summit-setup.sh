#!/usr/bin/env bash

oc create -f summit/projects.yaml

for i in 3scale-operator 3scale-apimanager api-dev api-uat api-prod
do
  oc -n $i create secret generic redhat-registry --from-file=.dockerconfigjson=$HOME/.docker/config.json --type=kubernetes.io/dockerconfigjson
  oc -n $i secrets link default redhat-registry --for=pull
done

for i in `ls deploy/crds/*3scale.hosted*_crd.yaml`; do oc create -f $i ; done

oc -n 3scale-operator create -f deploy/service_account.yaml
oc create -f deploy/cluster_role.yaml
oc create -f deploy/cluster_role_binding.yaml

oc -n 3scale-operator create -f deploy/operator.yaml
oc create -f summit/api-manager.yaml