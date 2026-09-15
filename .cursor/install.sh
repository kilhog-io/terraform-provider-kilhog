#!/usr/bin/env bash
# Copyright IBM Corp. 2021, 2026
# SPDX-License-Identifier: MPL-2.0

# Cloud Agent install script for terraform-provider-kilhog.
#
# Prepares a full development environment: Go module cache, golangci-lint,
# the Terraform CLI, and a local build of the Kilhog API used by the
# acceptance tests. Safe to run repeatedly (idempotent).
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

GOBIN_DIR="$(go env GOPATH)/bin"
export PATH="${GOBIN_DIR}:${PATH}"

# Go toolchain version this module targets (e.g. 1.26.5). Tools are built with
# a matching toolchain so their build version is not older than go.mod's target.
GO_VERSION="$(awk '/^go / {print $2; exit}' go.mod)"

# Pinned tool versions.
GOLANGCI_VERSION="v2.6.1"
TERRAFORM_VERSION="1.14.0"

# Local checkout of the Kilhog API (dependency for acceptance tests).
KILHOG_API_DIR="${HOME}/kilhog-api"
KILHOG_API_REPO="https://github.com/kilhog-io/kilhog.git"

echo "==> Downloading Go modules"
go mod download
(cd tools && go mod download)

echo "==> Installing golangci-lint ${GOLANGCI_VERSION} (built with go${GO_VERSION})"
if ! /usr/local/bin/golangci-lint version 2>/dev/null | grep -q "${GOLANGCI_VERSION#v}"; then
  # Prebuilt golangci-lint binaries lag the Go release used here, and it refuses
  # to run when built with a Go older than the target module. Build from source
  # with the module's Go toolchain to guarantee a compatible build version.
  GOTOOLCHAIN="go${GO_VERSION}" go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_VERSION}"
  sudo install -m 0755 "${GOBIN_DIR}/golangci-lint" /usr/local/bin/golangci-lint
fi

echo "==> Installing Terraform ${TERRAFORM_VERSION}"
if ! /usr/local/bin/terraform version 2>/dev/null | grep -q "v${TERRAFORM_VERSION}"; then
  tmpdir="$(mktemp -d)"
  curl -sSfLo "${tmpdir}/terraform.zip" \
    "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip"
  unzip -o -q "${tmpdir}/terraform.zip" -d "${tmpdir}"
  sudo install -m 0755 "${tmpdir}/terraform" /usr/local/bin/terraform
  rm -rf "${tmpdir}"
fi

echo "==> Building the provider"
go build -v ./...

echo "==> Preparing the Kilhog API at ${KILHOG_API_DIR}"
if [ ! -d "${KILHOG_API_DIR}/.git" ]; then
  git clone --depth 1 "${KILHOG_API_REPO}" "${KILHOG_API_DIR}"
else
  git -C "${KILHOG_API_DIR}" pull --ff-only || true
fi
(cd "${KILHOG_API_DIR}" && make build)

# Expose the local API endpoint and dev key to interactive/login shells so
# `make testacc` works out of the box. The Kilhog API rejects writes unless an
# API key is configured, so acceptance tests require a matching key. This value
# is a throwaway local development constant, not a secret.
echo "==> Configuring shell environment for acceptance tests"
MARKER_BEGIN="# >>> terraform-provider-kilhog dev env >>>"
MARKER_END="# <<< terraform-provider-kilhog dev env <<<"
if ! grep -qF "${MARKER_BEGIN}" "${HOME}/.bashrc" 2>/dev/null; then
  {
    echo ""
    echo "${MARKER_BEGIN}"
    echo "export KILHOG_BASE_URL=http://127.0.0.1:8080"
    echo "export KILHOG_API_KEY=ci-secret"
    echo "${MARKER_END}"
  } >> "${HOME}/.bashrc"
fi

echo "==> Install complete"
