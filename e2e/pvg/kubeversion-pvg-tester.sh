#!/bin/bash
# Licensed Materials - Property of IBM
#
# (C) Copyright IBM Corp. 2017 All Rights Reserved
#
# US Government Users Restricted Rights - Use, duplicate or
# disclosure restricted by GSA ADP Schedule Contract with
# IBM Corp.
# encoding: utf-8

export PVG_TEST_UTILS=dove/dove-tools/py/pvg_run_test_utils.sh
. $PVG_TEST_UTILS

. ./pvg/pvg-tester-common-functions.sh

export PVG_BX_DASH_S="ibmc-file-e2e-test"

bx login -a $PVG_BX_DASH_A -u $BLUEMIX_USER -p $BLUEMIX_PASSWORD -c $PVG_BX_DASH_C -o $PVG_BX_DASH_O -s $PVG_BX_DASH_S
bx cs init --host $ARMADA_API_ENDPOINT
bx cs clusters
bx cs credentials-set --infrastructure-username $INFRA_USER --infrastructure-api-key $INFRA_KEY

export CLUSTER_LOCATION=$(get_first_zone)
export ZONE_PUBLIC_VLAN=$(get_public_vlan $CLUSTER_LOCATION)
export ZONE_PRIVATE_VLAN=$(get_private_vlan $CLUSTER_LOCATION)
export PVG_CLUSTER_LOCATION=$CLUSTER_LOCATION
export TEST_CLUSTER_NAME="stg_e2e_kube_${TEST_REGION}_${CLUSTER_LOCATION}"

# Delete the TEST_CLUSTER_NAME if it exists
cluster_id=$(bx cs clusters | awk "/$/"'{print $2}')
if [ "$cluster_id" != "" ]; then
    set +e
    rm_cluster $TEST_CLUSTER_NAME
    check_cluster_deleted $TEST_CLUSTER_NAME
    set -e
fi

export PVG_PHASE="${TEST_REGION}"
sed -i "s/PVG_PHASE/"$PVG_PHASE"/g" common/constants.go

set -x
# Check the cluster after creation and update
function set_cluster_config_and_run_tests {
    bx login -a $PVG_BX_DASH_A -u $BLUEMIX_USER -p $BLUEMIX_PASSWORD -c $PVG_BX_DASH_C -o $PVG_BX_DASH_O -s $PVG_BX_DASH_S
    bx cs init --host $ARMADA_API_ENDPOINT

    # Verify cluster is up and running
    echo "Checking the cluster for deployed state..."
    check_cluster_state $TEST_CLUSTER_NAME

    echo "Checking the worker nodes for ready state..."
    check_worker_state $TEST_CLUSTER_NAME

    bx cs cluster-config $TEST_CLUSTER_NAME | grep ^export | cut -d '=' -f 2
    configfile=$(bx cs cluster-config $TEST_CLUSTER_NAME | grep ^export | cut -d '=' -f 2)
    export KUBECONFIG=$configfile

    # Use Gate common function, it will setup the export KUBECONFIG safely
    setKubeConfig $TEST_CLUSTER_NAME

    # Use Gate common function
    addFullPathToCertsInKubeConfig
    export KUBECONFIG=$configfile
    export API_SERVER=$(kubectl config view | grep server | cut -f 2- -d ":" | tr -d " ")
    sleep 7m
    make KUBECONFIGPATH=$KUBECONFIG PVG_PHASE=$PVG_PHASE armada-storage-e2e-test
}

bx cs cluster-create --name $TEST_CLUSTER_NAME --zone $CLUSTER_LOCATION --public-vlan $ZONE_PUBLIC_VLAN  --private-vlan $ZONE_PRIVATE_VLAN --workers $PVG_BX_CLUSTER_WORKERS_COUNT --machine-type $PVG_BX_MACHINE_TYPE --kube-version $KUBE_VERSION

kube_versions=(1.12 1.13 1.14)
set_cluster_config_and_run_tests
for version in "${kube_versions[@]}"
do
    bx cs cluster-update --cluster $TEST_CLUSTER_NAME --kube-version $version -f
    check_cluster_update $TEST_CLUSTER_NAME $version
    set_cluster_config_and_run_tests
done

storage_plugin_pod=$(kubectl get pods -n kube-system| awk "/ibm-file-plugin-/"'{print $1}')
kubectl  logs $storage_plugin_pod -n kube-system > storage_pod_logs.txt

bx login -a $PVG_BX_DASH_A -u $BLUEMIX_USER -p $BLUEMIX_PASSWORD -c $PVG_BX_DASH_C -o $PVG_BX_DASH_O -s $PVG_BX_DASH_S
bx cs init --host $ARMADA_API_ENDPOINT

# Delete the TEST_CLUSTER_NAME if it exists
cluster_id=$(bx cs clusters | awk "/$/"'{print $2}')
if [ "$cluster_id" != "" ]; then
    rm_cluster $TEST_CLUSTER_NAME
    check_cluster_deleted $TEST_CLUSTER_NAME
fi

echo "Finished armada storage basic tests"
