#!/usr/bin/env bash
SCRIPT_DIR=`dirname $0`

oc create -f ${SCRIPT_DIR}/projects.yaml

for i in 3scale-operator 3scale-apimanager api-dev api-uat api-prod
do
  oc -n $i create secret generic redhat-registry --from-file=.dockerconfigjson=$HOME/.docker/config.json --type=kubernetes.io/dockerconfigjson
  oc -n $i secrets link default redhat-registry --for=pull
done

for i in `ls deploy/crds/*3scale.hosted*_crd.yaml`; do oc create -f $i ; done

oc -n 3scale-operator create -f ${SCRIPT_DIR}/../deploy/service_account.yaml
oc create -f ${SCRIPT_DIR}/../deploy/cluster_role.yaml
oc create -f ${SCRIPT_DIR}/../deploy/cluster_role_binding.yaml

oc -n 3scale-operator create -f ${SCRIPT_DIR}/../deploy/operator.yaml
oc create -f ${SCRIPT_DIR}/api-manager.yaml
