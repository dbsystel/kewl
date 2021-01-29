![Go](https://github.com/dbsystel/kewl/workflows/Go/badge.svg) [![codecov](https://codecov.io/gh/dbsystel/kewl/branch/master/graph/badge.svg?token=E123SJUGFD)](https://codecov.io/gh/dbsystel/kewl) [![Go Reference](https://pkg.go.dev/badge/github.com/dbsystel/kewl/.svg)](https://pkg.go.dev/github.com/dbsystel/kewl/)
# KEWL - K8s Easy Webhook Library

## Description

This library aims to facilitate the implementation of k8s webhooks
for [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
.

## Features

- easy implementation of validators/mutators for k8s objects
- multiple validators and mutators can be added at the same time
- supports v1 and v1beta1 AdmissionReview from the same URLs
- exposes metrics for validators and mutators
- custom handlers for an admission-review can be easily implemented
- validation responses contain the cause of the validation error with the fields and messages
- mutation responses contain an [RFC6902](https://tools.ietf.org/html/rfc6902) compatible JSON patch

## Usage

Add the following line to your ``go.mod`` file and you're all setup:

```
github.com/dbsystel/kewl v1.0.0
```

### Examples

- [Creating a webhook server](examples/server.go)
- [Creating a validator](examples/validator.go) and [Testing it](examples/validator_test.go)
- [Creating a mutator](examples/mutator.go) and [Testing it](examples/mutator_test.go)

### Exposed paths

- `/healthz` for health checks
- `/metrics` for prometheus metrics
- `/validate` for validation hooks
- `/mutate` for mutation hooks

## Metrics and health

### Healthz

The webhook exposes and endpoint `/healthz` which can be used to check, if the server still runs fine.

### Prometheus metrics

Also, prometheus summaries are exposed via `/metrics` for the following:

#### HTTP requests

A prometheus summary is exposed for all requests as `webhook_http_request_seconds_sum` labeled by:

- request `method`
- request `path`
- response `status` code.

Example:

```
webhook_http_request_seconds_sum{method="POST",path="/validate",status="200"} 7.3844e-05
webhook_http_request_seconds_count{method="POST",path="/validate",status="200"} 
```

#### Invoked validations

Invoked validations are registered in a summary named `webhook_handler_validation_sum` labeled by:

- version of the admission review (`admission_review_version`)
- group of the reviewed object: `obj_group`
- kind of the reviewed object: `obj_kind`
- version of the reviewed object: : `obj_version`
- namespace of the reviewed object (`obj_namespace`)
- result of the review (`result`), which can be the following
    - `allowed` - the validation was successful (admission was allowed)
    - `denied` - the validation was unsuccessful (admission was denied)
    - `error` - an error occurred in the server (or validator)

Example:

```
webhook_handler_validation_sum{admission_review_version="v1",group="",kind="Pod",result="allowed",target_namespace="test",version="v1"} 2.9475e-05
webhook_handler_validation_count{admission_review_version="v1",group="",kind="Pod",result="allowed",target_namespace="test",version="v1"} 1
```

#### Invoked mutations

Invoked mutations are registered in a summary named `webhook_handler_mutation_sum` labeled by:

- version of the admission review (`admission_review_version`)
- group of the reviewed object: `obj_group`
- kind of the reviewed object: `obj_kind`
- version of the reviewed object: : `obj_version`
- namespace of the reviewed object (`obj_namespace`)
- result of the review (`result`), which can be the following
    - `allowed` - object was not modified (admission was allowed)
    - `mutated` - object was mutated (admission was allowed)
    - `error` - an error occurred in the server (or mutator)

Example:

```
webhook_handler_mutation_sum{admission_review_version="v1",group="",kind="Pod",result="mutated",target_namespace="test",version="v1"} 4.258e-05
webhook_handler_mutation_count{admission_review_version="v1",group="",kind="Pod",result="mutated",target_namespace="test",version="v1"} 1
```

## License

This project is licensed under Apache License v2.0, which is [included in the repository](./LICENSE.txt).

## Contributions

Contributions are very welcome, please refer to the [Contribution guide](./CONTRIBUTING.md)

## Code of conduct
Our code of conduct can be found [here](./CODE_OF_CONDUCT.md).
