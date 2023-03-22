# greenstar

GreenSTAR is an accounting application built for the new generation.

https://greenstar.analytics
https://greenstar.app

## Setup a local development environment

### Prerequisite tools

```shell
$ brew install jq yq                                # JSON and YAML manipulation
$ brew install go node                              # Programming languages 
$ brew install kubernetes-cli kustomize helm kind   # Kubernetes tools
$ brew install skaffold                             # Kubernetes development tool
```

### Point `*.greenstar.test` to localhost (`127.0.0.1`)

The idea here is that we set up a local DNS redirection so every DNS lookup for subdomains of `greenstar.test` (e.g.
`foobar.greenstar.test` would point to `127.0.0.1`). This would enable us to then point different services of the
application like Neo4j, Redis, Postgres, Greenstar API, Greenstar UI, etc to their respective DNS subdomains in a
consistent manner for every developer - and also use these values in external settings like a development Auth0 tenant.

The methodology is details [in this article](https://allanphilipbarku.medium.com/setup-automatic-local-domains-with-dnsmasq-on-macos-ventura-b4cd460d8cb3)
but here's the gist of it to save you time:

```shell
$ brew install dnsmasq
$ mkdir -pv $(brew --prefix)/etc/
$ echo 'address=/.greenstar.test/127.0.0.1' >> $(brew --prefix)/etc/dnsmasq.conf
$ sudo brew services start dnsmasq
$ dig greenstar.test @127.0.0.1 # verification step, expecting to see 127.0.0.1 in dig response
$ dig foobar.greenstar.test @127.0.0.1 # verification step, expecting to see 127.0.0.1 in dig response
$ sudo mkdir -v /etc/resolver
$ sudo bash -c 'echo "nameserver 127.0.0.1" > /etc/resolver/test'
$ ping -c 1 apple.com # verification step, expected to succeed
$ ping -c 1 another-sub-domain.test # verification step, expected to fail (no such domain)
$ ping -c 1 greenstar.test # verification step, expected to succeed
$ ping -c 1 foobar.greenstar.test # verification step, expected to succeed
$ nslookup greenstar.test 127.0.0.1 # 2nd argument is the DNS server to use; verification step, expected to succeed
$ nslookup foobar.greenstar.test 127.0.0.1 # 2nd argument is the DNS server to use; verification step, expected to succeed
```

### Create a `kind` cluster

```shell
$ cat <<EOF | kind create cluster --config=- --name=local
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
$ helm repo add nats https://nats-io.github.io/k8s/helm/charts/
$ helm install my-nats nats/nats --set nats.jetstream.enabled=true
```

### Install NGINX Ingress Controller

```shell
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```
