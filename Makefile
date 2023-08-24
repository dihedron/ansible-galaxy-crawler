.PHONY: build
build:
	@go build

.PHONY: test
test: reset build
	@./ansible-galaxy-grabber --collections=@_test/input.json --directory=_test/download

.PHONY: reset
reset:
	@rm -rf _test/download