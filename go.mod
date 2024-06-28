module github.com/microsoft/moc-sdk-for-go

go 1.16

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20210608160410-67692ebc98de
	github.com/Azure/go-autorest/autorest v0.11.29
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/google/uuid v1.4.0
	github.com/microsoft/moc v0.18.1
	google.golang.org/grpc v1.59.0
	k8s.io/klog v1.0.0
)

require (
	github.com/Microsoft/go-winio v0.6.1
	github.com/golang/protobuf v1.5.4
	github.com/spf13/viper v1.18.2
	github.com/stretchr/testify v1.8.4
	google.golang.org/protobuf v1.34.2
)

replace (
	github.com/Azure/go-autorest v11.1.2+incompatible => github.com/Azure/go-autorest/autorest v0.10.0
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/microsoft/moc => ../moc
	github.com/miekg/dns => github.com/miekg/dns v1.1.25
	github.com/nats-io/nkeys => github.com/nats-io/nkeys v0.4.6
	golang.org/x/net => golang.org/x/net v0.0.0-20220822230855-b0a4917ee28c
	golang.org/x/sys => golang.org/x/sys v0.0.0-20220823224334-20c2bfdbfe24
)
