# echoperator 🤖

[![CI](https://github.com/mmontes11/echoperator/actions/workflows/ci.yml/badge.svg)](https://github.com/mmontes11/echoperator/actions/workflows/ci.yml)
[![Release](https://github.com/mmontes11/echoperator/actions/workflows/release.yml/badge.svg)](https://github.com/mmontes11/echoperator/actions/workflows/release.yml)
[![Deploy](https://github.com/mmontes11/echoperator/actions/workflows/deploy.yml/badge.svg)](https://github.com/mmontes11/echoperator/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmontes11/echoperator)](https://goreportcard.com/report/github.com/mmontes11/echoperator)
[![Go Reference](https://pkg.go.dev/badge/github.com/mmontes11/echoperator.svg)](https://pkg.go.dev/github.com/mmontes11/echoperator)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/echoperator)](https://artifacthub.io/packages/search?repo=echoperator)

Simple kubernetes operator for handling echo CRDs.

#### Installation 🌱

```bash
helm repo add mmontes https://charts.mmontes-dev.duckdns.org
```
```bash
helm install echoperator mmontes/echoperator
```