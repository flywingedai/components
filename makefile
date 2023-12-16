install:
	go build
	go install

mocks:
	export WD=$(PWD) && go generate ./...
	goimports -w .