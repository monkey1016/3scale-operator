#!/usr/bin/env bash

for i in `oc get -n 3scale-apimanager bindings.capabilities.3scale.hosted.net -o name`
do
  oc delete $i;
done

for i in `oc get -n 3scale-apimanager apis.capabilities.3scale.hosted.net -o name`
do
  oc delete $i;
done

for i in `oc get -n 3scale-apimanager metrics.capabilities.3scale.hosted.net -o name`
do
  oc delete $i;
done

for i in `oc get -n 3scale-apimanager plans.capabilities.3scale.hosted.net -o name`
do
  oc delete $i;
done

for i in `oc get -n 3scale-apimanager mappingrules.capabilities.3scale.hosted.net -o name`
do
  oc delete $i;
done

oc delete -f summit/api-manager.yaml

oc -n 3scale-operator delete -f deploy/operator.yaml

for i in `ls -r deploy/crds/*3scale.hosted*_crd.yaml`; do oc delete -f $i ; done

oc delete -f deploy/cluster_role_binding.yaml
oc delete -f deploy/cluster_role.yaml
oc -n 3scale-operator delete -f deploy/service_account.yaml


for i in 3scale-operator 3scale-apimanager api-dev api-uat api-prod api-cicd
do
  oc delete project $i
done
