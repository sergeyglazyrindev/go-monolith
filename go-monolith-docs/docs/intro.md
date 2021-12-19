---
sidebar_position: 1
---

# Tutorial

Let's discover **GoMonolith in less than 30 minutes**.

## Getting Started

Get started by **creating a new project**.

## Hot reloading of your go project on file changes

For this we propose you to use "reflex" go package in a way: `reflex -r '\.go' -s -- sh -c "go run cmd/go-monolith/main.go admin serve"`. And on file changes your project would be reloaded.

## Generate a new project

1. Make sure you use go1.16.
2. Please configure environment variable **GOMONOLITH_PATH** that is a path to the root of your project.
3. Create your project .mod file, as example you can take this file: [gomonolith go.mod example](https://github.com/sergeyglazyrindev/go-monolith/blob/master/examples/go-monolith-example/go.mod)
4. Add file cmd/{{YOUR_PROJECT_NAME}}/main.go with content similar to: [example of how to configure gomonolith and use it's commands](https://github.com/sergeyglazyrindev/go-monolith/blob/master/examples/go-monolith-example/cmd/go-monolith-example/main.go).
5. Add config for your project, for this please create folder: configs in the root of your project. And put there your .yml config file, you can use example from this file: [example of environment](https://github.com/sergeyglazyrindev/go-monolith/blob/master/examples/go-monolith-example/configs/dev.yml)
6. Build your go binary with a command: `go build cmd/go-monolith/main.go`. Your binary will be available in the root of your project with name "main".

## Add new blueprint

1. You can add new blueprint with a name `example` using command: `./main blueprint create -m 'blueprint for gomonolith example' -n example`. Add model if you want. You can check how it's done here: [model file](https://github.com/sergeyglazyrindev/go-monolith/blob/master/examples/go-monolith-example/blueprint/example/models/models.go) and [migration file](https://github.com/sergeyglazyrindev/go-monolith/blob/master/examples/go-monolith-example/blueprint/example/migrations/initial_1631027794.go)
2. Apply all migrations: `./main migrate up`.

## Start your admin panel

1. Please create superuser, so you can sign in into admin panel: `./main superuser create -n adminadmin -e admin@example.com`
2. Add administration panel for the model you created, you can find how to do that in this file: [example of how to add administration panel for your model](https://github.com/sergeyglazyrindev/go-monolith/blob/master/examples/go-monolith-example/blueprint/example/example.go#L16)
3. Start admin panel using command: `./main admin serve`.
4. Open in browser: 127.0.0.1:8080/admin and you are in the admin panel. Sign in using user credentials you created before.
