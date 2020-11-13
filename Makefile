# Copyright (C) 2020 Yunify, Inc.

VERBOSE = no

.PHONY: help
help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "  format  to format the code"
	@echo "  vet     to run golang vet"
	@echo "  lint    to run the staticcheck"
	@echo "  check   to format, vet, lint"
	@echo "  test    to run test case"
	@echo "  bench   to run benchmark test case"
	@exit 0


.PHONY: test
test:
	@[[ ${VERBOSE} = "yes" ]] && set -x; go test -race -v . -test.count=1 -failfast

.PHONY: bench
bench:
	@[[ ${VERBOSE} = "yes" ]] && set -x; go test -test.bench="." -test.run="Benchmark" -benchmem -count=1


.PHONY: format
format:
	@[[ ${VERBOSE} = "yes" ]] && set -x; go fmt ./...;

.PHONY: vet
vet:
	@[[ ${VERBOSE} = "yes" ]] && set -x; go vet ./...;

.PHONY: lint
lint:
	@[[ ${VERBOSE} = "yes" ]] && set -x; staticcheck ./...;


.PHONY: check
check: format vet lint

.DEFAULT_GOAL = help

# Target name % means that it is a rule that matches anything, @: is a recipe;
# the : means do nothing
%:
	@:
