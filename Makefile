
list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v ':=' | grep -v '^\.' | sed 's/:.*//g' | sed 's/://g' | sort

init:
	go mod tidy

bootstrap:
	go mod init $(service)
	make init

test:
	./test.sh

# clean cache
cache-clean:
	go clean -cache
