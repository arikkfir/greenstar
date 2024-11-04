BREW_PREFIX = $(shell brew --prefix)
DOMAIN = greenstar.local
OP_ACCOUNT = my.1password.com
REST_GENERATOR_SRC = $(shell find scripts/generate-rest)

.PHONY: setup-toolchain
setup-toolchain:
	# TODO: why do we need the "nss" package? who uses it and when? perhaps it was for setting up openssl to use mkcert?
	# TODO: why do we need gettext?
	brew install -q google-cloud-sdk 1password-cli kubernetes-cli yq gettext go node skaffold helm kind telepresenceio/telepresence/telepresence-oss mkcert nss dnsmasq
	mkcert -install

hack/generate-rest-client: $(REST_GENERATOR_SRC)
	go build -o hack/generate-rest-client ./scripts/generate-rest/cmd/client/main.go

hack/generate-rest-server: $(REST_GENERATOR_SRC)
	go build -o hack/generate-rest-server ./scripts/generate-rest/cmd/server/main.go

.PHONY:
rest: hack/generate-rest-client hack/generate-rest-server
	./hack/generate-rest-client --api-file=api.yaml
	./hack/generate-rest-server --api-file=api.yaml

hack/currencyapi_key:
	if [ -z "$(CURRENCYAPI_KEY)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/currencyapi/api key" | xargs printf "%s" > hack/currencyapi_key; \
  	else \
		printf "$(CURRENCYAPI_KEY)" > hack/currencyapi_key; \
	fi

hack/descope_project_id:
	if [ -z "$(DESCOPE_PROJECT_ID)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/project_id" | xargs printf "%s" > hack/descope_project_id; \
  	else \
		printf "$(DESCOPE_PROJECT_ID)" > hack/descope_project_id; \
	fi

hack/descope_management_key_id:
	if [ -z "$(DESCOPE_MANAGEMENT_KEY_ID)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/management key id" | xargs printf "%s" > hack/descope_management_key_id; \
  	else \
		printf "$(DESCOPE_MANAGEMENT_KEY_ID)" > hack/descope_management_key_id; \
	fi

hack/descope_management_key_token:
	if [ -z "$(DESCOPE_MANAGEMENT_KEY_TOKEN)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/management key token" | xargs printf "%s" > hack/descope_management_key_token; \
  	else \
		printf "$(DESCOPE_MANAGEMENT_KEY_TOKEN)" > hack/descope_management_key_token; \
	fi

hack/descope_backend_access_key_id:
	if [ -z "$(DESCOPE_BACKEND_ACCESS_KEY_ID)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/backend access key id" | xargs printf "%s" > hack/descope_backend_access_key_id; \
  	else \
		printf "$(DESCOPE_BACKEND_ACCESS_KEY_ID)" > hack/descope_backend_access_key_id; \
	fi

hack/descope_backend_access_key_token:
	if [ -z "$(DESCOPE_BACKEND_ACCESS_KEY_TOKEN)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/GreenSTAR/descope/backend access key token" | xargs printf "%s" > hack/descope_backend_access_key_token; \
  	else \
		printf "$(DESCOPE_BACKEND_ACCESS_KEY_TOKEN)" > hack/descope_backend_access_key_token; \
	fi

hack/gcp-project-id:
	if [ -z "$(GCP_PROJECT_ID)" ]; then \
		gcloud config get-value project | xargs printf "%s" > hack/gcp-project-id; \
	else \
		printf "$(GCP_PROJECT_ID)" > hack/gcp-project-id; \
	fi

hack/gcp-project-number: hack/gcp-project-id
	if [ -z "$(GCP_PROJECT_NUMBER)" ]; then \
		gcloud projects describe $(shell cat hack/gcp-project-id) --format='value(projectNumber)' | xargs printf "%s" > hack/gcp-project-number; \
	else \
		printf "$(GCP_PROJECT_NUMBER)" > hack/gcp-project-number; \
	fi

hack/newrelic-license-key:
	if [ -z "$(NEWRELIC_LICENSE_KEY)" ]; then \
		op read --account $(OP_ACCOUNT) "op://Development/Newrelic/license keys/key1" | xargs printf "%s" > hack/newrelic_license_key; \
	else \
		printf "$(NEWRELIC_LICENSE_KEY)" > hack/newrelic-license-key; \
	fi

hack/gcp-ksa-iam-principal: hack/gcp-project-number
	printf "principal://iam.googleapis.com/projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/subject/system:serviceaccount:greenstar:local-greenstar-backend" > hack/gcp-ksa-iam-principal

deploy/local/ingress/wildcard-tls.crt deploy/local/ingress/wildcard-tls.key:
	mkcert \
		-cert-file deploy/local/ingress/wildcard-tls.crt \
		-key-file deploy/local/ingress/wildcard-tls.key \
		"$(DOMAIN)" "*.app.$(DOMAIN)" "*.$(DOMAIN)" "localhost" "127.0.0.1" "::1"

.PHONY: setup-local-dns
setup-local-dns:
	grep -Fxq 'conf-dir=$(BREW_PREFIX)/etc/dnsmasq.d/,*.conf' '$(BREW_PREFIX)/etc/dnsmasq.conf' \
		|| echo 'conf-dir=$(BREW_PREFIX)/etc/dnsmasq.d/,*.conf' >> '$(BREW_PREFIX)/etc/dnsmasq.conf'
	grep -Fxq 'address=/.$(DOMAIN)/127.0.0.1' '$(BREW_PREFIX)/etc/dnsmasq.d/$(DOMAIN).conf' \
		|| printf 'address=/.$(DOMAIN)/127.0.0.1' > '$(BREW_PREFIX)/etc/dnsmasq.d/$(DOMAIN).conf'
	sudo "$(BREW_PREFIX)/bin/brew" services restart dnsmasq; \
	sudo mkdir -p "/etc/resolver"
	sudo bash -c 'printf "nameserver 127.0.0.1" > /etc/resolver/$(DOMAIN)'

.PHONY: setup-local-cluster
setup-local-cluster: hack/newrelic-license-key hack/gcp-project-number hack/currencyapi_key hack/descope_project_id hack/descope_management_key_id hack/descope_management_key_token hack/descope_backend_access_key_id hack/descope_backend_access_key_token deploy/local/ingress/wildcard-tls.crt deploy/local/ingress/wildcard-tls.key
	# Create cluster
	kind create cluster --name greenstar --config=deploy/local/kind-cluster-config.yaml --wait "1m"
	kubectl --context=kind-greenstar get namespace kube-system --template='{{.metadata.uid}}' | cut -c -32 | xargs printf "%s" > hack/cluster-id
	kubectl --context=kind-greenstar get --raw /openid/v1/jwks > hack/cluster-jwks.json

	# Setup ingress
	curl -sSL "https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.1.0/experimental-install.yaml" | kubectl apply -f -
	helm upgrade --install traefik traefik --repo "https://traefik.github.io/charts" -n ingress --create-namespace -f deploy/local/traefik-values.yaml --version v30.1.0 --wait
	kubectl apply -k deploy/local/ingress

	# Build KSA audience
	printf "https://iam.googleapis.com/projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/providers/$$(cat hack/cluster-id)" > hack/gcp-ksa-audience

	# Setup workload identity pool provider for this cluster
	gcloud iam workload-identity-pools providers delete $$(cat hack/cluster-id) --location=global --workload-identity-pool=greenstar --quiet 2>/dev/null || true
	gcloud iam workload-identity-pools providers create-oidc $$(cat hack/cluster-id) \
    		--quiet \
    		--location=global \
    		--workload-identity-pool=greenstar \
    		--issuer-uri="https://kubernetes.default.svc.cluster.local" \
    		--attribute-mapping="google.subject=assertion.sub,attribute.namespace=assertion['kubernetes.io']['namespace'],attribute.service_account_name=assertion['kubernetes.io']['serviceaccount']['name'],attribute.pod=assertion['kubernetes.io']['pod']['name']" \
    		--attribute-condition="assertion['kubernetes.io']['namespace'] in ['greenstar']" \
    		--jwk-json-path="hack/cluster-jwks.json"

	gcloud iam workload-identity-pools create-cred-config \
    		projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/providers/$$(cat hack/cluster-id) \
    		--credential-source-file=/var/run/service-account/token \
    		--credential-source-type=text \
    		--output-file=hack/gcp-credential-configuration.json

	# Build chart values file
	cp -v deploy/local/chart-local-values.yaml hack/chart-hack-values.yaml
	yq -i '.currencyAPI.key = load_str("./hack/currencyapi_key")' hack/chart-hack-values.yaml
	yq -i '.descope.project.id = load_str("./hack/descope_project_id")' hack/chart-hack-values.yaml
	yq -i '.descope.managementKey.id = load_str("./hack/descope_management_key_id")' hack/chart-hack-values.yaml
	yq -i '.descope.managementKey.key = load_str("./hack/descope_management_key_token")' hack/chart-hack-values.yaml
	yq -i '.descope.accessKey.id = load_str("./hack/descope_backend_access_key_id")' hack/chart-hack-values.yaml
	yq -i '.descope.accessKey.key = load_str("./hack/descope_backend_access_key_token")' hack/chart-hack-values.yaml
	yq -i '.configMaps[0].data["credential-configuration.json"] = load_str("./hack/gcp-credential-configuration.json")' hack/chart-hack-values.yaml
	yq -i '.configMaps[0].data["credential-configuration.json"] style="literal"' hack/chart-hack-values.yaml
	yq -i '.backend.volumes[0].projected.sources[0].serviceAccountToken.audience = load_str("./hack/gcp-ksa-audience")' hack/chart-hack-values.yaml
	yq -i '.exchangeRatesCronJob.volumes[0].projected.sources[0].serviceAccountToken.audience = load_str("./hack/gcp-ksa-audience")' hack/chart-hack-values.yaml
	yq -i '.initJob.volumes[0].projected.sources[0].serviceAccountToken.audience = load_str("./hack/gcp-ksa-audience")' hack/chart-hack-values.yaml

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

.PHONY: deploy
deploy: rest
	skaffold build -q | skaffold deploy --load-images --build-artifacts -

.PHONY: undeploy
undeploy:
	skaffold delete

.PHONY: dev
dev: rest
	skaffold dev --tail=false --toot=true

.PHONY: teardown-local-cluster
teardown-local-cluster:
	if kind get clusters 2>/dev/null | grep -E '^greenstar$$'; then \
		kubectl --context=kind-greenstar get namespace kube-system --template='{{.metadata.uid}}' | cut -c -32 | xargs printf "%s" > hack/cluster-id; \
		gcloud iam workload-identity-pools providers delete $$(cat hack/cluster-id) --location=global --workload-identity-pool=greenstar --quiet || true; \
		kind delete cluster --name greenstar || true; \
	fi
	rm -fv hack/chart-hack-values.yaml || true
	rm -fv hack/gcp-credential-configuration.json || true
	rm -fv hack/gcp-ksa-audience || true
	rm -fv hack/gcp-ksa-iam-principal || true
	rm -fv hack/gcp-project-number || true
	rm -fv hack/gcp-project-id || true
	rm -fv hack/cluster-jwks.json || true
	rm -fv hack/cluster-id || true

.PHONY: e2e
e2e: deploy
	echo "Testing with Playwright"
