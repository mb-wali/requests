all: requests

install-swagger:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger.json: install-swagger
	GO111MODULE=on go mod vendor && GO111MODULE=off swagger generate spec -o ./swagger.json --scan-models

requests: swagger.json
	go build

clean:
	rm -rf swagger.json requests vendor

.PHONY: install-swagger clean all
