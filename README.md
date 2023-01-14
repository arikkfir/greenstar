# greenstar

GreenSTAR is an accounting application built for the new generation.

https://greenstar.analytics
https://greenstar.app

## Setup a local development environment

### Prerequisite tools

```shell
$ brew install jq yq
$ brew install go node
$ brew install kubernetes-cli kustomize helm kind
$ brew install skaffold
$ brew install postgresql redis neo4j
```

### Point `*.greenstar.test` to localhost (`127.0.0.1`)

Follow this article: https://allanphilipbarku.medium.com/setup-automatic-local-domains-with-dnsmasq-on-macos-ventura-b4cd460d8cb3
Here's the gist of it:

```shell
$ brew install dnsmasq
$ mkdir -pv $(brew --prefix)/etc/
$ echo 'address=/.greenstar.test/127.0.0.1' >> $(brew --prefix)/etc/dnsmasq.conf
$ sudo brew services start dnsmasq
$ dig arctic-jill.test @127.0.0.1 # verification step
$ sudo mkdir -v /etc/resolver
$ sudo bash -c 'echo "nameserver 127.0.0.1" > /etc/resolver/test'
$ ping -c 1 apple.com # verification step
$ ping -c 1 another-sub-domain.test # verification step
$ curl -i --resolve redis.localhost:80:127.0.0.1 http://redis.localhost/ # verification step
$ curl -i --resolve neo4j.localhost:80:127.0.0.1 http://neo4j.localhost/ # verification step
$ curl -i --resolve api.localhost:80:127.0.0.1 http://api.localhost/ # verification step
```

### Create a `kind` cluster

```shell
$ cat <<EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF
$ kind create cluster --config kind-config.yaml
```

### Install NGINX Ingress Controller

```shell
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```
