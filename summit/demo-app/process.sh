#!/usr/bin/env bash

oc process -f api-template.yaml \
-p=SERVICE_NAME="api-versioning-service-dev" \
-p=SYSTEM_MASTER_NAMESPACE="3scale-apimanager" \
-p=SECRET_TOKEN=562870a6fd3d192692c2400a023b8c5f \
-p=PRIVATE_BASE_URL=http://api-versioning-service-web.api-dev.svc.cluster.local:8162 \
-p=PUBLIC_BASE_URL=https://dev-api-dev.apps.cluster-nyc-6b8f.nyc-6b8f.example.opentlc.com:443 \
-p=TENANT_TOKEN_SECRET=tenant-dev \
-p=TENANT_TOKEN_NAMESPACE=api-dev | oc apply -f -