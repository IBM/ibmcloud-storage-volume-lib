/*
# Licensed Materials - Property of IBM
#
# (C) Copyright IBM Corp. 2017 All Rights Reserved
#
# US Government Users Restricted Rights - Use, duplicate or
# disclosure restricted by GSA ADP Schedule Contract with
# IBM Corp.
# encoding: utf-8
*/

package tests

import (
	commontest "github.com/IBM/ibmcloud-storage-volume-lib/e2e/common"
	"github.com/IBM/ibmcloud-storage-volume-lib/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	clientset "k8s.io/client-go/kubernetes"
)

var _ = framework.KubeDescribe("Volumes [Feature:VolumeMountTest]", func() {
	f := framework.NewDefaultFramework("storage")
	// filled in BeforeEach
	var c clientset.Interface
	var ns string

	BeforeEach(func() {
		c = f.ClientSet
		ns = f.Namespace.Name
	})

	framework.KubeDescribe("Test mount point: Create mount point, read, write...", func() {
		PIt("PRESTAGE: Test mount point: Create mount point, read, write... [Serial]", func() {
			By("Creating a claim with a dynamic provisioning annotation")
			claim := commontest.NewClaim(ns, f.PluginNameShort+"-bronze", "20Gi")
			defer func() {
				c.Core().PersistentVolumeClaims(ns).Delete(claim.Name, nil)
			}()
			claim, err := c.Core().PersistentVolumeClaims(ns).Create(claim)
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
