# This file is part of myhttp

# Licensed under the BSD-3-Clause license:
# https://opensource.org/licenses/BSD-3-Clause
# Copyright (c) 2020, Pablo S. Blum de Aguiar

# list all available targets
list:
	@sh -c "$(MAKE) -p no_targets__ | awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ {split(\$$1,A,/ /);for(i in A)print A[i]}' | grep -v '__\$$' | grep -v 'make\[1\]' | grep -v 'Makefile' | sort"
.PHONY: list

# required for list
no_targets__:

# run tests in verbose mode
test:
	@go test -v
.PHONY: test

# run tests collecting coverage statistics and displays the result
cover:
	@go test -coverprofile=coverage.out -covermode=count
	@go tool cover -func=coverage.out
.PHONY: cover

# display an HTML representation of the test coverage
cover-html: cover
	@go tool cover -html=coverage.out
.PHONY: cover-html

# display an HTML representation of the test coverage
build:
	@go build -x -v -race
.PHONY: build
