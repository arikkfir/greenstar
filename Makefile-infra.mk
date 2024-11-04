hack/gcp-project-id:
	gcloud config get-value project | xargs printf "%s" > hack/gcp-project-id

hack/gcp-project-number: hack/gcp-project-id
	gcloud projects describe $(shell cat hack/gcp-project-id) --format='value(projectNumber)' | xargs printf "%s" > hack/gcp-project-number

hack/gcp-ksa-iam-principal: hack/gcp-project-number
	printf "principal://iam.googleapis.com/projects/$(shell cat hack/gcp-project-number)/locations/global/workloadIdentityPools/greenstar/subject/system:serviceaccount:greenstar:local-greenstar-backend" > hack/gcp-ksa-iam-principal

.PHONY: setup-gcp-workload-identity-pool
setup-gcp-workload-identity-pool: hack/gcp-project-id
	# TODO: this belongs in the project IaC since it's static and identical for all workstations
	gcloud iam workload-identity-pools describe greenstar --location="global" --quiet 2>/dev/null \
		|| gcloud iam workload-identity-pools create greenstar --project=$(shell cat hack/gcp-project-id) --location=global --description=GreenSTAR --display-name=GreenSTAR

.PHONY: grant-ksa-gcp-permissions
grant-ksa-gcp-permissions: hack/gcp-project-id hack/gcp-ksa-iam-principal
	# TODO: this belongs in the project IaC since it's static and identical for all workstations
	gcloud projects remove-iam-policy-binding $(shell cat hack/gcp-project-id) --member=$(cat hack/gcp-ksa-iam-principal) --role=roles/storage.objectViewer --quiet 2>/dev/null || true
	gcloud projects add-iam-policy-binding projects/$(shell cat hack/gcp-project-id) --role=roles/storage.objectViewer --member=$(shell cat hack/gcp-ksa-iam-principal) --condition=None
