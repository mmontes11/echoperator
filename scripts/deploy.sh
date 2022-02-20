#!/usr/bin/env bash

set -e

repo="mmontes"
release="echoperator"
chart="$repo/$release"
namespace="echoperator"

helm repo add "$repo" https://charts.mmontes-dev.duckdns.org
helm repo update

echo "🚀 Deploying '$chart' with image version '$tag'..."
helm upgrade --install "$release" "$chart" --namespace "$namespace"
