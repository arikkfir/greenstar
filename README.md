# greenstar

GreenSTAR is an accounting application built for the new generation.

https://greenstar.analytics
https://greenstar.app

(the above do not work yet)

## Setup a local development environment

### Prerequisite tools

```shell
$ brew install jq yq                                # JSON and YAML manipulation
$ brew install go node                              # Programming languages 
$ brew install kubernetes-cli kustomize helm kind   # Kubernetes tools
$ brew install skaffold                             # Kubernetes development tool
```

### Create a `kind` cluster

```shell
$ kind create cluster --config=deploy/kind/cluster-config.yaml --name=local
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
$ cat > google-client-app.env <<EOF
GOOGLE_CLIENT_ID=<your_google_oauth_client_id>
GOOGLE_CLIENT_SECRET=<your_google_oauth_client_secret>
EOF
$ kubectl create secret generic greenstar-google-app \
    --save-config \
    --output=yaml \
    --dry-run=client \
    --from-env-file=google-client-app.env \
    | kubectl apply -f -
```

### Running

```shell
$ skaffold dev --kube-context=kind-local --toot=true --tolerate-failures-until-deadline=true
```

This will build all Docker images, deploy them to the `default` namespace of the `kind` cluster you just set up, and
then make the following URLs available:

- http://localhost - the GreenSTAR web application
- http://localhost/api/ - the GreenSTAR APIs
- http://localhost:9000 - the Redis web client
- localhost:6379 - the Redis server (this port is for a Redis client, not for a browser)
- http://localhost:7474 - the Neo4j web client
- localhost:7687 - the Neo4j server (this port is for a Neo4j client, not for a browser)
