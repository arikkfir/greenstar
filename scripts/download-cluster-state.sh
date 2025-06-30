#!/usr/bin/env bash
set -euo pipefail

# Get the most recent commit in the current branch
COMMIT=$(git rev-parse HEAD)
echo "Most recent commit: ${COMMIT}"

# Get the repository name from the remote URL
REPO_URL=$(git config --get remote.origin.url)
REPO_NAME=$(echo "${REPO_URL}" | sed -E 's|.*github\.com[:/]([^/]+/[^/]+)(\.git)?$|\1|')
echo "Repository: ${REPO_NAME}"

# Use GitHub CLI to find the last run of the verify workflow for that commit
if ! command -v gh &> /dev/null; then
    echo "GitHub CLI (gh) is not installed. Please install it to use this script."
    exit 1
fi

# Check if user is authenticated with GitHub CLI
if ! gh auth status &> /dev/null; then
    echo "You are not authenticated with GitHub CLI. Please run 'gh auth login' first."
    exit 1
fi

echo "Finding the last run of the verify workflow for commit ${COMMIT}..."
RUN_ID=$(gh run list --repo "${REPO_NAME}" --workflow=verify.yaml --commit="${COMMIT}" --json databaseId --jq '.[0].databaseId')

if [ -z "${RUN_ID}" ]; then
    echo "No verify workflow run found for commit ${COMMIT}"
    exit 1
fi

echo "Found workflow run ID: ${RUN_ID}"

# Create a temporary directory for the download
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "${TEMP_DIR}"' EXIT

# Download the cluster-state artifact
echo "Downloading cluster-state artifact from workflow run ${RUN_ID}..."
gh run download "${RUN_ID}" --repo "${REPO_NAME}" --name cluster-state --dir "${TEMP_DIR}"

if [ ! -d "${TEMP_DIR}/cluster-state" ]; then
    echo "Failed to download cluster-state artifact or artifact does not contain expected directory structure"
    exit 1
fi

# Prepare the target directory
TARGET_DIR="hack/cluster-state"
if [ -d "${TARGET_DIR}" ]; then
    echo "Removing existing ${TARGET_DIR} directory..."
    rm -rf "${TARGET_DIR}"
fi

# Create the target directory and move the files
mkdir -p "${TARGET_DIR}"
echo "Extracting artifact to ${TARGET_DIR}..."
cp -R "${TEMP_DIR}/cluster-state/"* "${TARGET_DIR}/"

echo "Cluster state successfully downloaded and extracted to ${TARGET_DIR}"
