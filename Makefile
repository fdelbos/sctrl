test: mocks
	go test \
		-count=1 \
		-cover \
		-race \
		-timeout 60s

mocks:
	rm -fr ./mocks
	mockery --all --output ./mocks

.PHONY: mocks test