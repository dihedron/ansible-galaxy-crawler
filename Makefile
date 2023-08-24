.PHONY: build
build:
	@go build

.PHONY: clean
clean:
	@rm -rf ansible-galaxy-grabber

.PHONY: test
test: build
	@rm -rf _test/download && ./ansible-galaxy-grabber --collections=@_test/input.json --directory=_test/download

.PHONY: reset
reset: clean
	@rm -rf _test/download