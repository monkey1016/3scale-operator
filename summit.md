# Notes for Red Hat Summit 3Scale

## Launching a Cluster
1. Assign 6 CPUs and at least 18GB or RAM to CodeReady Containers
2. `crc setup`
3. `crc start`
4. `oc create secret generic redhat-registry --from-file=.dockerconfigjson=$HOME/.docker/config.json --type=kubernetes.io/dockerconfigjson`
5. `oc secrets link default redhat-registry --for=pull`
6. `for i in `ls deploy/crds/*_crd.yaml`; do oc create -f $i ; done`
7. `make local`
8. In another window: `oc create -f api-manager.yaml`
9. Log in to https://master.apps-crc.testing/
10. Use `system-seed` secret with `MASTER_USER` and `MASTER_PASSWORD` as credentials
