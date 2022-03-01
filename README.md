# Path Change server

### Summary

A server for CI/CD pipelines to detect file change based on last commit, with [glob pattern](https://github.com/bmatcuk/doublestar#patterns), compatible
with [Tekton ClusterInterceptor](https://tekton.dev/docs/triggers/clusterinterceptors/#configuring-a-kubernetes-service-for-the-clusterinterceptor).

The behaviour will be similar to [github actions paths](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#example-including-paths) and [cloud build includeFile](https://cloud.google.com/build/docs/automating-builds/create-manage-triggers#build_trigger). Only trigger build when the file changes match the patterns

Not all CI server/service support file change detection feature. This project will make this available to CI servers, for example, Tekton.

### Explanation

Assume have the following config

```
repos:
  demo:
    url: https://github.com/username/repo
    paths:
      - "examples/**"
```

#### If last file changes matches the pattern

```
$ git --no-pager diff --name-only HEAD^
examples/main.go
```

It will return the following response, with `continue` set to `true`

```json
{
  "extensions": {
    "paths": "MATCH"
  },
  "continue": true,
  "status": {}
}
```

#### If last file changes not matches

```
$ git --no-pager diff --name-only HEAD^
util/util.go
```

It will return the following response, with `continue` set to `false`

```json
{
  "extensions": {
    "paths": "NOT_MATCH"
  },
  "continue": false,
  "status": {}
}
```

### Installation

#### Docker

#### Kubernetes

#### config

#### git credentials


#### Setting `ref` to be read
Most of the webhook has `ref` on the request body, for example
[github](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#push), [gitlab](https://docs.gitlab.com/ee/user/project/integrations/webhook_events.html#push-events), [gitea](https://docs.gitea.io/en-us/webhooks/#event-information), [gitee](https://gitee.com/help/articles/4186#article-header1)

To customise it, 


#### Install
```
kubectl apply -f examples
tkn hub install task git-clone -n tekton
```
#### Test
```
 curl localhost:8080 -d '{"ref": "refs/heads/main"}' -H 'Content-Type: application/json' -v
```
