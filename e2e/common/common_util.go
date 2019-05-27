/*
Copyright 2017 The IBM Storage Plugin Author.

*/

package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/IBM/ibmcloud-storage-volume-lib/e2e/framework"
	"github.com/opencontainers/runc/libcontainer/selinux"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/apis/storage/v1beta1"
	"k8s.io/client-go/pkg/runtime"
	utilyaml "k8s.io/client-go/pkg/util/yaml"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// StorageClassAnnotation represents the storage class associated with a resource.
const StorageClassAnnotation = "volume.beta.kubernetes.io/storage-class"

var (
	pluginName  = PluginName
	claimPrefix = "local-"
	billingType = "hourly"
)
var accessMode v1.PersistentVolumeAccessMode

func init() {
	// Initialise default access modes to be used PVC creation
	switch accessModeEnv := os.Getenv("PVC_ACCESS_MODE"); accessModeEnv {
	case "RWO":
		accessMode = v1.ReadWriteOnce
	case "RWM":
		accessMode = v1.ReadWriteMany
	case "ROM":
		accessMode = v1.ReadOnlyMany
	default:
		accessMode = v1.ReadWriteMany
	}

	pluginNameEnv := os.Getenv("PLUGIN_NAME")
	if len(pluginNameEnv) > 0 {
		pluginName = pluginNameEnv
	}
	claimPrefixEnv := os.Getenv("PVG_PHASE")
	if len(claimPrefixEnv) > 0 {
		claimPrefix = claimPrefixEnv + "-"
	}
}

func TestDynamicProvisioning(client clientset.Interface, claim *v1.PersistentVolumeClaim) {
	pv := TestCreate(client, claim)
	TestWrite(client, claim)
	TestRead(client, claim)
	TestDelete(client, claim, pv)
}

func TestCreate(client clientset.Interface, claimExpected *v1.PersistentVolumeClaim) *v1.PersistentVolume {
	err := framework.WaitForPersistentVolumeClaimPhase(v1.ClaimBound, client, claimExpected.Namespace, claimExpected.Name, framework.Poll, framework.ClaimProvisionTimeout)
	Expect(err).NotTo(HaveOccurred())

	By("Checking the claim: " + claimExpected.Name)
	// Get new copy of the claim
	claim, err := client.Core().PersistentVolumeClaims(claimExpected.Namespace).Get(claimExpected.Name)
	Expect(err).NotTo(HaveOccurred())

	// Get the bound PV
	pv, err := client.Core().PersistentVolumes().Get(claim.Spec.VolumeName)
	Expect(err).NotTo(HaveOccurred())

	// Check sizes
	pvCapacity := pv.Spec.Capacity[v1.ResourceName(v1.ResourceStorage)]
	claimCapacity := claim.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)]
	Expect(claimCapacity.Value()).To(Equal(pvCapacity.Value()))

	// Check PV properties
	Expect(pv.Spec.PersistentVolumeReclaimPolicy).To(Equal(v1.PersistentVolumeReclaimDelete))
	expectedAccessModes := claimExpected.Spec.AccessModes
	Expect(pv.Spec.AccessModes).To(Equal(expectedAccessModes))
	Expect(pv.Spec.ClaimRef.Name).To(Equal(claim.ObjectMeta.Name))
	Expect(pv.Spec.ClaimRef.Namespace).To(Equal(claim.ObjectMeta.Namespace))

	return pv
}

// We start two pods, first in testWrite and second in testRead:
// - The first writes 'hello word' to the MountPath (= the volume).
// - The second one runs grep 'hello world' on MountPath.
// If both succeed, Kubernetes actually allocated something that is
// persistent across pods.
func TestWrite(client clientset.Interface, claim *v1.PersistentVolumeClaim) {
	By("Checking the created volume is writable")
	RunInPodWithVolume(client, claim.Namespace, claim.Name, "echo 'hello world' > "+MountPath+"/data")

	// Unlike cloud providers, kubelet should unmount NFS quickly
	By("Sleeping to let kubelet destroy pods")
	time.Sleep(5 * time.Second)
}

func TestRead(client clientset.Interface, claim *v1.PersistentVolumeClaim) {
	By("Checking the created volume is readable and retains data")
	RunInPodWithVolume(client, claim.Namespace, claim.Name, "ls "+MountPath+"/data")

	// Unlike cloud providers, kubelet should unmount NFS quickly
	By("Sleeping to let kubelet destroy pods")
	time.Sleep(5 * time.Second)
}

func TestDelete(client clientset.Interface, claim *v1.PersistentVolumeClaim, pv *v1.PersistentVolume) {
	By("Deleting the claim: " + claim.Name)
	framework.ExpectNoError(client.Core().PersistentVolumeClaims(claim.Namespace).Delete(claim.Name, nil))

	// Wait for the PV to get deleted too.
	//framework.ExpectNoError(framework.WaitForPersistentVolumeDeleted(client, pv.Name, 5*time.Second, 5*time.Minute))
}

func NewClaim(ns string, storageClassName string, requestedStorageSize string) *v1.PersistentVolumeClaim {
	labels := make(map[string]string)
	labels[billingType] = "hourly"

	claim := v1.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: claimPrefix,
			Namespace:    ns,
			Labels:       labels,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				accessMode,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse(requestedStorageSize),
				},
			},
		},
	}

	claim.Annotations = map[string]string{
		StorageClassAnnotation: storageClassName,
	}

	return &claim
}

// RunInPodWithVolume runs a command in a pod with given claim mounted to /mnt directory.
func RunInPodWithVolume(c clientset.Interface, ns, claimName, command string) {
	pod := &v1.Pod{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			GenerateName: "pvc-volume-tester-",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "volume-tester",
					Image:   BusyboxTestImage,
					Command: []string{"/bin/sh"},
					Args:    []string{"-c", command},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "my-volume",
							MountPath: MountPath,
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
			Volumes: []v1.Volume{
				{
					Name: "my-volume",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: claimName,
							ReadOnly:  false,
						},
					},
				},
			},
		},
	}
	pod, err := c.Core().Pods(ns).Create(pod)
	defer func() {
		framework.ExpectNoError(c.Core().Pods(ns).Delete(pod.Name, nil))
	}()
	framework.ExpectNoError(err, "Failed to create pod: %v", err)
	framework.ExpectNoError(framework.WaitForPodSuccessInNamespace(c, pod.Name, pod.Namespace))
}

func NewStorageClass(storageClassName string, storageClassType string, iopsPerGB string, sizeRange string, mountOptions string, billingType string) *v1beta1.StorageClass {
	return &v1beta1.StorageClass{
		TypeMeta: unversioned.TypeMeta{
			Kind: "StorageClass",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: storageClassName,
			Labels: map[string]string{
				"kubernetes.io/cluster-service": "true",
			},
		},
		Provisioner: pluginName,
		Parameters: map[string]string{
			"type":         storageClassType,
			"iopsPerGB":    iopsPerGB,
			"sizeRange":    sizeRange,
			"mountOptions": mountOptions,
			"billingType":  billingType,
		},
	}
}

func StartProvisionerPod(c clientset.Interface, ns string) *v1.Pod {
	podClient := c.Core().Pods(ns)
	By("startProvisionerPod: Init")
	provisionerPod := &v1.Pod{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name: ProvisionerPodName,
			Labels: map[string]string{
				"role": ProvisionerPodName,
			},
		},

		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  ProvisionerPodName,
					Image: PluginImage,
					SecurityContext: &v1.SecurityContext{
						Capabilities: &v1.Capabilities{
							Add: []v1.Capability{"DAC_READ_SEARCH"},
						},
					},
					Args: []string{
						fmt.Sprintf("-provisioner=%s", pluginName),
						"-grace-period=90",
					},
					Ports: []v1.ContainerPort{
						{Name: "nfs", ContainerPort: 2049},
						{Name: "mountd", ContainerPort: 20048},
						{Name: "rpcbind", ContainerPort: 111},
						{Name: "rpcbind-udp", ContainerPort: 111, Protocol: v1.ProtocolUDP},
					},
					Env: []v1.EnvVar{
						{
							Name: "POD_IP",
							ValueFrom: &v1.EnvVarSource{
								FieldRef: &v1.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						},
					},
					ImagePullPolicy: v1.PullIfNotPresent,
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "export-volume",
							MountPath: VolumeMountPath,
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "export-volume",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}
	By("startProvisionerPod:Defined Spec")
	provisionerPod, err := podClient.Create(provisionerPod)
	framework.ExpectNoError(err, "Failed to create %s pod: %v", provisionerPod.Name, err)

	framework.ExpectNoError(framework.WaitForPodRunningInNamespace(c, provisionerPod))

	By("Locating the provisioner pod")
	pod, err := podClient.Get(provisionerPod.Name)
	framework.ExpectNoError(err, "Cannot locate the provisioner pod %v: %v", provisionerPod.Name, err)

	By("Sleeping a bit to give the provisioner time to start")
	time.Sleep(5 * time.Second)
	return pod
}

func StartProvisionerDeployment(c clientset.Interface, ns string) (*v1.Service, *extensions.Deployment) {
	gopath := os.Getenv("GOPATH")
	service := SvcFromManifest(path.Join(gopath, DeploymentGOFilePATH))

	deployment := DeployFromManifest(path.Join(gopath, DeploymentGOFilePATH))

	tmpDir, err := ioutil.TempDir("", ProvisionerPodName+"-deployment")
	Expect(err).NotTo(HaveOccurred())
	if selinux.SelinuxEnabled() {
		fcon, err := selinux.Getfilecon(tmpDir)
		Expect(err).NotTo(HaveOccurred())
		context := selinux.NewContext(fcon)
		context["type"] = "svirt_sandbox_file_t"
		err = selinux.Chcon(tmpDir, context.Get(), false)
		Expect(err).NotTo(HaveOccurred())
	}
	deployment.Spec.Template.Spec.Volumes[0].HostPath.Path = tmpDir
	deployment.Spec.Template.Spec.Containers[0].Image = PluginImage
	deployment.Spec.Template.Spec.Containers[0].Args = []string{
		fmt.Sprintf("-provisioner=%s", pluginName),
		"-grace-period=90",
	}

	service, err = c.Core().Services(ns).Create(service)
	framework.ExpectNoError(err, "Failed to create %s service: %v", service.Name, err)

	deployment, err = c.Extensions().Deployments(ns).Create(deployment)
	framework.ExpectNoError(err, "Failed to create %s deployment: %v", deployment.Name, err)

	framework.ExpectNoError(framework.WaitForDeploymentPodsRunning(c, ns, deployment.Name))

	By("Sleeping a bit to give the provisioner time to start")
	time.Sleep(5 * time.Second)

	return service, deployment
}

// SvcFromManifest reads a .json/yaml file and returns the json of the desired
func SvcFromManifest(fileName string) *v1.Service {
	var service v1.Service
	data, err := ioutil.ReadFile(fileName)
	Expect(err).NotTo(HaveOccurred())

	r := ioutil.NopCloser(bytes.NewReader(data))
	decoder := utilyaml.NewDocumentDecoder(r)
	var chunk []byte
	for {
		chunk = make([]byte, len(data))
		_, err := decoder.Read(chunk)
		chunk = bytes.Trim(chunk, "\x00")
		Expect(err).NotTo(HaveOccurred())
		if strings.Contains(string(chunk), "kind: Service") {
			break
		}
	}

	json, err := utilyaml.ToJSON(chunk)
	Expect(err).NotTo(HaveOccurred())
	Expect(runtime.DecodeInto(api.Codecs.UniversalDecoder(), json, &service)).NotTo(HaveOccurred())

	return &service
}

// deployFromManifest reads a .json/yaml file and returns the json of the desired
func DeployFromManifest(fileName string) *extensions.Deployment {
	var deployment extensions.Deployment
	data, err := ioutil.ReadFile(fileName)
	Expect(err).NotTo(HaveOccurred())

	r := ioutil.NopCloser(bytes.NewReader(data))
	decoder := utilyaml.NewDocumentDecoder(r)
	var chunk []byte
	for {
		chunk = make([]byte, len(data))
		_, err := decoder.Read(chunk)
		chunk = bytes.Trim(chunk, "\x00")
		Expect(err).NotTo(HaveOccurred())
		if strings.Contains(string(chunk), "kind: Deployment") {
			break
		}
	}

	json, err := utilyaml.ToJSON(chunk)
	Expect(err).NotTo(HaveOccurred())
	Expect(runtime.DecodeInto(api.Codecs.UniversalDecoder(), json, &deployment)).NotTo(HaveOccurred())

	return &deployment
}

func ScaleDeployment(c clientset.Interface, ns, name string, newSize int32) {
	deployment, err := c.Extensions().Deployments(ns).Get(name)
	Expect(err).NotTo(HaveOccurred())
	deployment.Spec.Replicas = &newSize
	updatedDeployment, err := c.Extensions().Deployments(ns).Update(deployment)
	Expect(err).NotTo(HaveOccurred())
	framework.ExpectNoError(framework.WaitForDeploymentPodsRunning(c, ns, updatedDeployment.Name))
	// Above is not enough. Just sleep to prevent conflict when doing Update.
	// kubectl Scaler would be ideal. or WaitForDeploymentStatus
	time.Sleep(5 * time.Second)
}
