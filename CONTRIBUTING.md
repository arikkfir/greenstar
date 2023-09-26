# Contributing

Hi there! We're thrilled that you'd like to contribute to this project. Your help is essential for keeping it great.

Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in
this project you agree to abide by its terms.

## Setup

### TL;DR:

```shell
$ make setup-local-dns
$ make setup-local-cluster
```

Optional:

```shell
$ echo "export JAVA_TOOL_OPTIONS=\"-Djavax.net.ssl.trustStoreType=KeychainStore\"" > ~/.oh-my-zsh/custom/greenstar-openjdk.zsh
$ brew tap raggi/ale && brew install openssl-osx-ca && brew services start openssl-osx-ca
```

To delete the cluster:

```shell
$ make teardown-local-cluster
```

> [!NOTE]
> Currently, macOS is the only supported development platform. We expect Linux support to come soon.

> [!IMPORTANT]
> The `setup-local-dns` target invokes `sudo` to set up the local DNS resolving - therefore, you might get some password
> prompts during invocation. The commands that use `sudo` can be found in the `Makefile` for review.
> The reason for using `sudo` is that some DNS changes require elevated permissions, such as:
> - `brew services restart dnsmasq` to restart `dnsmasq` after configuring it
> - Placing a file in `/etc/resolver/greenstar.local` to direct macOS to use `dnsmasq` when resolving `greenstar.local`

### What does `setup-local-dns` do?

The `setup-local-dns` will enable resolving the `*.greenstar.local` DNS entries in your workstation. You can think of it
as an alias to `localhost`, with the only difference that you can generate TLS certificates for it. In fact, the setup
does exactly that as part of the `setup-local-cluster` and installs these certificates in the local cluster's Ingress
controller (Traefik).

This way, you will be able to use https://acme.app.greenstar.local (note the use of `HTTPS` protocol!) without getting
any prompts from the browser about bad or missing certificates. In other words, full TLS development right here in your
workstation!

#### Self-signed certificates & macOS

Some tools do not use macOS's native key-chain of trusted certificates, and thus fail when trying to connect via TLS to
the `greenstar.local` domain. This can usually be fixed for common platforms - here are a few:

##### OpenJDK

To make OpenJDK use the native macOS key-chain trust store:

```shell
export JAVA_TOOL_OPTIONS="-Djavax.net.ssl.trustStoreType=KeychainStore"
```

This, of course, will only apply to the shell session you run it in. You can add this to your `zsh` or `bash` startup
scripts so it gets executed for every session.

See details
here: https://stackoverflow.com/questions/14280578/how-to-set-up-java-vm-to-use-the-root-certificates-truststore-handled-by-mac-o

##### OpenSSL

The Homebrew-installed version of OpenSSL **does not** use macOS's Keychain Certificates, which means that it will not
consider the certificates returned as valid.

To make it use the native macOS key-chain trust store:

```shell
$ brew tap raggi/ale
$ brew install openssl-osx-ca
$ brew services start openssl-osx-ca
```

See details here: https://akrabat.com/syncing-macos-keychain-certificates-with-homebrews-openssl/

### What does `setup-local-cluster` do?

This target sets up a local Kubernetes cluster using `kind`, and installs the following components in it:

- Kubernetes Gateway API CRDs
- Traefik Ingress Controller & Gateway API implementation
- Metrics Server
- Open Telemetry DaemonSet and Deployment
- Prometheus server
- Certificate Manager (`cert-manager`)
- Jaeger server
- Grafana

Once this target finishes, you can use your familiar `kubectl` to interact with the cluster, as well as deploy the
application into it (the Makefile has targets for that, see below).

## Developing

Once a local `kind` cluster has been created, the easiest way to get up & running is using Skaffold to (continuously)
deploy the Helm chart into the cluster. Skaffold will keep the deployed chart updated with any code changes you make
automatically, and can be run from the CLI, or from the VSCode/JetBrains IDEs (easier).

To run it from the CLI, simply run:

```shell
$ make dev
```

This will keep running in the shell foreground, and redeploy whenever you make changes to your source code. To stop,
just press CTRL+C and it will undeploy and stop.

If, however, you prefer to perform a one-off deployment without the Skaffold development loop, use this instead:

```shell
$ make deploy
```

To remove it, use this:

```shell
$ make undeploy
```

## Issues and PRs

If you have suggestions for how this project could be improved, or want to report a bug, open an issue! We'd love all
and any contributions. If you have questions, too, we'd love to hear them.

We'd also love PRs. If you're thinking of a large PR, we advise opening up an issue first to talk about it, though! Look
at the links below if you're not sure how to open a PR.

## Submitting a pull request

1. [Fork](https://github.com/arikkfir/devbot/fork) and clone the repository.
2. Make sure you set up your local development environment as described above.
3. Create a new branch: `git checkout -b my-branch-name`.
4. Make your change, add your feature/bug specific tests, and make sure the entire tests suite passes (see above)
5. Push to your fork, and submit a pull request
6. Pat your self on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Write and update tests.
- Keep your changes as focused as possible
    - Smaller PRs make faster & easier reviews, which make faster acceptance & merges
    - Keep each PR focused on one specific change
- Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
- Provide any and all necessary information in the PR description to help reviewers understand the context and impact of
  your change.

Work in Progress pull requests are also welcome to get feedback early on, or if there is something blocked you.

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)
