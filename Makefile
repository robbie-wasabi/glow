
list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v ':=' | grep -v '^\.' | sed 's/:.*//g' | sed 's/://g' | sort

bootstrap:
	go mod init $(service)
	make init

init:
	go mod tidy


test:
	./test.sh