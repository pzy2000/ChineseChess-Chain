GIT_BRANCH_NAME := $(shell git rev-parse --abbrev-ref HEAD)

build:
	go mod tidy && cd src && go build -o ../scripts/chainmaker-browser.bin

run:
	go mod tidy && cd src  && go build -o ../scripts/chainmaker-browser.bin  && cd ../ && ./scripts/chainmaker-browser.bin

run-dev:
	go mod tidy && cd src  && go build -o ../scripts/chainmaker-browser.bin  && cd ../ && ./scripts/chainmaker-browser.bin --env dev

docker-stop-mysql:
	cd docker && docker-compose -f docker-compose.yml down

docker-start-mysql:
	cd docker && docker-compose -f docker-compose.yml up -d

docker-stop-clickhouse-v:
	cd docker && docker-compose -f docker-compose-clickhouse.yml down -v

docker-stop-clickhouse:
	cd docker && docker-compose -f docker-compose-clickhouse.yml down

docker-start-clickhouse:
	cd docker && docker-compose -f docker-compose-clickhouse.yml up -d

docker-build:
	docker build . -t chainmaker-explorer-backend:latest -f Dockerfile

docker-push:
	docker image tag chainmaker-explorer-backend hub-dev.cnbn.org.cn/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}
	docker push hub-dev.cnbn.org.cn/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}
	docker image tag chainmaker-explorer-backend chainmaker1.tencentcloudcr.com/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}
	docker push chainmaker1.tencentcloudcr.com/opennet/chainmaker-explorer-backend:${GIT_BRACH_NAME}

ut:
	./scripts/ut_cover.sh

lint:
	golangci-lint run ./...

docker_kub_pod:
	kubectl get pod  -o wide | grep  explorer-backend

#算力平台测试链configma
k8s_configmap_test:
	kubectl create cm configmap-client --from-file=./k8s-test/configs/crypto-config/node1/user/client1/client1.key --dry-run=client -o yaml > ./k8s-test/configmap-client.yaml
	kubectl create cm configmap-config --from-file=./k8s-test/configs/ --dry-run=client -o yaml > ./k8s-test/configmap-config.yaml
	kubectl apply -f ./k8s-test/configmap-client.yaml
	kubectl apply -f ./k8s-test/configmap-config.yaml

#算力平台测试链
k8s_apply_test:
	kubectl delete -f ./k8s-test/deploy-explorer-backend.yaml
	kubectl apply -f ./k8s-test/deploy-explorer-backend.yaml

#联调测试链configmap
k8s_configmap_test_181:
	kubectl create cm configmap-client --from-file=./k8s-test-181/configs/crypto-config/node1/user/client1/client1.key --dry-run=client -o yaml > ./k8s-test-181/configmap-client.yaml
	kubectl create cm configmap-config --from-file=./k8s-test-181/configs/ --dry-run=client -o yaml > ./k8s-test-181/configmap-config.yaml
	kubectl apply -f ./k8s-test-181/configmap-client.yaml
	kubectl apply -f ./k8s-test-181/configmap-config.yaml

#联调测试链
k8s_apply_test_181:
	kubectl delete -f ./k8s-test-181/deploy-explorer-backend.yaml
	kubectl apply -f ./k8s-test-181/deploy-explorer-backend.yaml


#算力平台主链configmap
k8s_configmap:
	kubectl create cm configmap-client --from-file=./k8s/configs/crypto-config/node1/user/client1/client1.key --dry-run=client -o yaml > ./k8s/configmap-client.yaml
	kubectl create cm configmap-config --from-file=./k8s/configs/ --dry-run=client -o yaml > ./k8s/configmap-config.yaml
	kubectl apply -f ./k8s/configmap-client.yaml
	kubectl apply -f ./k8s/configmap-config.yaml

#算力平台主链
k8s_apply:
	kubectl delete -f ./k8s/deploy-explorer-backend.yaml
	kubectl apply -f ./k8s/deploy-explorer-backend.yaml

#联调主链configmap
k8s_configmap_181:
	kubectl create cm configmap-client --from-file=./k8s-181/configs/crypto-config/node1/user/client1/client1.key --dry-run=client -o yaml > ./k8s-181/configmap-client.yaml
	kubectl create cm configmap-config --from-file=./k8s-181/configs/ --dry-run=client -o yaml > ./k8s-181/configmap-config.yaml
	kubectl apply -f ./k8s-181/configmap-client.yaml
	kubectl apply -f ./k8s-181/configmap-config.yaml

#联调主链
k8s_apply_181:
	kubectl delete -f ./k8s-181/deploy-explorer-backend.yaml
	kubectl apply -f ./k8s-181/deploy-explorer-backend.yaml

#k8s重启前端服务
k8s_restart_front_formal:
	kubectl rollout restart deployment/explorer-front-formal
	kubectl rollout restart deployment/explorer-front-testnet

zip:
	cd k8s && zip -r ../explorer_config.zip ./configs && cd ../
