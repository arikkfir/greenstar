# greenstar

GreenSTAR is an accounting application built for the new generation.

- https://app.admin.greenstar.kfirs.com
- https://app.greenstar.kfirs.com
- https://app.operations.greenstar.kfirs.com

(the above do not work yet)

## Local development environment

### Prerequisite tools

```shell
$ brew install jq yq                                # JSON and YAML manipulation
$ brew install go node                              # Programming languages 
$ brew install kubernetes-cli kustomize helm kind   # Kubernetes tools
$ brew install skaffold                             # Kubernetes development tool
$ brew install dnsmasq mkcert                       # For local TLS domains
```

### Local testing domains

Given that the application has multiple components, and each needs to be addressable separately via a separate domain,
we need to use the same convention even when running locally.

To do that, we will utilize `dnsmasq` to resolve the `greenstar.test` domain to localhost. To do that, perform the
following actions (once):

```shell
$ cat >> /opt/homebrew/etc/dnsmasq.conf <<EOF
address=/.greenstar.test/127.0.0.1
server=8.8.8.8
server=8.8.4.4
EOF
$ sudo brew services restart dnsmasq
$ dig @127.0.0.1 foobar.greenstar.test # sanity test
$ sudo mkdir /etc/resolver
$ cat <<EOF | sudo tee /etc/resolver/greenstar
domain greenstar.test
search greenstar.test
nameserver 127.0.0.1
EOF
$ sudo killall -HUP mDNSResponder
$ scutil --dns # sanity test; verify resolver for "greenstar.test" appears and working

... replace DNS server with 127.0.0.1 in macOS network manager...

$ dig foobar.greenstar.test # sanity test (notice it's not the same sanity as above)
```

See [this blog post](https://mjpitz.com/blog/2020/10/21/local-ingress-domains-kind/) for more information.

### Local TLS

It is not sufficient to just resolve the `*.greenstar.test` domain to localhost, we also need it to work with TLS, since
we are utilizing Auth0 which uses cryptographic APIs, which require secure hosts. To do that, follow the following
actions:

```shell
$ mkcert -install
$ cd deploy/local
$ mkcert '*.greenstar.test'
```

See [this blog post](https://web.dev/how-to-use-local-https/) for more information.

### Setup cluster

Replace the relevant values with something representing you, e.g. your local username, or something similar. The app
(when in dev mode) will create a few transient resources (e.g. Cloud Pub/Sub topics, subscriptions, etc) in the GCP
project associated with this ID, to differentiate it from similar resources used by other developers. Make sure you
don't use something secret - a name or nickname is perfectly fine here.

```shell
$ kind create cluster --name=greenstar --config=deploy/etc/kind-cluster-config.yaml
```

### Personal secrets & values

Some values needed during local development are not shared, and are personally related to the development team. To make
those available while running, perform the following (assuming you have access to the secrets, or use your own):

```shell
$ touch frontend/apply-patches.sh
$ cat > deploy/local/backend-secrets.env <<EOF
AUTH_DESCOPE_PROJECT_ID=descope_project_id
AUTH_DESCOPE_MANAGEMENT_KEY=descope_management_key
DEV_MODE_ID=your_own_made_up_id
EOF
$ cat > deploy/local/frontend-secrets.env <<EOF
EOF
```

### Running

#### Traefik

Greenstar is relying on Traefik as an ingress controller, and thus needs it to be deployed in the cluster. Since there's
no need to redeploy it on every change, you can deploy it once to your local cluster and forget about it, like so:

```shell
$ skaffold run -f deploy/etc/skaffold-traefik.yaml
```

This will install Traefik to the `default` namespace of the `kind` cluster you just set up. It will start responding to
requests to `*.greenstar.test` though it will not have any backends to route to (see below).

#### Application

To run the application in a dev-loop mode, so that changes you make will get redeployed to the local cluster, simply
run this:

```shell
$ skaffold dev
```

This will build all container images, deploy them to the `default` namespace of the `kind` cluster you just set up, and
then make the following URLs available:

- https://api.greenstar.test - application backend (tenant agnostic)
- https://neo4j.greenstar.test - neo4j console
- https://global.greenstar.test - administration frontend (no tenant)
- [https://*.greenstar.test](https://*.greenstar.test) - tenant frontend (replace `*` with the tenant ID)

#### Application & Traefik together

You can also deploy both Traefik and the application in one dev-loop command, so that they are both redeployed on
changes you make, like so:

```shell
$ skaffold dev -p traefik
```

This activates the `traefik` Skaffold profile, which simply includes the Traefik skaffold module as a dependency.
