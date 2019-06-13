#!/bin/bash
set -eou pipefail

if [ "$#" -ne 2 ]; then
   echo "$0 <release_name> <namespace>"
   exit 1
fi

RELEASE_NAME=$1
NAMESPACE=$2

echo "Migrating Helm v2 '${RELEASE_NAME}' state to Helm v3 state in namespace '${NAMESPACE}' ........."

# Get the v2 state for ${RELEASE_NAME}
kubectl get configmap ${RELEASE_NAME} -n kube-system -o json > ${RELEASE_NAME}-cm.json
cp ${RELEASE_NAME}-cm.json ${RELEASE_NAME}-secret.json

# Update fields and values to correspond to v3 state secret object
sed -i -e 's/ConfigMap/Secret/g' ./${RELEASE_NAME}-secret.json
sed -i -e 's/MODIFIED_AT/modifiedAt/g' ./${RELEASE_NAME}-secret.json
sed -i -e 's/NAME/name/g' ./${RELEASE_NAME}-secret.json
sed -i -e 's/OWNER/owner/g' ./${RELEASE_NAME}-secret.json
sed -i -e 's/STATUS/status/g' ./${RELEASE_NAME}-secret.json
sed -i -e 's/VERSION/version/g' ./${RELEASE_NAME}-secret.json
sed -i -e 's/configmaps/secrets/g' ./${RELEASE_NAME}-secret.json
sed -i -e "s/kube-system/${NAMESPACE}/g" ./${RELEASE_NAME}-secret.json
sed -i -e 's/TILLER/helm/g' ./${RELEASE_NAME}-secret.json
STATUS=`jq '.metadata.labels.status' ${RELEASE_NAME}-secret.json | tr '[:upper:]' '[:lower:]'`
jq ".metadata.labels.status=${STATUS}" ${RELEASE_NAME}-secret.json > ${RELEASE_NAME}-secret.tmp && mv ${RELEASE_NAME}-secret.tmp ${RELEASE_NAME}-secret.json

# Add "helm.sh/release" type
sed -i -e 'x; ${s/.*/    },/;p;x}; 1d' ./${RELEASE_NAME}-secret.json
sed -i -e '$ i\    "type": "helm.sh/release"' ./${RELEASE_NAME}-secret.json

# Deploy the ${RELEASE_NAME} secret into the ${NAMESPACE} namespace
kubens ${NAMESPACE}
kubectl create -f ${RELEASE_NAME}-secret.json

echo "'${RELEASE_NAME}' v3 state to v3 state in namespace '${NAMESPACE}' migrated."
