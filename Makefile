SHELL=/bin/bash

BLACK     := $(shell tput setab 7)$(shell tput -Txterm setaf 0)
RED       := $(shell tput setab 0)$(shell tput -Txterm setaf 1)
GREEN     := $(shell tput setab 0)$(shell tput -Txterm setaf 2)
YELLOW    := $(shell tput setab 0)$(shell tput -Txterm setaf 3)
BLUE      := $(shell tput setab 4)$(shell tput setaf 7)
PURPLE    := $(shell tput setab 0)$(shell tput -Txterm setaf 5)
LIGHTBLUE := $(shell tput setab 0)$(shell tput -Txterm setaf 6)
WHITE     := $(shell tput setab 0)$(shell tput -Txterm setaf 7)

RESET := $(shell tput -Txterm sgr0)
CLEAR := -en "\033c"

PROJECTNAME=$(shell basename "$(PWD)")
# HAS_GIT := $(shell command -v  'ls -l .git' 2> /dev/null)
HAS_GIT := $(shell ls .git 2> /dev/null)
BRANCHENAME=
ifneq ("$(HAS_GIT)","")
	BRANCHENAME=$(shell git branch --show-current)
endif

TEXTSTART := $(GREEN)─────┤░░░▒▒▒▓▓▓$(RESET)
TEXTEND := $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME) $(GREEN)▓▓▓▒▒▒░░░├─────$(RESET)
READY := "\n$(TEXTSTART) Ready: $(TEXTEND)\n"

GOTEST_PRESENT := $(shell command -v gotest 2> /dev/null)
GOTEST=go test
ifneq ("$(GOTEST_PRESENT)","")
 GOTEST=gotest
endif

# Default make target
.DEFAULT_GOAL:=default

# Print help message with short target descriptions and check all tools are installed
default: help
	@echo ""

%::
	make
	@echo "$(RED) > type one of the targets above$(RESET)"
	@echo

help: ## Print this message
makefile: help
help: Makefile
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Choose a make command from the following: $(TEXTEND)"
	@echo
	@awk 'BEGIN {FS = ":.*##"; } /^[a-zA-Z_-]+:.*?##/ { printf "$(LIGHTBLUE)%-20s$(WHITE) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo
	@echo "GOPATH=$(GOPATH)"
	@echo "GOBIN=$(GOBIN)"
	@echo -e ${READY}

## colors: show all the colors
colors:
	@echo ${CLEAR}
	@echo "${BLACK}BLACK${RESET}"
	@echo "${RED}RED${RESET}"
	@echo "${GREEN}GREEN${RESET}"
	@echo "${YELLOW}YELLOW${RESET}"
	@echo "${BLUE}DARKBLUE${RESET}"
	@echo "${PURPLE}PURPLE${RESET}"
	@echo "${LIGHTBLUE}LIGHTBLUE${RESET}"
	@echo "${WHITE}WHITE${RESET}"
	@echo -e ${READY}

## install: install all project tools
install:
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Installing all project tools $(TEXTEND)"
	@echo

	@echo "> golangci-lint : a linter for checking source code"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo

	@echo "> identypo : finds typos in identifiers"
	@go install github.com/alexkohler/identypo/cmd/identypo@latest
	@echo

	@echo "> nakedret : a naked return linter"
	@go install github.com/alexkohler/nakedret@latest
	@echo

	@echo "> staticcheck : a static analysis linter"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo
	
	@echo "> gotest : a more colorfull test"
	@go install github.com/marcelloh/gotest@latest
	@echo

	@echo "> godoc : for looking at generated inline documentation"
	@go install golang.org/x/tools/cmd/godoc@latest
	@echo
	
	@echo -e ${READY}

## update: Updating mod-file and vendor directory
update:
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Updating mod-file $(TEXTEND)"
	@echo
	@rm -f go.sum

	@go env -w GOPRIVATE=github.com/zale144
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod tidy
	@echo -e ${READY}

# Build all
.PHONY: all
all: ## Runs following targets: test, lint
	@echo ${CLEAR}
	@echo "$(GREEN)     ┏━━━━━━━━━━━┓"
	@echo "$(GREEN)═════┃ $(RESET)Start All$(GREEN) ┃═════$(RESET)"
	@echo "$(GREEN)     ┗━━━━━━━━━━━┛"
	@echo 

	@make test_next
	@make lint_next

	@echo 
	@echo "$(GREEN)     ┏━━━━━━━━━━━━━━┓"
	@echo "$(GREEN)═════┃ $(RESET)All is Ready$(GREEN) ┃═════$(RESET)"
	@echo "$(GREEN)     ┗━━━━━━━━━━━━━━┛"

.PHONY: test
test: ## Runs all test targets (no race condition check, only short tests)
	@echo ${CLEAR}
	@make test_next
test_next:
	@echo "$(TEXTSTART) Runs all test targets (withrace condition check, only short tests) $(TEXTEND)"
	@echo
	@$(GOTEST) -failfast ./... -short -count=1 -race
	@echo -e ${READY}

.PHONY: test-g
test-g: ## generate all testdata
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Generate all test-data $(TEXTEND)"
	@echo
	@$(GOTEST) -failfast ./... --tags=generatetests
	@echo -e ${READY}

.PHONY: test-d
test-d: ## run all tests verbose within a directory (Exmaple: make test-d d=./interfaces/controllers)
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Testing directory $(WHITE)'$(d)' $(TEXTEND)"
	@echo
ifdef d
	@$(GOTEST) $(VENDOR) -v -count=1 $(TESTTIMEOUT) ${d}/...
else
	@echo "$(RED)d=??? not given$(RESET)"
endif
	@echo -e ${READY}

.PHONY: test-p
test-p: ## run all tests verbose within a package (Exmaple: make test-p p=handler)
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Testing package $(TEXTEND)"
	@echo
ifdef p
	@$(GOTEST) $(VENDOR) -v -count=1 -coverprofile coverage.out ./$(p)
	@go tool cover -func=coverage.out | tail -n 1
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover -html coverage.out
else
	@echo "$(RED)p=??? not given$(RESET)"
endif
	@echo -e ${READY}

.PHONY: test-v
test-v: ## Runs all test targets verbose (no race condition check, only short tests)
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Runs all test targets verbose (no race condition check, only short tests) $(TEXTEND)"
	@echo
	@$(GOTEST) -v -failfast ./... -short -count=1
	@echo -e ${READY}

.PHONY: test-c
test-c: ## Runs all test targets for codecoverage (no race condition check, only short tests)
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Runs all test targets verbose (no race condition check, only short tests) with code codecoverage $(TEXTEND)"
	@echo
	@$(GOTEST) -v -failfast -short -count=1 -coverprofile coverage.out ./...
	@echo "Coverage overall:"
	@go tool cover -func=coverage.out | tail -n 1
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go tool cover -html coverage.out
	@echo -e ${READY}

.PHONY: test-h
test-h: ## Runs all test targets and creates heatmap
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Runs all test targets and creates heatmap $(TEXTEND)"
	@echo
	@$(GOTEST) -failfast -short -count=1 -coverprofile coverage.out ./...
	@go-cover-treemap -coverprofile coverage.out > heatmap.svg
	@echo -e ${READY}

.PHONY: lint
lint: ## Runs the linter (uses .golangci.yml with 36 rules) https://golangci-lint.run/usage/install/ @ gopath/bin
	@echo ${CLEAR}
	@make lint_next
lint_next:
	@echo "$(TEXTSTART) Checking linters $(TEXTEND)"
	@echo
	@golangci-lint --version
	@golangci-lint run ./...
	@echo "identypo"
	@identypo ./...
	@echo "nakedret"
	@nakedret ./...
	@echo "gosec"
	@gosec -fmt=golint -quiet ./...
	@staticcheck -version
	@staticcheck -checks all ./...
	@echo -e ${READY}

.PHONY: doc
doc: ## look at documentation
	@echo ${CLEAR}
	@echo "$(TEXTSTART) Checking with golangci-lint $(TEXTEND)"
	@echo
	@echo ${CLEAR}
	@echo "$(LIGHTBLUE) > Starting DOC $(BLUE) $(PROJECTNAME) $(YELLOW) $(BRANCHENAME)$(RESET)"
	@open http://localhost:6060/pkg/ube/
	@echo "$(PURPLE)CTRL-c to exit$(RESET)"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) godoc -http=:6060
	@echo -e ${READY}
