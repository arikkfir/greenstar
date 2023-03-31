# greenstar

GreenSTAR is an accounting application built for the new generation.

- https://greenstar.kfirs.com

(the above do not work yet)

## Setup a local development environment

### Prerequisite tools

```shell
$ brew install jq yq                                # JSON and YAML manipulation
$ brew install go node                              # Programming languages 
$ brew install kubernetes-cli kustomize helm kind   # Kubernetes tools
$ brew install skaffold                             # Kubernetes development tool
```

### Setup

Replace the `NON_SECRET_ID` with something representing you, e.g. your local username, or something similar. The app
(when in dev mode) will create a few transient resources (e.g. Cloud Pub/Sub topics, subscriptions, etc) in the GCP
project associated with this ID, to differentiate it from similar resources used by other developers. Make sure you
don't use something secret - a name or nickname is perfectly fine here.

```shell
$ cat deploy/kind/cluster-config.yaml | yq ".nodes[0].extraMounts[0].hostPath = \"$(cd ~/.config/gcloud && pwd)\"" \
    | kind create cluster --config=- --name=local
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
$ cat > deploy/local/greenstar-secrets.env <<EOF
auth0_domain=AUTH0_DOMAIN
auth0_api_client_id=AUTH0_API_APP_CLIENT_ID
auth0_api_client_secret=AUTH0_API_APP_CLIENT_SECRET
auth0_app_client_id=AUTH0_SPA_APP_CLIENT_ID
EOF
$ cat > deploy/local/local-config.env <<EOF
DEV_MODE_ID=YOUR_OWN_MADEUP_ID
GCP_PROJECT_ID=GCP_PROJECT_TO_USE
REACT_APP_AUTH0_ORG_NAME=DEFAULT_ORG_NAME_IN_AUTH0_TO_USE_IN_LOCALHOST
EOF
```

### Running

```shell
$ skaffold dev --kube-context=kind-local --keep-running-on-failure=true
```

This will build all Docker images, deploy them to the `default` namespace of the `kind` cluster you just set up, and
then make the following URLs available:

- http://localhost - the GreenSTAR web application
- http://localhost/api/ - the GreenSTAR APIs
- localhost:6379 - the Redis server

### Authentication & Authorization

We use Auth0 to manage authentication and authorization. The following objects are used:

- One Auth0 app for the greenstar web application (aka. "SPA")
- One Auth0 app for the greenstar backend API
- Multiple Auth0 APIs (in Auth0, you can define an "API" object)
- Manage a set of Roles
  - Each role defines the set of scopes it grants, for one or more of the APIs above
  - Roles are then granted to users
- User logs into the app, granted some scopes over some of the APIs
  - Upon login, an access token is received by the SPA
- The access token is then used (via the HTTP authorization bearer header) to access the backend API
