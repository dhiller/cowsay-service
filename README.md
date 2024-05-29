# cowsay-service

## Goal

Create a demo environment that consists of
1. an "external" VM that hosts a service and
2. a local k8s cluster that consumes the service on the VM

# About this repository

This repository holds all setup to show a demo importing a virtual machine into a Kubernetes cluster where [KubeVirt](kubevirt.io) is installed, and then exercising some of the possibilities of KubeVirt.

It contains the setup for two services

* the [cowsay service](./cmd/cowsay-service)
* the [fortune service](./cmd/fortune-service)

## Requirements
* a linux system acting as the demo host
* [golang](https://go.dev/doc/install)
* [go-tasks](https://taskfile.dev/installation/)
* [libvirt](https://libvirt.org/compiling.html#installing-from-distribution-repositories)
* [kubectl](https://kubernetes.io/docs/reference/kubectl/)

## Preparation

1. Configure a virtual machine as described in [./vm/README.md](./vm/README.md)
2. `source envrc`
3. `go-task vm:prepare`
4. `go-task demo:prepare`

This should give you the initial demo environment, where you will have a vm running, and a local kubernetes cluster consuming the external service.
It should launch the fortune service that gives you a fortune embedded inside a random cowsay character.

## Demo

### Step 0: call the fortune service

Test it:
```bash
watch -n 4 curl -s localhost:9090
```

You may run this command in a separate terminal to see what happens during changes to the services and vms.

### Step 1: import the VM image into the cluster

This step imports the prepared vm into the Kubernetes cluster, starts it and redirects the fortune service to use the internal service backed by the vm.

```bash
# start vm import, change the service to using the internal vm
go-task demo:fortune-internal

# show virtual machines in the cluster
kubectl get vmis
```

### Step x: TODO
