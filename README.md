# ibmcloud-storage-volume-lib
This library is for IBM Cloud Kubernetes Services persistent volume and snapshot management(create, delete, modify etc).
As of now this only have block snapshot and volume creation functionalities but going forward it will have more functionalities.

# Purpose:
By using this library user can enable multiple volume providers as well by just adding/modifying configuration, as of now this library has only Softlayer provider(block only) support.

# How to use:
To use this library user need to see what all capabilities are there in this library by just checking `ibmcloud-storage-volume-lib/lib/provider/*_manager.go` files and concreate implementation for each provider which can be seen under `ibmcloud-storage-volume-lib/volume-providers` directory.

NOTE: as of now this library has only Softlayer block support

## Steps to use:
Following steps has to followed by user to use this library
### Step 1:
Update the glide.yaml to get the source code of this library, also only export the following packages to their application

```
"github.ibm.com/IBM/ibm-volume-lib/config"
"github.ibm.com/IBM/ibm-volume-lib/lib/provider"
"github.ibm.com/IBM/ibm-volume-lib/provider/local"
provider_util "github.ibm.com/IBM/ibm-volume-lib/provider/utils"
```

User can see the refrence from `ibmcloud-storage-volume-lib/main.go` sample file

### Step 2:
How to provide configuration
##### Option 1:
User need to modify `ibmcloud-storage-volume-lib/etc/libconfig.toml` for enabling supported providers and provide IAaS a/c details to connect to the IAaS provider.
As of now this librray has only Softlayer implementation so user just need to provide value of `softlayer_block_enabled`,  `softlayer_block_provider_name`, `softlayer_username` , `softlayer_api_key` , `softlayer_datacenter` to enable Softlayer block although other provider can also be enabled by updating configuration file

#### Option 2:
User can also create a secret from  `ibmcloud-storage-volume-lib/etc/libconfig.toml` file and mount the secret into the pod which uses this library and mount point

```
$kubectl create secret generic volume-lib-secret --from-file=./libconfig.toml

and use as follows in the deployment file

apiVersion: extensions/v1beta1
kind: Deployment
metadata:
...
...
spec:
...
...
      containers:
			....
			....
          - name: SECRET_CONFIG_PATH
            value: /etc/volume_secret_path
        volumeMounts:
          - name: volume-lib-secret
            mountPath: /etc/volume_secret_path
      volumes:
      - name: volume-lib-secret
        secret:
          secretName: volume-lib-secret
```

### Step 3:
From the implementaion file which uses use this library, user need to initilaze the providers and open the sessions to backend IAaS provider and to do that user just need to call the following method from there

`github.ibm.com/IBM/ibm-volume-lib/provider/utils` utility packages, reference code can be found `ibmcloud-storage-volume-lib/main.go`

```
providerRegistry, err := provider_util.InitProviders(conf, logger)

and

sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, conf.Softlayer.SoftlayerBlockProviderName, logger)
```

`OpenProviderSession` call will be required for all enabled provider.
Once session is ready the user can all any method/interface which are part of `ibmcloud-storage-volume-lib/lib/provider/*_manager.go`. Sample code can be found `ibmcloud-storage-volume-lib/main.go`

## How to build the sample code:
To build the sample code, just need to follow following commands

### Step 1:
Set the GOPATH as per your system and directory lets say `/home/user/test`   and create a sub directory under this by using following command

```
$mkdir -p src/github.com/IBM
```

### Step 2:
Change the directory to `src/github.com/IBM`  nad check out the code
```
git clone https://github.com/IBM/ibmcloud-storage-volume-lib.git
```

### Step 3:
Change the directory `$cd ibmcloud-storage-volume-lib`  and run the following command to build the sample code

```
$make deps
$make build
```

It will create `libSample` executable which you can run by using `./libSample` and you will see some options to validate lib functionalities.
