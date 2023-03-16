module github.com/microsoft/moc-sdk-for-go

go 1.20

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20210608160410-67692ebc98de
	github.com/Azure/go-autorest/autorest v0.9.0
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/google/uuid v1.3.0
	github.com/microsoft/moc v0.10.28-alpha.2
	github.com/spf13/viper v1.15.0
	google.golang.org/grpc v1.52.0
	k8s.io/klog v1.0.0
)

require (
	github.com/Azure/go-autorest/autorest/adal v0.5.0 // indirect
	github.com/Azure/go-autorest/logger v0.1.0 // indirect
	github.com/Azure/go-autorest/tracing v0.5.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hectane/go-acl v0.0.0-20190604041725-da78bae5fc95 // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20221227171554-f9683d7f8bef // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Azure/go-autorest v11.1.2+incompatible => github.com/Azure/go-autorest/autorest v0.10.0
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/microsoft/moc => github.com/microsoft/moc v0.11.0-alpha.1.0.20230316182544-e188475f87dc
	github.com/miekg/dns => github.com/miekg/dns v1.1.25
)
