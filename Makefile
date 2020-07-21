TEST = go test -v -failfast -cover ./...

all: install

install:
	go get -u golang.org/x/lint/golint
	go get github.com/fzipp/gocyclo
	go get github.com/securego/gosec/cmd/gosec/...

lint:
	go fmt -x
	gocyclo -over 10 .
	golint
	go vet
	gosec ./...

run:
	go run synacor.go

test:
	$(TEST)

test-run:
ifdef test
	$(TEST) -run $(test)
else
	@echo Syntax is 'make $@ test=<test name>'
endif
