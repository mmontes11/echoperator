# 🤖 echoperator 

[![CI](https://github.com/mmontes11/echoperator/actions/workflows/ci.yml/badge.svg)](https://github.com/mmontes11/echoperator/actions/workflows/ci.yml)
[![Release](https://github.com/mmontes11/echoperator/actions/workflows/release.yml/badge.svg)](https://github.com/mmontes11/echoperator/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmontes11/echoperator)](https://goreportcard.com/report/github.com/mmontes11/echoperator)
[![Go Reference](https://pkg.go.dev/badge/github.com/mmontes11/echoperator.svg)](https://pkg.go.dev/github.com/mmontes11/echoperator)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/echoperator)](https://artifacthub.io/packages/search?repo=echoperator)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](https://opensource.org/licenses/MIT)

Simple Kubernetes operator built from scratch with [client-go](https://github.com/kubernetes/client-go).

[Kubernetes operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) implementation using the [client-go](https://github.com/kubernetes/client-go) library. Altough there are a bunch of frameworks for doing this ([kubebuilder](https://book.kubebuilder.io/), [operator framework](https://operatorframework.io/) ...), this example operator uses the tools provided by [client-go](https://github.com/kubernetes/client-go) for simplicity and flexibility reasons. 

[Medium article](https://betterprogramming.pub/building-a-highly-available-kubernetes-operator-using-golang-fe4a44c395c2) that explains how to build this operator step by step.

### Features

- Simple example to understand how a Kubernetes operator works.
- Manages [Echo CRDs](https://github.com/mmontes11/charts/blob/main/charts/echoperator/crds/echo.yml) for executing an `echo` inside a pod.
- Manages [ScheduledEcho CRDs](https://github.com/mmontes11/charts/blob/main/charts/echoperator/crds/scheduledecho.yml) for scheduling the execution of an `echo` inside a pod.
- High Availability operator using Kubernetes [lease](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#lease-v1-coordination-k8s-io) objects.
- Prometheus metrics.
- [Helm chart](https://github.com/mmontes11/charts/tree/main/charts/echoperator).


### Versioning 

|Echo|ScheduledEcho|Job|CronJob|Lease|Kubernetes|
|----|-------------|---|-------|-----|----------|
|v1alpha1|v1alpha1|v1|v1|v1|v1.21.x|

### Installation

```bash
helm repo add mmontes https://charts.mmontes-dev.duckdns.org
```
```bash
helm install echoperator mmontes/echoperator
```

### Custom Resource Definitions (CRDs)

The helm chart installs automatically the [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) needed for this operator to work. However, if you wanted to install them manually, you can find them in the [helm chart repo](https://github.com/mmontes11/charts/tree/main/charts/echoperator/crds).

### Example use cases

###### Hello world

- Client creates a [hello world Echo CRD](./manifests/examples/hello-world.yml).
- Operator receives a `Echo` added event.
- Operator reads the `message` property from the `Echo` and creates a `Job` resource.
- The `Job` resource creates a `Pod` that performs a `echo` command with the `message` property.

###### Scheduled hello world

- Client creates a [hello world ScheduledEcho CRD](./manifests/examples/hello-world-scheduled.yml).
- Operator receives a `ScheduledEcho` added event.
- Operator reads the `message` and `schedule` property from the `ScheduledEcho` and creates a `CronJob`.
- The `CronJob` schedules a `Job` creation using the `schedule` property.
- When scheduled, the `Job` resource creates a `Pod` that performs a `echo` command with the `message` property. 
