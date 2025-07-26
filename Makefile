# ==============================================================================
# Define dependencies

GOLANG          := golang:1.24
ALPINE          := alpine:3.22
POSTGRES        := postgres:17.5
KIND            := kindest/node:v1.33.1
GRAFANA         := grafana/grafana:11.6.0
PROMETHEUS      := prom/prometheus:v3.4.0
TEMPO           := grafana/tempo:2.7.0
LOKI            := grafana/loki:3.5.0
PROMTAIL        := grafana/promtail:3.5.0

KIND_CLUSTER    := iniciar-starter-cluster
NAMESPACE       := sales-system
SALES_APP       := sales
AUTH_APP        := auth
METRICS_APP     := metrics
BASE_IMAGE_NAME := localhost/iniciar
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/$(METRICS_APP):$(VERSION)

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize

dev-docker:
	docker pull $(GOLANG) & \
	docker pull $(ALPINE) & \
	docker pull $(KIND) & \
	docker pull $(POSTGRES) & \
	docker pull $(GRAFANA) & \
	wait;

# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"


statsviz:
	open http://localhost:3010/debug/statsviz

# ==============================================================================
# Building containers

build: sales auth  metrics




metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t $(METRICS_IMAGE) \
		--build-arg BUILD_REF=0.0.1 \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

auth:
	docker build \
		-f zarf/docker/dockerfile.auth \
		-t $(AUTH_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.
sales:
	docker build -f zarf/docker/dockerfile.sales \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner
	wait;

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# ==============================================================================
# Administration
pgcli:
	pgcli postgresql://postgres:postgres@localhost

# ------------------------------------------------------------------------------

IMAGE_LIST := \
	$(GOLANG) \
	$(ALPINE) \
	$(KIND) \
	$(POSTGRES) \
	$(GRAFANA) \
	$(PROMETHEUS) \
	$(TEMPO) \
	$(LOKI) \
	$(PROMTAIL)

load-images:
	@echo "Pulling Docker images..."
	@for image in $(IMAGE_LIST); do \
		docker pull $$image || exit 1; \
	done

	@echo "Saving images as tarballs..."
	@for image in $(IMAGE_LIST); do \
		name=$$(echo $$image | tr '/:' '__'); \
		docker save $$image -o $$name.tar || exit 1; \
	done

	@echo "Copying tarballs into kind node..."
	@for image in $(IMAGE_LIST); do \
		name=$$(echo $$image | tr '/:' '__'); \
		docker cp $$name.tar iniciar-starter-cluster-control-plane:/$$name.tar || exit 1; \
	done

	@echo "Importing images into containerd inside kind node..."
	@for image in $(IMAGE_LIST); do \
		name=$$(echo $$image | tr '/:' '__'); \
		docker exec iniciar-starter-cluster-control-plane ctr -n k8s.io images import /$$name.tar || exit 1; \
	done

	@echo "Cleaning up tarballs from kind node..."
	@for image in $(IMAGE_LIST); do \
		name=$$(echo $$image | tr '/:' '__'); \
		docker exec iniciar-starter-cluster-control-plane rm /$$name.tar || true; \
	done

	@echo "✅ All images loaded successfully into the kind cluster."

dev-load-db:
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(AUTH_IMAGE) --name $(KIND_CLUSTER) &\
	kind load docker-image $(METRICS_IMAGE) --name $(KIND_CLUSTER) &\
	wait;

dev-apply:
	kustomize build zarf/k8s/dev/grafana | kubectl apply -f -
	kustomize build zarf/k8s/dev/prometheus | kubectl apply -f -
	kustomize build zarf/k8s/dev/tempo | kubectl apply -f -
	kustomize build zarf/k8s/dev/loki | kubectl apply -f -
	kustomize build zarf/k8s/dev/promtail | kubectl apply -f -

	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/auth | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(AUTH_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready


dev-restart:
	kubectl rollout restart deployment $(AUTH_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)

dev-run: build dev-up dev-load dev-apply

dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

dev-database-restart:
	kubectl rollout restart statefulset database --namespace=$(NAMESPACE)

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run api/tooling/logfmt/main.go -service=$(SALES_APP)

dev-logs-auth:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(AUTH_APP) --all-containers=true -f --tail=100 | go run apis/tooling/logfmt/main.go

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

readiness:
	curl -i http://localhost:3000/v1/readiness


liveness:
	curl -i http://localhost:3000/v1/readiness


migrate:
	export SALES_DB_HOST=localhost; go run api/tooling/admin/main.go migrate

seed: migrate
	export SALES_DB_HOST=localhost; go run api/tooling/admin/main.go seed

compose-up:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml -p compose up -d

compose-build-up: build compose-up

compose-down:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml down

compose-logs:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml logs

pgcli:
	pgcli postgresql://postgres:postgres@localhost


token:
	curl -i \
	--user "admin@example.com:gophers" http://localhost:6000/v1/auth/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

curl-create:
	curl -i -X POST \
	-H "Authorization: Bearer ${TOKEN}" \
	-H 'Content-Type: application/json' \
	-d '{"name":"bill","email":"b@gmail.com","roles":["ADMIN"],"department":"ITO","password":"123","passwordConfirm":"123"}' \
	http://localhost:3000/v1/users

grafana:
	open http://localhost:3100/