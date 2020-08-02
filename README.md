<h1> <img src="https://covergates.com/logo.png" alt="logo" width="48" height=48> Covergates - Portal Gates to Code Quality</h1>

[![badge](https://covergates.com/api/v1/reports/bsi5dvi23akg00a0tgl0/badge?)](https://covergates.com/report/github/covergates/covergates)
![CI](https://github.com/covergates/covergates/workflows/CI/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/covergates/covergates)](https://pkg.go.dev/github.com/covergates/covergates)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![CLA assistant](https://cla-assistant.io/readme/badge/covergates/covergates)](https://cla-assistant.io/covergates/covergates)

## Purpose

Managing coverage reports for private repositories should not be hard!
**Covergates** is to make the easiest way to setup a self-hosted coverage report service.
Use it to improve coding review and quality management flow for your internal projects.
Want to try? Visit [covergates.com](https://covergates.com) before you starting.

## Using

To start service, download [covergates-server.**version**.**os**.zip](https://github.com/covergates/covergates/releases), and try:

```sh
unzip covergates-server.<version>.<os>.zip
./covergates-server
```

To upload report, download [covergates.**version**.**os**.zip](https://github.com/covergates/covergates/releases). Unzip and place it to `/path/to/bin`. Use it for example:

```sh
covergates upload -report <report id> -type go coverage.out
```

## Configure

`covergates-server` uses environment variables to change configurations.
Below is the list of variables for basic configuration:

- `GATES_SERVER_ADDR` Default `http://localhost:8080`
- `GATES_SERVER_BASE` Default `/`
- `GATES_DB_DRIVER` Default `sqlite3`. Other options are `postgres` and `cloudrun`
- `GATES_DB_HOST` Required host for `postgres` and `cloudrun`
- `GATES_DB_PORT` Required port for `postgres` and `cloudrun`
- `GATES_DB_USER` Required user for`postgres` and `cloudrun`
- `GATES_DB_NAME` Required database name for `postgres` and `cloudrun`
- `GATES_DB_PASSWORD` Required password for `postgres` and `cloudrun`
- `GATES_GITEA_SERVER` Default `https://try.gitea.io/`, gitea server address
- `GATES_GITEA_CLIENT_ID` Required for Gitea OAuth login
- `GATES_GITEA_CLIENT_SECRET` Required for Gitea OAuth login
- `GATES_GITHUB_SERVER` Default `https://github.com`
- `GATES_GITHUB_API_SERVER` Default `https://api.github.com`
- `GATES_GITHUB_CLIENT_ID` Required for GitHub OAuth login
- `GATES_GITHUB_CLIENT_SECRET` Required for GitHub OAuth login

## Development

The build is split into `backend`, `cli` and `frontend`. To build backend, run:

```sh
go build -o covergates-server ./cmd/server
```

To build CLI, run:

```sh
export SERVER_API_URL=http://localhost:8080/api/v1
go build -o covergates -ldflags="-X main.CoverGatesAPI=$SERVER_API_URL" ./cmd/cli
```

You may change `SERVER_API_URL` to your self-hosted **covergates-server** address.

If your are behind firewall or proxy,
you may also download source package with `vendor` modules from [covergates.**version**.src.zip
](https://github.com/covergates/covergates/releases). To build with `vendor` modules, run:

```
go build -o covergates-server -mod vendor ./cmd/server
```

To build frontend, it requires:

1. [Node.js v12](https://nodejs.org/en/download/)
2. [togo](https://github.com/bradrydzewski/togo)

Read [web/README.md](https://github.com/covergates/covergates/blob/master/web/README.md) for more details.

## Contributing

It would be highly appreciated if you could contribute to the project.
There are many ways in which you can participate in the project:

1. Contributing directly to the code base

    The expected workflow is [GitHub flow](https://guides.github.com/introduction/flow/).
    Read [CONTRIBUTING.md](https://github.com/covergates/covergates/blob/master/CONTRIBUTING.md) before getting start.

2. [Submit feature requests and bugs](https://github.com/covergates/covergates/issues)

    Especially for the new language support.
    **Covergates** supports `go` and `perl`.
    It would be great if you could provide coverage report examples for other languages.

3. Testing, both unit testing and e2e testing are welcome.

## Milestones

Refer to [TODO.md](https://github.com/covergates/covergates/blob/master/TODO.md) for details.

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](https://github.com/covergates/covergates/blob/master/LICENSE) file for the full license text.

## Screenshots

![report](https://raw.githubusercontent.com/covergates/brand/master/screenshots/covergates.png)


![files](https://raw.githubusercontent.com/covergates/brand/master/screenshots/covergates_code.png)


![setting](https://raw.githubusercontent.com/covergates/brand/master/screenshots/covergates_setting.png)

![pull request](https://raw.githubusercontent.com/covergates/brand/master/screenshots/covergates_pr.png)
