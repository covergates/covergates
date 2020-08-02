# Contribution Guidelines

Thanks for your interest in contributing **covergates**.
Before you start, please read this document.
It will help you set yourself up with a good development environment.

## Signing the CLA

First, please read and sign the
[Electronic Covergates Contribution License Agreement](https://cla-assistant.io/covergates/covergatesx).
By signing a contributor license agreement, it is ensured that the community is free to use your contribution.

## Discuss your Design

Please let everyone knows what you are working on by [creating a issue](https://github.com/covergates/covergates/issues).
This process prevents duplication of effort, and ensure that the idea is discussed and fits the project goal.

## Build and Run

To get starting, it assumes a certain level of comfort with `Go`,
`Vue`, `typescript` and `Docker`. To get some helps, you could refer to below links:

1. [A Tour of Go](https://tour.golang.org/welcome/1)
2. [What is Vue.js](https://vuejs.org/v2/guide/)
3. [TypeScript in 5 mins](https://www.typescriptlang.org/docs/handbook/typescript-in-5-minutes.html)
4. [docker for beginners](https://docker-curriculum.com/)

The project also uses [google/wire](https://github.com/google/wire) for dependency injection.
You may refer to [Wire Tutorial](https://github.com/google/wire/blob/master/_tutorial/README.md) for more detail.

### Getting the Source Code

Fork the repository so that you can make a pull request. Then clone your fork by:

```sh
git clone https://github.com/<<<your-github-account>>>/covergates.git
cd covergates
```

Create a new branch for your change.
Manage any merge conflicts, commit them, and then push them to your fork.
It is suggested to have one feature or bug-fix per branch.

### Setup SCM (source control management)

**Covergates** requires to connect with SCM to work.
Although you could connect it with GiHub,
it is suggest to you setup a self-hosted SCM for testing.
It's recommended to choose [Gitea](https://gitea.io/en-us/). You can refer to
[Installation with Docker](https://docs.gitea.io/en-us/install-with-docker/) for the detail.
You can also setup Gitea by running our container for testing:

```sh
docker run -d -p 3000:3000 blueworrybear/gitea-docker
```

Visit [http://localhost:3000](http://localhost:3000) to see if it works.  The admin user and password is `gitea` and `gitea`.

### Run

To test your changes, run:

```sh
# If using blueworrybear/gitea-docker container
export GATES_GITEA_CLIENT_ID="c8c6a2cc-f948-475c-8663-f420c8fc15ab"
export GATES_GITEA_CLIENT_SECRET="J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE="
export GATES_GITEA_SERVER="http://localhost:3000"
export GATES_GITEA_SKIP_VERITY="true"
go run ./cmd/server
```

Visit [http://localhost:8080](http://localhost:8080) for the `covergates` server.

### Debug

If yor are using `VS code` as your ID, you could setup the `launch.json` as blow and press `F5` to debug:

```yml
{
 "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "env": {
        "DEBUG": "true",
        "DEBUG_LOGIN": "gitea",
        "DEBUG_EMAIL": "gitea@gmail.com",
        "DEBUG_AVATAR": "http://localhost:3000/user/avatar/gitea/-1",
        "DEBUG_GITEA_TOKEN": "1749a6106454f05f689051c331680c13d78d81b7",
        "DEBUG_GITEA_LOGIN": "gitea",
        "GATES_GITEA_CLIENT_ID": "c8c6a2cc-f948-475c-8663-f420c8fc15ab",
        "GATES_GITEA_CLIENT_SECRET": "J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE=",
        "GATES_GITEA_SERVER": "http://localhost:3000",
        "GATES_GITEA_SKIP_VERITY": "true"
      },
      "buildFlags": "-tags='debug'"
    }
  ]
}
```

### Working with Frontend

If you would like to change frontend, you may need to change backend server to run on different port by:

```sh
export DEBUG="true"
export DEBUG_LOGIN="gitea"
export DEBUG_EMAIL="gitea@gmail.com"
export DEBUG_AVATAR="http://localhost:3000/user/avatar/gitea/-1"
export DEBUG_GITEA_TOKEN="1749a6106454f05f689051c331680c13d78d81b7"
export DEBUG_GITEA_LOGIN="gitea"
export GATES_GITEA_CLIENT_ID="c8c6a2cc-f948-475c-8663-f420c8fc15ab"
export GATES_GITEA_CLIENT_SECRET="J8YYirhYOZY9a9RepaoORN-8EFcSO-sbwjSGvGo4NwE="
export GATES_GITEA_SERVER="http://localhost:5900"
export GATES_GITEA_SKIP_VERITY="true"
go run ./cmd/server
```

> The **DEBUG** environment variables is required to keep you login to covergates.


And then run:

```sh
cd web
npm install
npm run server
```

It will start a server at [http://localhost:8080](http://localhost:8080).
You could also modify backend proxy by editing `VUE_APP_PROXY` at `web/.env.development`.

### Testing

To test backend, simply run:

```sh
go test ./...
```

or testing with Gitea:

```
go test -tags="gitea" ./...
```

To test frontend, run:

```
cd web
npm run test:unit
```

### Build

To build backend, run:

```sh
go build -o covergates-server ./cmd/server
```

To build CLI, run:

```sh
go build -o covergates ./cmd/cli
```

To build frontend, you require installing `togo` by running below command outside go module folder:

```sh
go get github.com/bradrydzewski/togo
```

And then run:

```sh
go generate ./web
```

## Convention

### Linting is Critical

Please make sure your code pass both `golint` and `eslint` by run:

```sh
golint ./...
cd web
npm run lint
```

## Pull Request

You have to sign your commit before we can accept a pull request from you.
Please refer to [Signing commits
](https://docs.github.com/en/github/authenticating-to-github/signing-commits) for the detail.
You could also refer to [GitHub flow](https://guides.github.com/introduction/flow/) to learn
how to properly create a pull request.

## Be Responding to Feedback

The project owners may recommend adjustment to your code. Don's get discourage!
Any suggestion means that we care about your work! We just hope your work meets
the project goal and could shine over decades.

*At last, thank you for contributing to the covergates open source project!*
