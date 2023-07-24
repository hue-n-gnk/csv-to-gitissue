.PHONY: install format test clean build run

# Set default env as develop.
ifeq ($(ENV),)
ENV := develop
endif

# Set current Git revision unless given explicitly.
ifeq ($(GIT_REVISION),)
GIT_REVISION := $(shell git rev-parse --short HEAD)
endif

-include ./.env.${ENV}

BUILD := ./build

SCHEMA_DIR := ../schema/build

PACKAGE_NAME := hue-n-gnk/miro-github-issue

API_DOC_DIR := docs
API_DOC := $(API_DOC_DIR)/api.yml

REST_DIR := api
REST_DTOS_DIR := dtos

DAO_DIR := database/dao

GORUN := GIT_REVISION=$(GIT_REVISION) go run

# DB_GEN := gentool -db postgres -dsn "host=${SECRET_POSTGRES_HOSTNAME} port=${SECRET_POSTGRES_PORT} \
# 	user=${SECRET_POSTGRES_USER}  password=${SECRET_POSTGRES_PASSWORD} dbname=${SECRET_POSTGRES_DATABASE} sslmode=disable" \
# 	-fieldNullable -fieldWithIndexTag -fieldWithTypeTag -outPath ${DAO_DIR} -fieldSignable -withUnitTest -modelPkgName entities

TEST := GIT_REVISION=$(GIT_REVISION) go test -cover -coverprofile coverage.log

TEST_TOOL := go tool cover -html=coverage.log

DB_GEN := cd script/dbgen/ && go run .

WORKER_DIR := workers

define _apigen
	oapi-codegen -generate chi-server -include-tags="$(1)" -package="$(1)" $(API_DOC) > "$(REST_DIR)/$(1)/$(1).gen.go"
endef

define _apigen_w_query
	arg=$(1) oapi-codegen -config="docs/cfg.yml" $(API_DOC)
endef

define _jsongen
	go run helpers/transformengine/jsonfiles/genaratejsontools.go -in=helpers/transformengine/jsonfiles/rawjson/$(2).json -out=helpers/transformengine/jsonfiles/$(1)/$(2)
endef

define _genconfig
	echo "" > $(API_DOC_DIR)/cfg.yml
	echo "package: $(1)" >> $(API_DOC_DIR)/cfg.yml
	echo "generate:" >> $(API_DOC_DIR)/cfg.yml
	echo "  chi-server: true" >> $(API_DOC_DIR)/cfg.yml
	echo "output-options:" >> $(API_DOC_DIR)/cfg.yml
	echo "  include-tags:" >> $(API_DOC_DIR)/cfg.yml
	echo "    - $(1)" >> $(API_DOC_DIR)/cfg.yml
	echo "additional-imports:" >> $(API_DOC_DIR)/cfg.yml
	echo "  - alias: ." >> $(API_DOC_DIR)/cfg.yml
	echo "    package: $(PACKAGE_NAME)/$(REST_DTOS_DIR)" >> $(API_DOC_DIR)/cfg.yml
	echo "output: $(REST_DIR)/$(1)/$(1).gen.go" >> $(API_DOC_DIR)/cfg.yml
endef

define _build
	$(eval target := $(abspath $(strip $1)))
	$(eval output := $(abspath $(BUILD)/$(strip $2)))

	@echo "Compiling $(target) --> $(output)"
	cd $(target) && GOOS=linux GOARCH=amd64 go build -o $(output) $(target)
	@echo "Done!\n"
endef

run:
	air

format:
	go fmt ./...

dbgen:
	@echo "Generating entity models, DAO from database"
	$(DB_GEN)
	@echo "Generate completed"

api-doc:
	@echo "Copying api doc.."
	@mkdir -p $(API_DOC_DIR)
	@\cp -rf $(SCHEMA_DIR)/index.yaml ./$(API_DOC)
	@echo "Success!"


dtogen: api-doc
	@echo "Start modelgen..."
	@mkdir -p $(REST_DTOS_DIR)
	oapi-codegen -generate "types,skip-prune" -package="dtos" $(API_DOC) > "$(REST_DTOS_DIR)/dto.gen.go"
	@echo "Success codegen! Check it out $(REST_DTOS_DIR)/dto.gen.go"

apigen-help:
	@echo ""
	@echo "  Syntax:"
	@echo "  make apigen tag=<api tag> query=1"
	@echo ""
	@echo "  Example:"
	@echo "  make apigen tag=health"
	@echo "  If query=1, that means the gen will imports dto package to apigen"
	@echo ""

apigen: api-doc
	@[ "${tag}" ] || ( make apigen-help ; exit 1 )
	@echo "Start apigen $(tag) ..."
	mkdir -p $(REST_DIR)/$(tag)
	@echo "Generate yml config file to $(API_DOC_DIR)"
	@$(call _genconfig,$(tag))

ifeq ($(query),1)
	@$(call _apigen_w_query,$(tag))
else
	@$(call _apigen,$(tag))
endif
	@echo "Success generate API! Check it out $(REST_DIR)/$(tag)/$(tag).gen.go"
test:
	$(TEST) ./... && $(TEST_TOOL)

install:
	go mod download

clean:
	rm -Rf $(BUILD)/

# Build Go artifacts including api and workers.
# The artifacts is for Linux x86_64 environment, and the directory structure
# must be in sync with Dockerfile.hisys.release.
#
# Directory Structure:
#   build/
#     api
#     fetch-worker
#     delivery-worker
build: clean
	@$(call _build, "", api)
	@$(call _build, workers/fetchexam, fetch-worker)
	@$(call _build, workers/deliverychat, delivery-worker)

