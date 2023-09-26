BREW_PREFIX = $(shell brew --prefix)
DOMAIN = greenstar.local
OP_ACCOUNT = my.1password.com
REST_GENERATOR_SRC = $(shell find scripts/generate-rest)

$(BREW_PREFIX)/bin/gcloud:
	brew install -q google-cloud-sdk

$(BREW_PREFIX)/bin/op:
	brew install -q 1password-cli ]

$(BREW_PREFIX)/bin/kubectl:
	brew install -q kubernetes-cli

$(BREW_PREFIX)/bin/yq:
	brew install -q yq

$(BREW_PREFIX)/bin/gettext:
	# TODO: why do we need gettext?
	brew install -q gettext

$(BREW_PREFIX)/bin/go:
	brew install -q go

$(BREW_PREFIX)/bin/node:
	brew install -q node

$(BREW_PREFIX)/bin/skaffold:
	brew install -q skaffold

$(BREW_PREFIX)/bin/helm:
	brew install -q helm

$(BREW_PREFIX)/bin/kind:
	brew install -q kind

$(BREW_PREFIX)/bin/telepresence:
	brew install -q telepresenceio/telepresence/telepresence-oss

$(BREW_PREFIX)/bin/mkcert:
	brew install -q mkcert
	mkcert -install

$(BREW_PREFIX)/sbin/dnsmasq:
	# TODO: why do we need the "nss" package? who uses it and when? perhaps it was for setting up openssl to use mkcert?
	brew install -q nss dnsmasq

hack/currencyapi_key: $(BREW_PREFIX)/bin/op
	op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/currencyapi/api key" | xargs printf "%s" > hack/currencyapi_key

hack/descope_project_id: $(BREW_PREFIX)/bin/op
	op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/project_id" | xargs printf "%s" > hack/descope_project_id

hack/descope_management_key_id: $(BREW_PREFIX)/bin/op
	op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/management key id" | xargs printf "%s" > hack/descope_management_key_id

hack/descope_management_key_token: $(BREW_PREFIX)/bin/op
	op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/management key token" | xargs printf "%s" > hack/descope_management_key_token

hack/descope_backend_access_key_id: $(BREW_PREFIX)/bin/op
	op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/backend access key id" | xargs printf "%s" > hack/descope_backend_access_key_id

hack/descope_backend_access_key_token: $(BREW_PREFIX)/bin/op
	op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/backend access key token" | xargs printf "%s" > hack/descope_backend_access_key_token

deploy/local/ingress/wildcard-tls.crt deploy/local/ingress/wildcard-tls.key: $(BREW_PREFIX)/bin/mkcert
	mkcert \
		-cert-file deploy/local/ingress/wildcard-tls.crt \
		-key-file deploy/local/ingress/wildcard-tls.key \
		"$(DOMAIN)" "*.app.$(DOMAIN)" "*.$(DOMAIN)" "localhost" "127.0.0.1" "::1"

.PHONY: setup-local-dns
setup-local-dns: $(BREW_PREFIX)/sbin/dnsmasq
	grep -Fxq 'conf-dir=$(BREW_PREFIX)/etc/dnsmasq.d/,*.conf' '$(BREW_PREFIX)/etc/dnsmasq.conf' \
		|| echo 'conf-dir=$(BREW_PREFIX)/etc/dnsmasq.d/,*.conf' >> '$(BREW_PREFIX)/etc/dnsmasq.conf'
	grep -Fxq 'address=/.$(DOMAIN)/127.0.0.1' '$(BREW_PREFIX)/etc/dnsmasq.d/$(DOMAIN).conf' \
		|| printf 'address=/.$(DOMAIN)/127.0.0.1' > '$(BREW_PREFIX)/etc/dnsmasq.d/$(DOMAIN).conf'
	sudo "$(BREW_PREFIX)/bin/brew" services restart dnsmasq; \
	sudo mkdir -p "/etc/resolver"
	sudo bash -c 'printf "nameserver 127.0.0.1" > /etc/resolver/$(DOMAIN)'

.PHONY: hack/gcp-project-id
hack/gcp-project-id: $(BREW_PREFIX)/bin/gcloud
	gcloud config get-value project | xargs printf "%s" > hack/gcp-project-id

.PHONY: hack/gcp-project-number
hack/gcp-project-number: $(BREW_PREFIX)/bin/gcloud hack/gcp-project-id
	gcloud projects describe $(shell cat hack/gcp-project-id) --format='value(projectNumber)' | xargs printf "%s" > hack/gcp-project-number

hack/gcp-ksa-iam-principal: hack/gcp-project-number
	printf "principal://iam.googleapis.com/projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/subject/system:serviceaccount:greenstar:local-greenstar-backend" > hack/gcp-ksa-iam-principal

.PHONY: setup-gcp-workload-identity-pool
setup-gcp-workload-identity-pool: $(BREW_PREFIX)/bin/gcloud hack/gcp-project-id
	# TODO: this belongs in the project IaC since it's static and identical for all workstations
	gcloud iam workload-identity-pools describe greenstar --location="global" --quiet 2>/dev/null \
		|| gcloud iam workload-identity-pools create greenstar --project=$(shell cat hack/gcp-project-id) --location=global --description=GreenSTAR --display-name=GreenSTAR

.PHONY: grant-ksa-gcp-permissions
grant-ksa-gcp-permissions: $(BREW_PREFIX)/bin/gcloud hack/gcp-project-id hack/gcp-ksa-iam-principal
	# TODO: this belongs in the project IaC since it's static and identical for all workstations
	gcloud projects remove-iam-policy-binding $(shell cat hack/gcp-project-id) --member=$(cat hack/gcp-ksa-iam-principal) --role=roles/storage.objectViewer --quiet 2>/dev/null || true
	gcloud projects add-iam-policy-binding projects/$(shell cat hack/gcp-project-id) --role=roles/storage.objectViewer --member=$(shell cat hack/gcp-ksa-iam-principal) --condition=None

.PHONY: setup-local-cluster
setup-local-cluster hack/cluster-id hack/gcp-ksa-audience hack/gcp-credential-configuration.json: $(BREW_PREFIX)/bin/telepresence $(BREW_PREFIX)/bin/kind $(BREW_PREFIX)/bin/kubectl hack/gcp-project-id hack/gcp-project-number hack/gcp-ksa-iam-principal setup-gcp-workload-identity-pool deploy/local/ingress/wildcard-tls.crt deploy/local/ingress/wildcard-tls.key $(BREW_PREFIX)/bin/gcloud hack/currencyapi_key hack/descope_project_id hack/descope_management_key_id hack/descope_management_key_token hack/descope_backend_access_key_id hack/descope_backend_access_key_token
	# (Re)create cluster
	kind delete cluster --name=greenstar -q || true
	kind create cluster --name greenstar --config=deploy/local/kind-cluster-config.yaml --wait "1m"
	kubectl --context=kind-greenstar get namespace kube-system --template='{{.metadata.uid}}' | cut -c -32 | xargs printf "%s" > hack/cluster-id
	kubectl --context=kind-greenstar get --raw /openid/v1/jwks > hack/cluster-jwks.json

	# Setup Traefik
	curl -sSL "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.1.0/experimental-install.yaml" | kubectl apply -f -
	helm upgrade --install traefik traefik --repo "https://traefik.github.io/charts" -n ingress --create-namespace -f deploy/local/traefik-values.yaml --version v30.1.0 --wait
	kubectl apply -k deploy/local/ingress

	# Setup Telepresence
	telepresence helm install
	telepresence login

	# Save some useful information as local files
	printf "https://iam.googleapis.com/projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/providers/$$(cat hack/cluster-id)" > hack/gcp-ksa-audience

	# Ensure an OIDC workload identity pool provider exists for this local cluster
	gcloud iam workload-identity-pools providers delete $$(cat hack/cluster-id) --location=global --workload-identity-pool=greenstar --quiet 2>/dev/null || true
	gcloud iam workload-identity-pools providers create-oidc $$(cat hack/cluster-id) \
		--quiet \
		--location=global \
		--workload-identity-pool=greenstar \
		--issuer-uri="https://kubernetes.default.svc.cluster.local" \
		--attribute-mapping="google.subject=assertion.sub,attribute.namespace=assertion['kubernetes.io']['namespace'],attribute.service_account_name=assertion['kubernetes.io']['serviceaccount']['name'],attribute.pod=assertion['kubernetes.io']['pod']['name']" \
		--attribute-condition="assertion['kubernetes.io']['namespace'] in ['greenstar']" \
		--jwk-json-path="hack/cluster-jwks.json";
	gcloud iam workload-identity-pools create-cred-config \
		projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/providers/$$(cat hack/cluster-id) \
		--credential-source-file=/var/run/service-account/token \
		--credential-source-type=text \
		--output-file=hack/gcp-credential-configuration.json

	cp -v deploy/local/chart-local-values.yaml hack/chart-hack-values.yaml
	yq -i '.currencyAPI.key = load_str("./hack/currencyapi_key")' hack/chart-hack-values.yaml
	yq -i '.descope.project.id = load_str("./hack/descope_project_id")' hack/chart-hack-values.yaml
	yq -i '.descope.managementKey.id = load_str("./hack/descope_management_key_id")' hack/chart-hack-values.yaml
	yq -i '.descope.managementKey.key = load_str("./hack/descope_management_key_token")' hack/chart-hack-values.yaml
	yq -i '.descope.accessKey.id = load_str("./hack/descope_backend_access_key_id")' hack/chart-hack-values.yaml
	yq -i '.descope.accessKey.key = load_str("./hack/descope_backend_access_key_token")' hack/chart-hack-values.yaml
	yq -i '.configMaps[0].data["credential-configuration.json"] = load_str("./hack/gcp-credential-configuration.json")' hack/chart-hack-values.yaml
	yq -i '.configMaps[0].data["credential-configuration.json"] style="literal"' hack/chart-hack-values.yaml
	yq -i '.initJob.volumes[0].projected.sources[0].serviceAccountToken.audience = load_str("./hack/gcp-ksa-audience")' hack/chart-hack-values.yaml
	yq -i '.backend.volumes[0].projected.sources[0].serviceAccountToken.audience = load_str("./hack/gcp-ksa-audience")' hack/chart-hack-values.yaml

#	helm upgrade --install metrics-server metrics-server --repo https://kubernetes-sigs.github.io/metrics-server -n observability --create-namespace --set 'args={"--kubelet-insecure-tls"}' --version "3.12.1" --wait
#	helm upgrade --install otel-agent opentelemetry-collector --repo https://open-telemetry.github.io/opentelemetry-helm-charts -n observability --create-namespace -f ./deploy/local/otel-collector-agent-values.yaml --version "0.97.1" --wait
#	helm upgrade --install otel-gateway opentelemetry-collector --repo https://open-telemetry.github.io/opentelemetry-helm-charts -n observability --create-namespace -f ./deploy/local/otel-collector-gateway-values.yaml --version "0.97.1" --wait
#	helm upgrade --install prometheus prometheus --repo https://prometheus-community.github.io/helm-charts -n observability --create-namespace -f ./deploy/local/prometheus-values.yaml --version "25.24.1" --wait
#	helm upgrade --install cert-manager cert-manager --repo https://charts.jetstack.io -n cert-manager --create-namespace -f ./deploy/local/cert-manager-values.yaml --version "v1.15.1" --wait
#	kubectl apply -f ./deploy/local/jaeger-rbac-ingress.yaml # Jaeger Helm chart is missing those for some reason
#	helm upgrade --install jaeger-operator jaeger-operator --repo https://jaegertracing.github.io/helm-charts -n observability --create-namespace -f ./deploy/local/jaeger-values.yaml --version "2.54.0" --wait
#	kubectl apply -k ./deploy/local/grafana-dashboards
#	helm upgrade --install grafana grafana --repo https://grafana.github.io/helm-charts -n observability --create-namespace -f ./deploy/local/grafana-values.yaml --version "8.3.2" --wait
#	kubectl apply -k ./deploy/local/observability

hack/generate-rest-client: $(BREW_PREFIX)/bin/go $(REST_GENERATOR_SRC)
	go build -o hack/generate-rest-client ./scripts/generate-rest/cmd/client/main.go

hack/generate-rest-server: $(BREW_PREFIX)/bin/go $(REST_GENERATOR_SRC)
	go build -o hack/generate-rest-server ./scripts/generate-rest/cmd/server/main.go

.PHONY:
rest: hack/generate-rest-client hack/generate-rest-server
	./hack/generate-rest-client --api-file=api.yaml
	./hack/generate-rest-server --api-file=api.yaml

.PHONY: deploy
deploy: $(BREW_PREFIX)/bin/skaffold $(BREW_PREFIX)/bin/helm $(BREW_PREFIX)/bin/kubectl rest
	skaffold build -q | skaffold deploy --load-images --build-artifacts -

.PHONY: undeploy
undeploy: $(BREW_PREFIX)/bin/skaffold $(BREW_PREFIX)/bin/helm $(BREW_PREFIX)/bin/kubectl
	skaffold delete

.PHONY: dev
dev: $(BREW_PREFIX)/bin/skaffold $(BREW_PREFIX)/bin/helm $(BREW_PREFIX)/bin/kubectl rest
	skaffold dev --tail=false --toot=true

.PHONY: teardown-local-cluster
teardown-local-cluster: $(BREW_PREFIX)/bin/gcloud $(BREW_PREFIX)/bin/kind $(BREW_PREFIX)/bin/kubectl
	kind delete cluster --name greenstar || true
	gcloud iam workload-identity-pools providers delete $(shell cat hack/cluster-id) --location=global --workload-identity-pool=greenstar --quiet || true
	rm -fv hack/chart-hack-values.yaml || true
	rm -fv hack/gcp-credential-configuration.json || true
	rm -fv hack/gcp-ksa-audience || true
	rm -fv hack/gcp-ksa-iam-principal || true
	rm -fv hack/gcp-project-number || true
	rm -fv hack/gcp-project-id || true
	rm -fv hack/cluster-jwks.json || true
	rm -fv hack/cluster-id || true

.PHONY: telepresence-connect
telepresence-connect:
	telepresence connect --namespace greenstar

.PHONY: telepresence-disconnect
telepresence-disconnect: telepresence-leave-frontend telepresence-leave-backend
	telepresence uninstall --all-agents || true
	telepresence quit --stop-daemons || true

.PHONY: telepresence-intercept-backend
telepresence-intercept-backend:
	# TODO: re-enable mounting once sshfs is set up
	telepresence intercept local-greenstar-backend --env-file backend/telepresence.env --mount false --port 8080:80 --workload local-greenstar-backend

.PHONY: telepresence-leave-backend
telepresence-leave-backend:
	telepresence leave local-greenstar-backend || true

.PHONY: telepresence-intercept-frontend
telepresence-intercept-frontend:
	# TODO: re-enable mounting once sshfs is set up
	telepresence intercept local-greenstar-frontend --env-file frontend/telepresence.env --mount false --port 3000:80 --workload local-greenstar-frontend

.PHONY: telepresence-leave-frontend
telepresence-leave-frontend:
	telepresence leave local-greenstar-frontend || true
