# greenstar

GreenSTAR is an accounting application built for the new generation.

## Local development environment

### Prerequisite tools

```shell
$ brew install jq yq          # Just useful tools for JSON and YAML manipulation
$ brew install go node        # Programming languages
$ brew install kubernetes-cli # Kubernetes CLI
```

### Telepresence setup

Install Telepresence by following [this guide](https://www.getambassador.io/docs/telepresence-oss/latest/install).

It's also possible you'll need to [apply these changes](https://www.getambassador.io/docs/telepresence-oss/latest/troubleshooting#volume-mounts-are-not-working-on-macos)
if you need to have volume mounts working on macOS.

You'll also need to connect to the cluster once with a durable connection, like so:

```shell
$ telepresence connect
```

### Running

The developer story for GreenSTAR is architected on the following principals:

1. Every change is done in a feature or bugfix branch
2. This branch is pushed to GitHub and a pull request is created for it (PR creation has been automated)
3. PR creation triggers a dedicated environment to be created for the PR
    1. This environment is deployed to a dedicated namespace in the cluster
    2. The environment is accessible via a dedicated URL
4. Telepresence is then used to run the backend or frontend locally yet replacing their counterpart in the remote
   namespace. For instance, the backend can be run locally, but traffic arriving to the dedicated environment's backend
   will reach the locally-running process of the backend (thanks to Telepresence).
5. Once the PR is merged, the dedicated environment is deleted

The IntelliJ IDEA project is stored in the repository (the parts that should be shared across engineers) and it contains
run configurations for the backend and frontend. These run configurations are configured to run the application locally.

However, you need to first start a Telepresence intercept, like so:

```shell
# If you want to intercept the backend locally:
$ telepresence intercept --namespace=MY-ENV greenstar-frontend --port 8080:http --env-file ./frontend/intercept.env

# If you want to intercept the frontend locally:
$ telepresence intercept --namespace=MY-ENV greenstar-backend --port 8000:http --env-file ./backend/intercept.env
```
