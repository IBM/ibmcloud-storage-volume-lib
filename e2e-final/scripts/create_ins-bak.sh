#!/bin/sh
echo "Generate ssh key";
ssh-keygen -t rsa -N "" -f ssh-key

export ssh_key=$(cat ssh-key.pub)
export vpc_api_endpoint="https://us-south-genesis-dal-dev45-etcd.iaasdev.cloud.ibm.com"
export profile_name="bx2-4x16"
export image_id="r134-9f6b534b-6061-40f4-ac42-5aba4dd0da7f"
E2ETEST_IC_API_EP="test.cloud.ibm.com"
IC_REGION="us-south"
api_version="2022-01-01"

ibmcloud login -a $E2ETEST_IC_API_EP --apikey  $IC_API_KEY -r $IC_REGION
if [[ $rc -ne 0 ]]; then echo "Error: IBM Cloud Login failed!!!"; exit 1; fi

iam_token=$(ibmcloud iam oauth-tokens | awk -F "token: " '{ print $2}')
echo $iam_token

echo "Adding ssh keys"
key=$(curl -s -X POST "$vpc_api_endpoint/v1/keys?version=$api_version&generation=2"   -H "Content-Type: application/json" -H "Authorization:$iam_token"  -d "{\"name\": \"test-ssh-key\", \"public_key\": \"$ssh_key\"}")
if [[ $key == *"errors"* ]]; then  >&2 echo "Adding ssh key failed!!! Error: $key"; exit 1; fi
key_id=$(echo $key | jq -r '.id')

if [[ $key_id == "" ]]; then >&2 echo "Error: Adding ssh key failed!!!"; exit 1; fi
echo "created key $key_id"


echo "Creating an instance"
instance=$(curl -s -X POST "$vpc_api_endpoint/v1/instances?version=$api_version&generation=2"   -H "Authorization:$iam_token"   -d '{
        "name": "my-instance",
        "zone": {
          "name": "us-south-3"
        },
        "vpc": {
          "id": "'$vpc'"
        },
        "primary_network_interface": {
          "subnet": {
            "id": "'$subnet'"
          }
        },
        "keys":[{"id": "'$key_id'"}],
        "profile": {
          "name": "'$profile_name'"
         },
        "image": {
          "id": "'$image_id'"
         }
        }')
if [[ $instance == *"errors"* ]]; then  >&2 echo "Creating instance failed!!! Error: $instance"; exit 1; fi
instance_id=$(echo $instance | jq -r '.id')

if [[ $instance_id == "" ]]; then >&2 echo "Error: Creating instance failed!!!"; exit 1; fi
echo "created instance $instance_id"
echo "Get network interface of the instance"
network_interface=$(curl -s -X GET "$vpc_api_endpoint/v1/instances/$instance_id?version=$api_version&generation=2"   -H "Authorization: $iam_token" | jq -r '.primary_network_interface.id')

if [[ $network_interface == "" ]]; then >&2 echo "Error: Getting network interface failed!!!"; exit 1; fi
echo "created network_interface $network_interface"
echo "Create a floating ip for the instance"
ip=$(curl -s -X POST "$vpc_api_endpoint/v1/floating_ips?version=$api_version&generation=2"   -H "Authorization:$iam_token"   -d '{
        "name": "my-floatingip",
        "target": {
            "id":"'$network_interface'"
        }
      }
')
if [[ $ip == *"errors"* ]]; then >&2 echo "Creating ip failed!!! Error: $ip"; exit 1; fi
ip_address=$(echo $ip | jq -r '.address')
floating_ip_id=$(echo $ip | jq -r '.id')
if [[ $ip_address == "" ]]; then >&2 echo "Error: Getting floating ip failed!!!"; exit 1; fi
echo "created ip $ip_address"
echo "output : Instance_id=$instance_id, Instance_ip=$ip_address, key_id=$key_id, floating_ip_id=$floating_ip_id"




