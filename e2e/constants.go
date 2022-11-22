package e2e

// TODO copied form vpc .. Needs to be cleaned up
const (
	vpcZone                = "us-east-1"
	resourceGroupID        = "f2075e07c1a362e26bdfc856771798a7"
	volumeSize             = 10
	iops                   = "0"
	volumeName             = "e2e-storage-volume"
	volumeType             = "vpc-block"
	generation             = "gc"
	vpcProfile             = "general-purpose"
	vpcConfigFilePath      = "/src/github.com/IBM/ibmcloud-storage-volume-lib/e2e/config/vpc-config.toml"
	volumeEncryptionKeyCRN = "ENCRYPTION_KEY_CRN"
)
