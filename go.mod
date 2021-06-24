module github.com/microsoft/moc-sdk-for-go

go 1.15

require (
	code.cloudfoundry.org/bytefmt v0.0.0-20210608160410-67692ebc98de
	github.com/Azure/go-autorest/autorest v0.9.0
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/google/uuid v1.2.0
	github.com/microsoft/moc v0.10.11-alpha.4
	github.com/spf13/viper v1.7.1
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.27.1
	k8s.io/klog v1.0.0
)

replace (
	github.com/Azure/go-autorest v11.1.2+incompatible => github.com/Azure/go-autorest/autorest v0.10.0
	github.com/miekg/dns => github.com/miekg/dns v1.1.25
)
