module github.com/kevin0120/GoScrewdriverWebApi

go 1.16

require (
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/atomic v1.11.0
	golang.org/x/net v0.17.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/resty.v1 v1.12.0
)

replace go.uber.org/atomic => github.com/uber-go/atomic v1.9.0
