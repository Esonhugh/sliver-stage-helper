module github.com/Esonhugh/sliver-stage-helper

go 1.23.5

require (
	github.com/Binject/debug v0.0.0-20230508195519-26db73212a7a
	github.com/bishopfox/sliver v1.5.43
	github.com/google/uuid v1.6.0
	github.com/keystone-engine/keystone v0.0.0-20230519220125-fb92f32391c6
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	google.golang.org/grpc v1.64.1
	google.golang.org/protobuf v1.34.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
)

replace gvisor.dev/gvisor v0.0.0-20250321200759-3a9ba1735157 => github.com/google/gvisor v0.0.0-20250321200759-3a9ba1735157
