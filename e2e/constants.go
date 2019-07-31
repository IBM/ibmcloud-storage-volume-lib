/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package e2e

const (
	vpcZone                = "eu-gb-1"
	volumeSize             = 10
	iops                   = "0"
	volumeName             = "e2e-storage-volume"
	volumeType             = "vpc-block"
	generation             = "gt"
	vpcProfile             = "general-purpose"
	vpcConfigFilePath      = "/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e/config/vpc-config.toml"
	volumeEncryptionKeyCRN = "crn:v1:bluemix:public:kms:us-south:a/3198c72555b38a6bce0f48460003676d:b04bda72-e50c-481a-a689-35f6d2dd2cfd:key:28b772ef-5d00-40eb-b597-47973ccab82b"
)
