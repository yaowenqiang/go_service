> cd $GOPATH/pkg/
>  mod
> go clean --modcache
> rm -rf sumdb

> GOCACHE
> GOMODCACHE

> export GO111MODULE=on
> export GOPROXY=https://goproxy.io
> go mod init github.com/yaowenqiang/service
> go mod tidy

## cached here

> $GOMODCACH/github.com/pkg/errors@v0.9.1/

## private mod mirror

> https://github.com/gomods/athens
> https://jfrog.com/artifactory/

> GONOPROXY

> export GONOPROXY="github.com"
