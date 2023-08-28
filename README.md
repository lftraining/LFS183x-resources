# Zero Trust Labs

## Prerequisites

In order to run these labs, you will need access to a device running a Linux OS. These labs have been tested on Ubuntu 22.04.

You will require the following tools to be installed:

- [Golang](https://go.dev/doc/install)
- [curl](https://everything.curl.dev/get/linux)
- [OpenSSL](https://www.openssl.org/)
- [Docker](https://docs.docker.com/desktop/install/linux-install/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [jq](https://jqlang.github.io/jq/download/)

## Lab 1 - How Asymmetric Encryption Works

[Toy RSA Algorithm](lab-1-toy-rsa/)

## Lab 2 - PKI in Practice

[Simple local PKI using openssl](lab-2-simple-pki/)

## Lab 3 - Mutual TLS

[mTLS between two golang services](lab-3-mtls-golang/)

## Lab 4 - Simple SPIRE Deployment

[Simple SPIRE Deployment](lab-4-simple-spire/)

## Lab 5 - Policy Engines

[Making Policy Decisions Using OPA](lab-5-policy-engines/)

## Lab 6 - Policy Enforcement

[Enforcing OPA policy decisions using Envoy Proxy](lab-6-opa-envoy/)
