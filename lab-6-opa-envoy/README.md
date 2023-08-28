# OPA-Envoy Lab

## Introduction

In the previous lab, we have seen how OPA can make policy decisions based on structured data. In this lab, we will see how Envoy Proxy can integrate with OPA to act as a Policy Enforcement Point for these decisions. We will use the example setup provided in OPA's [Standalone Envoy Tutorial](https://www.openpolicyagent.org/docs/latest/envoy-tutorial-standalone-envoy/) running in a [kind](https://kind.sigs.k8s.io/) cluster. We encourage you to look through this tutorial before starting this lab, where we will apply the concepts to an application based on the example from the previous lab, involving Alice, Bob, Charlie, Dan and Eve.

## Background

As per the OPA [Standalone Envoy Tutorial](https://www.openpolicyagent.org/docs/latest/envoy-tutorial-standalone-envoy/), we will run an application container in a Kubernetes Pod alongside an OPA container and a standalone envoy proxy container. Envoy External Authorization will be configured as per the OPA tutorial to delegate authorization decisions to OPA. Our application will be a simple API written in Golang which allows users to query employee records as per the scenario set out in the previous [OPA lab](../lab-5-policy-engines). The code is all contained within [main.go](main.go).

Looking through the code, we can see that the following requests are possible:

- `GET` requests to `/api/employees` to fetch all employee records
- `GET` requests to `/api/employees/<employee_id>` to fetch a specific employee's record
- `POST` requests to `/api/employees` to create a new employee record
- `PUT` requests to `/api/employees/<employee_id>` to edit an employee record
- `DELETE` requests to `/api/employees/<employee_id>` to delete a record

Implementing the business logic from the previous [OPA lab](../lab-5-policy-engines) will mean that administrators can perform any of these actions, whereas other users will only be able to send `GET` requests to `/api/employees/<employee_id>` to fetch the record of someone who they manage.

## Lab Setup

Create a kind cluster:

```bash
kind create cluster --image kindest/node:v1.27.3
```

Build a container image for our `employee-records` application, and load this into kind so that we can use the image in our manifests. Note that the Golang code uses a plaintext password for the Postgres database - please do not do this outside of a lab / demo setting!

```bash
docker build -t employee-records:v0.1 .
kind load docker-image employee-records:v0.1
```

Create config maps from our [OPA policy](policy.rego), the [external data](bundle/data.json) from the last lab, and the [Envoy Proxy configuration](config/envoy.yaml) (which is the same as in OPA's [Standalone Envoy Tutorial](https://www.openpolicyagent.org/docs/latest/envoy-tutorial-standalone-envoy/)):

```bash
kubectl create configmap authz-policy --from-file bundle/policy.rego
kubectl create configmap data-bundle --from-file bundle/data.json
kubectl create configmap proxy-config --from-file config/envoy.yaml
```

Have a look through the [OPA policy](bundle/policy.rego) and note how the Envoy-specific format of the `input` passed to OPA requires changes to the Rego that we created in the previous [OPA lab](../lab-5-policy-engines). We will be making decisions on the request `method` (i.e. `GET`, `POST`, `PUT`, etc.), the `path` (e.g. `/api/employees`) and the JWT `bearer`. You can see where these pieces of data sit within OPA's `input` from the `import` statements at the top of the Rego file.

In order for the OPA container to pick up the policy and external data, we will add them to a bundle, and run a bundle server within the cluster for OPA to reach out to:

```bash
kubectl apply -f config/bundle-server.yaml
kubectl wait pods -n default -l app=bundle-server --for condition=Ready --timeout=120s
```

Wait for the wait condition to be satisfied.

Our `employee-records` application needs to persist employee records, so let's spin up a Postgres pod. Note that the manifest uses a plaintext password for the Postgres database - please do not do this outside of a lab / demo environment!

```bash
kubectl apply -f config/postgres.yaml
kubectl wait pods -n default -l app=postgres --for condition=Ready --timeout=120s
```

Wait for the wait condition to be satisfied.

Now we can run our application pod (including application, OPA and Envoy containers):

```bash
kubectl apply -f config/employee-records.yaml
kubectl wait pods -n default -l app=employee-records --for condition=Ready --timeout=120s
```

Wait for the wait condition to be satisfied.

## Policy Enforcement in Action

Export example JWTs for Alice and Bob as environment variables (note that this time we are removing the leading and training double quotes by piping to `sed`, as this will help us later on). Note that you will require [jq](https://jqlang.github.io/jq/download/) to be installed to run these commands.

```bash
export ALICE_JWT=$(docker run -v .:/example openpolicyagent/opa eval -d /example/create_jwt.rego 'data.example.jwt.alice_token' | jq '.result[0].expressions[0].value' | sed -e 's/^"//' -e 's/"$//')
```

```bash
export BOB_JWT=$(docker run -v .:/example openpolicyagent/opa eval -d /example/create_jwt.rego 'data.example.jwt.bob_token' | jq '.result[0].expressions[0].value' | sed -e 's/^"//' -e 's/"$//')
```

We can simulate some requests to the applicaton by spinning up a `curl` container and sending requests to the `employee-records` service from within the cluster. Let's make a request using a JWT obtained by Alice (who is an administrator, as per the `is_admin` claim within the JWT), attempting to retrieve all employee records via a `GET` request to `/api/employees`:

```bash
kubectl run curl --restart=Never -it --rm --image curlimages/curl:8.1.2 \
    -- curl -H "Accept: application/json" \
    -H "Authorization: Bearer $ALICE_JWT" \
    employee-records/api/employees
```

Now let's try the same request but using a JWT obtained by Bob:

```bash
kubectl run curl --restart=Never -it --rm --image curlimages/curl:8.1.2 \
    -- curl -H "Accept: application/json" \
    -H "Authorization: Bearer $BOB_JWT" \
    employee-records/api/employees
```

You should see the employee records when using Alice's JWT, but not with Bob's. However, as per the logic from the previous [OPA lab](../lab-5-policy-engines), Bob is a manager who line manages Charlie and Eve, but not Dan. Let's check that Bob can fetch Charlie's record (whose `employee_id` is 3) via a `GET` request to `/api/employees/3`:

```bash
kubectl run curl --restart=Never -it --rm --image curlimages/curl:8.1.2 \
    -- curl -H "Accept: application/json" \
    -H "Authorization: Bearer $BOB_JWT" \
    employee-records/api/employees/3
```

Let's also check that Bob cannot fetch Dan's record (whose `employee_id` is 4) via a `GET` request to `/api/employees/4`:

```bash
kubectl run curl --restart=Never -it --rm --image curlimages/curl:8.1.2 \
    -- curl -H "Accept: application/json" \
    -H "Authorization: Bearer $BOB_JWT" \
    employee-records/api/employees/4
```

## Optional Challenge Step

Now that we have seen policy decisions and enforcement in action, try to grep out the OPA 'decision logs' from the OPA container running in the `employee-records` pod. (Hint - you can use `kubectl logs...`).

## Teardown

```bash
kind delete cluster
```
