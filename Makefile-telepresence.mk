setup-telepresence:
	telepresence helm install
	telepresence login

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
