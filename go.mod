module github.com/bhbosman/goCommsNetDialer

go 1.24.0

require (
	github.com/bhbosman/goConnectionManager v0.0.0-20250308162041-c3d88e4a7878
	github.com/bhbosman/gocommon v0.0.0-20250308155359-4baa9bec452e
	github.com/bhbosman/gocomms v0.0.0-20250308192115-8af5b0178806
	go.uber.org/fx v1.23.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.37.0
	golang.org/x/sync v0.12.0
)

require (
	github.com/bhbosman/goCommsDefinitions v0.0.0-20250308144130-64993b60920c // indirect
	github.com/bhbosman/goerrors v0.0.0-20250307194237-312d070c8e38 // indirect
	github.com/bhbosman/gomessageblock v0.0.0-20250308073733-0b3daca12e3a // indirect
	github.com/bhbosman/goprotoextra v0.0.2 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cskr/pubsub v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/icza/gox v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/reactivex/rxgo/v2 v2.5.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/dig v1.18.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/cskr/pubsub => github.com/bhbosman/pubsub v1.0.3-0.20250308162024-50f212a35484
	github.com/golang/mock => github.com/bhbosman/gomock v1.6.1-0.20250308071159-4cf72f668c72
)
