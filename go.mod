module github.com/IBM/ibmcloud-storage-volume-lib

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/IBM-Cloud/ibm-cloud-cli-sdk v0.6.7
	github.com/IBM/ibmcloud-volume-file-vpc v1.0.0-beta1
	github.com/IBM/ibmcloud-volume-interface v1.0.0-beta8
	github.com/IBM/ibmcloud-volume-vpc v1.0.0-beta12
	github.com/fatih/structs v1.1.0
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/jarcoal/httpmock v1.0.8 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/prometheus/client_golang v1.8.0
	github.com/renier/xmlrpc v0.0.0-20170708154548-ce4a1a486c03 // indirect
	github.com/softlayer/softlayer-go v0.0.0-20181027013155-82a74c5bf7ff
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	gopkg.in/yaml.v2 v2.3.0
)

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190313205120-8b27c41bdbb1
