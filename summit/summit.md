# Notes for Red Hat Summit 3Scale

## Building Operator Image

> This step is **Optional**  

This requires an account on quay.io
1. `docker login quay.io` and enter credentials
2. `make build IMAGE=quay.io/<your account>/3scale-bof-summit-2020 VERSION=summit-2.7`
3. `make push IMAGE=quay.io/<your account>/3scale-bof-summit-2020 VERSION=summit-2.7`
   (Careful, could take 15 or so minutes for it to become available)

## Provision a Cluster
For provisioning a cluster, you have several options:

* [OpenShift Installer](https://github.com/openshift/installer) (has cost associated with running in a cloud provider like AWS)
* RHPDS [Insert instructions for provisioning a cluster] (could have a cost associated with it)
* [CodeReady Containers](#crc_install) with sufficient resources

## Setup Script <a name="setup_script></a>
Once you have a cluster up and running, with cluster admin rights, update the `summit/api-manager.yaml` file
by providing the right value for `wildcardDomain`. For example, if your cluster name is located at
`cluster-nyc-3c57.nyc-3c57.example.opentlc.com`, the `wildcardDomain` could be
`hosted.apps.cluster-nyc-3c57.nyc-3c57.example.opentlc.com`. Once you've done that, you can run the following
script to get everything configured:

```bash
./summit-setup.sh
```

It may take up to 30 minutes to provision all the pods required to run 3scale.

> After all the pods are up, you'll need to create routes for the APICast deployments to be able to access them
> from outside the cluster.

### Setup Script Actions
The setup script creates a number of items in order:

1. Provisions the projects and necessary secrets for the demo
2. Sets up secrets to pull images from the Red Hat registry (requires you log in using `docker login` prior
to running this command)
3. Creates the necessary Custom Resource Definitions
4. Creates the service accounts, roles, and sets up the role bindings
5. Deploys the Operator
6. Provisions the API Manager and APICasts using the `APIManager` Custom Resource Definition

## Step by Step Install <a name="crc_install"></a>
1. Assign 6 CPUs and at least 18GB or RAM to CodeReady Containers
2. `crc setup`
3. `crc config set cpus 6`
4. `crc config set memory 18432`
3. `crc start`

Once the cluster is up, proceed to run the [setup script](#setup_script)

## Creating an API
See API Demo for a full example.

You can create a simple API using the provided template:

```bash
oc process -f summit/demo-app/api-template.yaml -p=SERVICE_NAME=echo-app \
-p=SYSTEM_MASTER_NAMESPACE=3scale-apimanager \
-p=SECRET_TOKEN=SomethingToReplaceLater123 \
-p=PRIVATE_BASE_URL=https://echo-api.3scale.net:443 \
-p=PUBLIC_BASE_URL=https://<route to apicast>:443 \
-p=TENANT_TOKEN_SECRET=tenant-dev \
-p=SECRET_NAMESPACE=api-dev | oc apply -f -
```
Make sure to replace `<route to apicast>` with the right address to your apicast based on the cluster you
provisioned.

## Cleaning Up
Once you're done with everything, you can clean up all the resources associated with this demo but executing
the delete script:

```bash
./summit/summit-delete.sh
```

## Notes
API Management as code
artefact file <- a file that contains a policy

1. Install 3scale somehow gateways in separate namespaces
2. API defined in 3scale
3. Build environment ready
4. Version 1 already exists
5. Git repo commit triggers pipeline
6. Build and deploy to Dev gateway
7. Update policy
8. Run test (BDD maybe)
9. Deploy to staging gateway

Each tenant has an endpoint the gateway can use to get data
https://github.com/3scale/APIcast/blob/master/doc/parameters.md#threescale_portal_endpoint