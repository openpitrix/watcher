# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

TARG.Name:=watcher
TRAG.Gopkg:=openpitrix.io/$(TRAG.Name)

build-image-%: ## build docker image
	@if [ "$*" = "latest" ];then \
	docker build -t $(TRAG.Gopkg):latest .; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t $(TRAG.Gopkg):$* .; \
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push $(TRAG.Gopkg):latest; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push $(TRAG.Gopkg):$*; \
	fi

build-push-image-%: ## build and push docker image
	make build-image-$*
	make push-image-$*

.PHONY: test
test: ## Run all tests
	make load-config-test
	make update-test
	@echo "test done"

load-config-test: ## Run test for LoadConf
	cd ./test && go test -v -run ^TestLoadConf$ && cd ..
	@echo "load-config-test done"

update-test: ## Run unit test for UpdateOpenPitrixEtcd
	cd ./test && go test -v -run ^TestWatchOpenPitrixConfig$ && cd ..
	@echo "update-openpitrix-etcd-test done"
