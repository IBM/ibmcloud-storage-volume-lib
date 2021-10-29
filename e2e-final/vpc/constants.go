/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package vpc

const (
	volumeSize                = 10
	iops                      = "0"
	volumeName                = "e2e-storage-volume"
	volumeType                = "vpc-block"
	generation                = "gt"
	vpcProfile                = "general-purpose"
	vpcConfigFilePath         = "/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e-final/config/vpc-config.toml"
	testCasesForBlock         = "/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e-final/config/test-cases-block.yml"
	testCasesForFile          = "/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e-final/config/test-cases-file.yml"
	testCase5iops             = "5iops"
	testCase10iops            = "10iops"
	testCaseCustom            = "custom"
	testCaseWithEncryption    = "with-encryption"
	testCaseWithoutEncryption = "without-encryption"
	testCaseSizes             = "sizes"
	fileAPIVersion            = "2021-10-20"
)
