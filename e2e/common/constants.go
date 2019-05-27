package common

// busyboxTestImage is the list of images used in common test. These images should be prepulled
// before a tests starts, so that the tests won't fail due image pulling flakes.
const (
	BusyboxTestImage     = "gcr.io/google_containers/busybox:1.24"
	PluginName           = "ibm.io/ibmc-file"
	MountPath            = "/mnt/test"
	NamespaceName        = "volume-lib-e2e-namespace"
	ProvisionerPodName   = "storage-deployment"
	DeploymentName       = "storage-deployment"
	ClaimPrefix          = "PVG_PHASE-"
	PluginImage          = "armada-master/storage-file-plugin:latest"
	VolumeMountPath      = "/export"
	DeploymentGOFilePATH = "src/github.com/IBM/ibmcloud-storage-volume-lib/e2e/deploy/kube-config/deployment.yaml"
)
