Storage file plugin deployment and testing steps:
=================================================

Steps to create storage plugin:
-------------------------------

1- Login to any carrier's worker node and get the latest storage plugin docker image from registry

	$docker pull 10.140.132.215:5001/armada/storage-file-plugin:latest

	or

	$docker pull alchemyregistry.hursley.ibm.com:5000/armada/storage-file-plugin:latest


2- Create the secret for storage plugin to order fileshare from SoftLayer. Note these secrets will be used by storage plugin to order fileshare from SL

	a- Create a toml file same as 'slclient.toml' but user need to provide correct info

	b- Use following command to create secret, please use the correct name so that same name should be used while creating storage plugin deployment
	
	$kubectl create secret generic <secret_name> --from-file=<path to toml file> -n kubx-masters
	Example:
		$kubectl create secret generic storage-secret --from-file=./slclient.toml -n kubx-masters

	User can verify secret by using following command
	
	$kubectl get secret -n kubx-masters


3- Deploy storage pod separately by using kubernetes deployment

	Following are the steps to deploy storage plugin on a perticular kubx cluster 

	Creating deployment for KubX master POD(for armada customer's cluster)

		a- Identify cluster id for which user want to create storage plugin, following command can be used for the same

			$kubectl get pod -n kubx-masters

			Get the cluster id from kubx master pod name, pattern for pod name is master-{cluster-id}-{some hex values}-{some hex values}
			Example:  master-cruiser1-446996936-qc8xd   this is one of the pod name and its cluster-id is 'cruiser1'

		b- Create a yaml file like 'storage_deployment_with_cert.yaml' or copy it and replace following ansible variable in this file
			{cluster-id} to actual cluster id which got from step-a

		c- Use following command to deploy storage plugin on carrier worker node 

			$kubectl create -f storage_deployment_with_cert.yaml

		d- Verify created storage pod by using following command

			$kubectl get pod -n kubx-masters | grep storage
			storage-deployment-3523242997-qfgw5                          1/1       Running            0          23m
			root@dev-mex01-carrier1-worker-01:~/Ara/storage-file#



Steps to test and verify storage plugin:
----------------------------------------

A- Please check the logs of the storage pod as follows

	$kubectl logs <name of the storage pod> -n kubx-masters

	It should show as follows ...
	
	"
		main.go:59] Provisioner kubernetes.io/ibmc-file specified
		controller.go:214] Starting nfs provisioner controller!
	"

B- For testing storage plugin please login to any of the kubx master's cluster worker node and follow steps 
   Sample files are available in the same githum location where this README is there, user can change values as per required. Please read the comments in the sample file as well
	
	1- Create the storage classes by using 'class.yaml' file, although storage classes(ibmc-file-bronze, ibmc-file-silver,ibmc-file-gold) would be pre-created
		
		$kubectl create -f class.yaml

		User can verify created class by using following command

		$kubectl get storageclass

	2- Create the pvc by using 'claim.yaml' file
		
		$kubectl create -f claim.yaml

		User can verify created claim by using following command

		$kubectl get pvc

		Expected output for successfully created pvc is as follows

		NAME            STATUS    VOLUME                                     CAPACITY   ACCESSMODES   AGE
		testclaimname   Bound     pvc-4ae1e5e8-e870-11e6-97eb-065bfd6dde1d   1Mi        RWX           1m


	3- Use created PVC in the kubernates pod and check the reading and writing to volume, sample 'claim_use.yaml' file can be used

		$kubectl create -f claim_use.yaml
