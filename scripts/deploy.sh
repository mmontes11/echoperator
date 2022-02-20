#!/usr/bin/env bash

set -e

repo="mmontes"
release="echoperator"
chart="$repo/$release"

git fetch --all
tag=$(git describe --abbrev=0 --tags)

helm repo add "$repo" https://charts.mmontes-dev.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --set image.tag=$tag --namespace echoperator
