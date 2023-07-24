.PHONY: install format test clean build run share

# Set default env as develop.
ifeq ($(ENV),)
ENV := develop
endif


# Git revision is a short hash of the current commit.
# Mainly used for DevOps purpose, to manage entire services.
GIT_REVISION := $(shell git rev-parse --short HEAD)
GIT_TAG := $(shell git describe --tags)

# Alias command for docker-compose.
COMPOSE := ENV=$(ENV) GIT_REVISION=$(GIT_REVISION) GIT_TAG=$(GIT_TAG) \
	COMPOSE_HTTP_TIMEOUT=300 docker-compose


up: ## Boot all containers
	$(COMPOSE) up -d

build:
	$(COMPOSE) build --no-cache

build-one:
	$(COMPOSE) build --no-cache $(tag)

down: ## Kill all containers
	$(COMPOSE) kill

# config: ## Pull configs from Vault service
# 	ENV=${ENV} cd infra/vault && $(MAKE) build &&  $(MAKE) pull

# share: ## Pull configs from Vault service
# 	ENV=${ENV} cd infra/vault && $(MAKE) share
