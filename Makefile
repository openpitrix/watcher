# Copyright 2019 The OpenPitrix Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.


login-dockerhub:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

build-image-%: ## build docker image
	@if [ "$*" = "latest" ];then \
	docker build -t openpitrix/watcher:latest .; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker build -t openpitrix/watcher:$* .; \
	fi

push-image-%: ## push docker image
	@if [ "$*" = "latest" ];then \
	docker push openpitrix/watcher:latest; \
	elif [ "`echo "$*" | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+"`" != "" ];then \
	docker push openpitrix/watcher:$*; \
	fi

build-push-image-%: ## build and push docker image
	make build-image-$*
	make login-dockerhub
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
