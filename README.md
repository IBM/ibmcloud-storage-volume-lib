# IBM Cloud storage common library

[![Build Status](https://travis-ci.com/IBM/ibmcloud-storage-volume-lib.svg?branch=master)](https://travis-ci.com/IBM/ibmcloud-storage-volume-lib)
[![Coverage](https://ibm.github.io/ibmcloud-storage-volume-lib/coverage/master/badge.svg)](https://ibm.github.io/ibmcloud-storage-volume-lib/coverage/master/cover.html)
[![e2e](https://ibm.github.io/ibmcloud-storage-volume-lib/coverage/master/e2e.svg)](https://travis-ci.com/IBM/ibmcloud-storage-volume-lib)

This library is for volume and snapshot management(create, delete, modify etc).
As of now this only have block snapshot and volume creation functionalities but going forward it will have more functionalities.

# Purpose
By using this library user can enable multiple volume providers by adding/modifying configuration.

# How to use
To use this library user need to check what all capabilities are supported in this library by checking `ibmcloud-storage-volume-lib/lib/provider/*_manager.go` files and concreate implementation for each provider which can be seen under `ibmcloud-storage-volume-lib/volume-providers` directory.

## Steps to use
Following steps has to followed by user to use this library
### Step 1
Update the glide.yaml to get the source code of this library, also only export the following packages to their application

```
"github.com/IBM/ibmcloud-storage-volume-lib/config"
"github.com/IBM/ibmcloud-storage-volume-lib/lib/provider"
"github.com/IBM/ibmcloud-storage-volume-lib/provider/local"
provider_util "github.com/IBM/ibmcloud-storage-volume-lib/provider/utils"
```

User can see the refrence from `ibmcloud-storage-volume-lib/samples`.

### Step 2
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
#### IAM secrets:
If you want to use your IAM credentials, please make sure the following properties are set in configuration file `ibmcloud-storage-volume-lib/etc/libconfig.toml`.  Replace value for `iam_api_key` and keep blank for `softlayer_username` and `softlayer_api_key`

```
[bluemix]
iam_url = "https://iam.bluemix.net"
iam_client_id = "bx"
iam_client_secret = "bx"
iam_api_key = "testIAM_KEY" # replace with IAM key and keep blank for APIkey auth
refresh_token = ""

[softlayer]
softlayer_block_enabled = true
softlayer_block_provider_name = "SOFTLAYER-BLOCK"
softlayer_file_enabled = false
softlayer_file_provider_name = "SOFTLAYER-FILE"
softlayer_username = ""# keep blank for IAM auth
softlayer_api_key = "" # keep blank for IAM auth
softlayer_endpoint_url = "https://api.softlayer.com/rest/v3"
softlayer_iam_endpoint_url = "https://api.softlayer.com/mobile/v3"
softlayer_datacenter = "dal12"
softlayer_api_timeout = "20s"

```


### Step 3
From the implementation file which uses use this library, user need to initialize the providers and open the sessions to backend IAaS provider and to do that user just need to call the following method from there

`github.com/IBM/ibmcloud-storage-volume-lib/provider/utils` utility packages, reference code can be found `ibmcloud-storage-volume-lib/main.go`

```
providerRegistry, err := provider_util.InitProviders(conf, logger)

and

sess, _, err := provider_util.OpenProviderSession(conf, providerRegistry, conf.Softlayer.SoftlayerBlockProviderName, logger)
```

`OpenProviderSession` call will be required for all enabled provider.
Once session is ready the user can call any method/interface which are part of `ibmcloud-storage-volume-lib/lib/provider/*_manager.go`. Sample code can be found `ibmcloud-storage-volume-lib/main.go`

## How to build the sample code
To build the sample code, just need to follow the following commands

### Step 1
Set the GOPATH as per your system and directory lets say `/home/user/test` and create a sub directory under this by using the following command

```
$mkdir -p src/github.com/IBM
```

Change the directory to `src/github.com/IBM` and check out the code

```
git clone https://github.com/IBM/ibmcloud-storage-volume-lib.git
```

### Step 3
Change the directory `$cd ibmcloud-storage-volume-lib` and run the following command to build the sample code

```
$make deps
$make build
```

It will create `libSample` executable which you can run by using `./libSample` and you will see some options to validate lib functionalities.
