.PHONY: build
build:
	@go build

.PHONY: test
test: build
	@rm -rf _test/download && ./ansible-galaxy-grabber --collections=@_test/input.json --directory=_test/download