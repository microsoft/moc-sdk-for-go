module github.com/microsoft/moc-sdk-for-go

go 1.14

require (
	github.com/Azure/go-autorest/autorest v0.9.0
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/google/uuid v1.2.0
	github.com/hectane/go-acl v0.0.0-20190604041725-da78bae5fc95 // indirect
	github.com/microsoft/moc v0.10.10-0.20210510205040-161b5e608074
	github.com/microsoft/moc-pkg v0.10.9-alpha.4.0.20210511211911-84d67d21f73c // indirect
	github.com/spf13/viper v1.6.2
	google.golang.org/grpc v1.27.1
	google.golang.org/protobuf v1.25.0 // indirect
	k8s.io/klog v1.0.0
)

replace github.com/Azure/go-autorest v11.1.2+incompatible => github.com/Azure/go-autorest/autorest v0.10.0
