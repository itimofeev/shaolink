

run-db:
	docker run -d --network=ac_test -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 --name db postgres:12.4

lint:
	# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0
	GO111MODULE=on GL_DEBUG=debug L_DEBUG=linters_output GOPACKAGESPRINTGOLISTERRORS=1 golangci-lint -v run

goimports:
	#go get golang.org/x/tools/cmd/goimports
	goimports -w -local github.com/itimofeev/shaolink internal cmd
	go mod tidy