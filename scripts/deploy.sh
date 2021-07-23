#!/usr/bin/env bash

set -e

repo="mmontes"
release="echoperator"
chart="$repo/$release"

git fetch --all
tag=$(git describe --tags $(git rev-list --tags --max-count=1))

helm repo add "$repo" https://charts.mmontes-dev.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --set image.tag=$tag
