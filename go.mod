module github.com/microsoft/moc-sdk-for-go

go 1.16

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20210608160410-67692ebc98de
	github.com/Azure/go-autorest/autorest v0.11.5
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/google/uuid v1.3.0
	github.com/microsoft/moc v0.11.0-alpha.9
	google.golang.org/grpc v1.54.0
	k8s.io/klog v1.0.0
)

require (
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Microsoft/go-winio v0.6.1
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88 // indirect
	google.golang.org/protobuf v1.30.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace (
	github.com/Azure/go-autorest v11.1.2+incompatible => github.com/Azure/go-autorest/autorest v0.10.0
	github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	github.com/miekg/dns => github.com/miekg/dns v1.1.25
	golang.org/x/net => golang.org/x/net v0.0.0-20220822230855-b0a4917ee28c
	golang.org/x/sys => golang.org/x/sys v0.0.0-20220823224334-20c2bfdbfe24

)
