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
> GONOSUM
> GOPRIVATE

> export GONOPROXY="github.com"
> export GOPRIVATE="github.com"

> MVS 

> go mod vendor
> go get github.com/divan/expvarmon

> brew install  staticcheck

# pprof
> hey -c 10 -n 15000 http://localhost:8000/v1/products
> go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
> (pprof) top
> (pprof) top -cum
> (pprof) web

> go get github.com/divan/expvarmon
> 

# kubernetes


> brew install kind
> brew install kustomize

> docker pull kindest/node:v1.20.0
> kubectl version --client

> docker pull postgres:13-alpine
> docker pull openzipkin/zipkin:2.23
> docker pull alpine:3.12.3

# test 
> go test -cover
> go test -coverprofile cover.out
> go tool cover -html cover.out
