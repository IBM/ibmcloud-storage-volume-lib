#!/bin/sh
instance=$1
key=$2
floating_ip=$3
export vpc_api_endpoint="https://us-south-genesis-dal-dev45-etcd.iaasdev.cloud.ibm.com"
E2ETEST_IC_API_EP="test.cloud.ibm.com"
IC_REGION="us-south"
api_version="2022-01-01"

ibmcloud login -a $E2ETEST_IC_API_EP --apikey  $IC_API_KEY -r $IC_REGION
if [[ $rc -ne 0 ]]; then echo "Error: IBM Cloud Login failed!!!"; exit 1; fi

iam_token=$(ibmcloud iam oauth-tokens | awk -F "token: " '{ print $2}')
echo $iam_token



echo "Delete an instance"
curl -s -X DELETE "$vpc_api_endpoint/v1/instances/$instance?version=$api_version&generation=2"   -H "Authorization:$iam_token"

sleep 10
echo "Delete ssh keys"
curl -s -X DELETE "$vpc_api_endpoint/v1/keys/$key?version=$api_version&generation=2"   -H "Content-Type: application/json" -H "Authorization:$iam_token"
sleep 10
echo "Delete floating ip"
curl -s -X DELETE "$vpc_api_endpoint/v1/floating_ips/$floating_ip?version=$api_version&generation=2"   -H "Authorization:$iam_token"


