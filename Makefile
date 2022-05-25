
list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v ':=' | grep -v '^\.' | sed 's/:.*//g' | sed 's/://g' | sort

bootstrap:
	go mod init $(service)
	make init

init:
	go mod tidy

flow-token-test:
	go test ./example/test -run TestDepositFlowTokens -v 3

# todo: add more tests

test:
	make flow-token-test	