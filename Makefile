GOCMD=go

test:
	$(GOCMD) test -cover -race ./...

bench:
	$(GOCMD) test -run=NONE -bench=. -benchmem ./...

.PHONY: test lint linters-install