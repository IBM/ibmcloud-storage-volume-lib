#!/bin/bash
# Licensed Materials - Property of IBM
#
# (C) Copyright IBM Corp. 2017 All Rights Reserved
#
# US Government Users Restricted Rights - Use, duplicate or
# disclosure restricted by GSA ADP Schedule Contract with
# IBM Corp.
# encoding: utf-8

# Before running this script, export the needed variables
#export PVG_BX_USER=<Bluemix User ID>
#export PVG_BX_PWD=<Bluemix Password>
#export PVG_SL_USERNAME=<SL account user name>
#export PVG_SL_API_KEY=<SL API Key>
#export PVG_BX_DASH_O=<Bluemix Org name>
#export PVG_BX_DASH_S=<Bluemix Space name>
#export PVG_BX_DASH_A=<Bluemix API for login>
#export ARMADA_API_ENDPOINT=<API Endpoint>
#export PVG_CRUISER_PRIVATE_VLAN=<Private VLAN>
#export PVG_CRUISER_PUBLIC_VLAN=<Public VLAN>
#export FREE_DATACENTER=<Datacenter name: mex01>

set -x
set -e
# Skip the tests during month end
TODAY=`/bin/date +%d`
TOMORROW=`/bin/date +%d -d "1 day"`
if [ $TOMORROW -lt $TODAY ] || [ $TODAY -eq 1 ]; then
  exit 0
fi

# Source Gate supplied common functions
. $PVG_TEST_UTILS
export PVG_CLUSTER_CRUISER="stg-e2e-cluster"

# Remove the created clusters
function rm_cluster {
    removed=1
    cluster_name=$1

    bx cs clusters
    set +e
    for i in {1..3}; do
        if bx cs cluster-rm $cluster_name -f; then
            removed=0
            break
        fi
        sleep 30
    done

    return $removed
}

# Delete the PVG_CLUSTER_CRUISER if it exists
set +e
rm_cluster $PVG_CLUSTER_CRUISER
check_cluster_deleted $PVG_CLUSTER_CRUISER
set -e

# Create a cruiser
cruiser_create $PVG_CLUSTER_CRUISER u2c.2x4 1

# Put a small delay to let things settle
sleep 30

bx cs clusters

unset DOCKER_HOST
docker -v
mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
mkdir -p $GOPATH/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e
DIR="$(pwd)"
echo "Present working directory: $DIR"
ls -altr $DIR
rm -rf .git
sed -i "s/PVG_PHASE/"$PVG_PHASE"/g" common/constants.go
rsync -az . $GOPATH/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e
cd $GOPATH/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e
ls -altr $GOPATH/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e

# Verify cluster is up and running
echo "Checking the cluster for deployed state..."
check_cluster_state $PVG_CLUSTER_CRUISER

echo "Checking the worker nodes for deployed state..."
check_worker_state $PVG_CLUSTER_CRUISER

# Run sniff tests against cluster
bx cs clusters
bx cs cluster-get $PVG_CLUSTER_CRUISER
bx cs workers $PVG_CLUSTER_CRUISER

# Use Gate common function, it will setup the export KUBECONFIG safely
setKubeConfig $PVG_CLUSTER_CRUISER
cat $KUBECONFIG

# Use Gate common function
addFullPathToCertsInKubeConfig
cat $KUBECONFIG

echo "Pods running on cluster: $PVG_CLUSTER_CRUISER"
kubectl get pods -n kube-system

echo "Checking storage plugin and watcher pods status on cluster: $PVG_CLUSTER_CRUISER"
function check_pods_state {
  attempts=0
  while true; do
    attempts=$((attempts+1))
    file_plugin_pod_status=$(kubectl get pods -n kube-system| awk "/ibm-file-plugin-/"'{print $2}')
    watcher_pod_status=$(kubectl get pods -n kube-system| awk "/ibm-storage-watcher-/"'{print $2}')

    if [ "$watcher_pod_status" = "1/1" -a "$file_plugin_pod_status" = "1/1" ]; then
      echo "Armada storage plugin and watcher pods were running well."
      break
    fi

    if [[ $attempts -gt 30 ]]; then
      echo "Armada storage plugin and watcher pods were not running well."
      kubectl get pods -n kube-system| awk "/ibm-storage-watcher-/"
      exit 1
    fi

    echo "ibm-file-plugin state == $file_plugin_pod_status; ibm-storage-watcher state == $watcher_pod_status.  Sleeping 10 seconds"
    sleep 10

  done
}

check_pods_state

kubectl cluster-info
kubectl config view

cd $GOPATH/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e
DIR="$(pwd)"
echo "Present working directory: $DIR"
ls -altr $DIR

echo "Starting armada storage basic e2e tests"
export API_SERVER=$(kubectl config view | grep server | cut -f 2- -d ":" | tr -d " ")
make KUBECONFIGPATH=$KUBECONFIG PVG_PHASE=$PVG_PHASE storage-e2e-test
echo "Finished armada storage basic e2e tests"

exit 0
