#!/usr/bin/env bash
set -euo pipefail

# Get the most recent commit in the current branch
COMMIT=$(git rev-parse HEAD)
echo "Most recent commit: ${COMMIT}" >&2

# Get the repository name from the remote URL
REPO_NAME="arikkfir/greenstar"

# Use GitHub CLI to find the last run of the verify workflow for that commit
if ! command -v gh &> /dev/null; then
    echo "GitHub CLI (gh) is not installed. Please install it to use this script." >&2
    exit 1
fi

# Check if user is authenticated with GitHub CLI
if ! gh auth status &> /dev/null; then
    echo "You are not authenticated with GitHub CLI. Please run 'gh auth login' first." >&2
    exit 1
fi

echo "Finding the last run of the verify workflow for commit ${COMMIT}..." >&2
RUN_ID=$(gh run list --repo "${REPO_NAME}" --workflow=Verify --commit="${COMMIT}" --json databaseId --jq '.[0].databaseId')

if [ -z "${RUN_ID}" ]; then
    echo "No verify workflow run found for commit ${COMMIT}" >&2
    exit 1
fi

echo "Found workflow run ID: ${RUN_ID}" >&2

# Create a temporary directory for the download
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "${TEMP_DIR}"' EXIT

# Download the cluster-state artifact
echo "Downloading cluster-state artifact from workflow run ${RUN_ID}..." >&2
gh run download "${RUN_ID}" --repo "${REPO_NAME}" --name cluster-state --dir "${TEMP_DIR}"

# Prepare the target directory
TARGET_DIR="hack/cluster-state"
if [ -d "${TARGET_DIR}" ]; then
    echo "Removing existing ${TARGET_DIR} directory..."
    rm -rf "${TARGET_DIR}"
fi

# Create the target directory and move the files
mkdir -p "${TARGET_DIR}"
echo "Extracting artifact to ${TARGET_DIR}..."
cp -vR "${TEMP_DIR}/"* "${TARGET_DIR}/"

echo "Cluster state successfully downloaded and extracted to ${TARGET_DIR}"
