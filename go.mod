module github.com/opensourceways/robot-gitee-approve

go 1.15

replace k8s.io/apimachinery => k8s.io/apimachinery v0.17.3

require (
	github.com/opensourceways/community-robot-lib v0.0.0-20220118064921-28924d0a1246
	github.com/opensourceways/go-gitee v0.0.0-20220120022149-6d34985edf4f
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)
