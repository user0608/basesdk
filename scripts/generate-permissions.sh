#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(dirname "$script_dir")

cd "$repo_root"

go run ./cmd/permissiongen \
	-input permissions \
	-output security/permissions/permissions.go \
	-ts-output ui/package/src/generated/permissions.ts \
	-package permissions
