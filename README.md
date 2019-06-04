# Kubernetes Hello World

## Introduction

This repository contains a small Go web application that will be deployed on a Kubernetes cluster. The application could have a customized coloured page (to test zero-downtime deployments). Further you can control the memory and CPU utilization of the application to simulate load. To control mem/CPU, you can set a mathematical expression (function of `n` iterations of time). The expression could be set using CAS [Expreduce](https://corywalker.github.io/expreduce-docs/). Examples of expressions to be used:

* `60 Sin[n (Pi/20)]`
* Periodic piecewise function: `x = Mod[n, 40]; Which[x <= 20, 60, 20 < x, 0]`
* Periodic Sawtooth function: `4 Mod[n, 20]`

**Note:** The expressions should be eventually a function of `n`. You can use Mathematica expressions. If the expression could't be evaluated or evaluated with negative value, it will be replaced with `0`. If the the evaluated expression is bigger than `100%` for CPU or `max-memory` for memory, then `100%` or `max-memory` will be used respectively.

**Warning:** The CPU utilization algorithm is not precise; it will be enhanced in the future.

Further feature is to set the initialization time which simulates delay until the application is initialized. This could be useful for testing `livenessProbe` and `readinessProbe` of Kubernetes deployments.

## Installation

A kustomize package has been created that have a base the could for test environment, and overlays that fit a real production setup. To install use `kusomoize` or latest `kubectl` (it supports kustomize):

```bash
kubectl apply -k install
```

or

```bash
kusomoize build install | kubectl apply -f -
```

A demo could be seen at https://kubernetes-hello-world.appwavelets.com

## Update a deployment

To update the deployment to green version:

```bash
kubectl apply -k install/prod/green
```

To watch the status of the deployment update:

```bash
kubectl rollout status deployment.v1.apps/kubernetes-hello-world
```

To rollout a deployment to previous versions:

```bash
kubectl rollout undo deployment.v1.apps/kubernetes-hello-world
```

To get the deployment’s rollout history:

```bash
kubectl rollout history deployment.v1.apps/kubernetes-hello-world
```

To rollout to specific version:

```bash
kubectl rollout undo deployment.v1.apps/kubernetes-hello-world --to-revision=1
```

To scale the deployment:

```bash
kubectl scale deployment.v1.apps/kubernetes-hello-world --replicas=5
```

To force update of the current deployment:

```bash
kubectl patch deployment kubernetes-hello-world -p \
  "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"
```

To check the HPA status

```bash
kubectl get hpa kubernetes-hello-world
```

## Configuration

You can configure the following argument:

* `color`: the color of the hello world web page
* `init-duration`: the delay duration to simulate initialization
* `max-memory`: the maximum memory to be allocated
* `memory-expression`: the mathematical expression of memory in MB as a function of n
* `cores`: the value of cores that is allocated to the application. It could be non-integer value for cores in mili. For example if the pod is allocated `2700m` then you should use `2.7` cores for more precise CPU allocation.
* `cpu-expression`: the mathematical expression of CPU percentage as a function of n
* `resolution`: the duration in seconds at which the expressions will be evaluated

## Building application

To build docker image:

```bash
docker build -t abdullahalmariah/kubernetes-hello-world:latest .
docker push abdullahalmariah/kubernetes-hello-world:latest
```

## High availability setup:

To have high available setup, the deployment will have the following feature:

* Replica of 3 instances
* HorizontalPodAutoscaler for memory and CPU utilization
* Pod AntiAffinity for having the pods distributed on several nodes. To enable this feature uncomment the `affinity` section in `install/prod/base/deployment.yaml`
* PodDisruptionPolicy for voluntary disruptions and rolling updates.
* Policy of rollingUpdate for ensuring zero-downtime deployments.
* Deployment with `livenessProbe` and `readinessProbe` for making sure that the rolling update goes without downtime.
