package e2e

const (
	vpcZone                = "VPC_ZONE"
	resourceGroupID        = "RESOURCEGROUP"
	volumeSize             = 10
	iops                   = "0"
	volumeName             = "e2e-storage-volume"
	volumeType             = "vpc-block"
	generation             = "gt"
	vpcProfile             = "general-purpose"
	vpcConfigFilePath      = "/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e/config/vpc-config.toml"
	volumeEncryptionKeyCRN = "ENCRYPTION_KEY_CRN"
)
