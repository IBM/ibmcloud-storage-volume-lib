# opensource-volume-lib-e2e

This is an implementation of code which is yaml based e2e execution for file and block storage libraries.


# Build the library

For building the e2e `GO` should be installed on the system

1. On your local machine, install [`Go`](https://golang.org/doc/install).
2. GO version should be >=1.16
3. Set the [`GOPATH` environment variable](https://github.com/golang/go/wiki/SettingGOPATH).
4. Build the library

   ## Clone the repo or your forked repo

   ```
   $ mkdir -p $GOPATH/src/github.com/IBM
   $ cd $GOPATH/src/github.com/IBM/
   $ git clone https://github.com/IBM/ibmcloud-storage-volume-lib.git
   $ cd ibmcloud-storage-volume-lib
   ```
   ## Build project and runs e2e testcases

  1. Edit e2e-final/config/vpc-config.toml and provide the respective vpc_volume_type, g2_riaas_endpoint_url, g2_resource_group_id and g2_api_key.
  2. Edit e2e-final/config/test-cases-block.yml and e2e-final/config/test-cases-file.yml and provide required zone,instanceID,clusterID, vpcID etc
  3. Set ENCRYPTION_KEY_CRN environment variable for running e2e based on encryptions else it will assume as provider managed encryption
  4. Add any new cases to the yaml and it will be picked up by the framework.

   ```
   $ make volume-lib-e2e-test
   ```


# How to contribute

If you have any questions or issues you can create a new issue [ here ](https://github.com/IBM/ibmcloud-storage-volume-lib/issues/new).

Pull requests are very welcome! Make sure your patches are well tested. Ideally create a topic branch for every separate change you make. For example:

1. Fork the repo

2. Create your feature branch (git checkout -b my-new-feature)

3. Commit your changes (git commit -am 'Added some feature')

4. Push to the branch (git push origin my-new-feature)

5. Create new Pull Request

6. Add the test results in the PR