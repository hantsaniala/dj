# dj — Django project CLI helper

`dj` is a Go CLI tool wrapping django-admin, docker, and development workflows.
Replaces ad-hoc shell scripts with structured, configurable commands.

## Install

```bash
go install github.com/hantsaniala/dj@latest
```

Or build from source:

```bash
git clone https://github.com/hantsaniala/dj
cd dj
go build -o dj .
```

## Usage

```
dj [command] [flags] [args...]
```

Run `dj` with no arguments to see the help screen.

## Quick reference

| Command | Description |
|---|---|
| `migrate [args]` | Apply database migrations |
| `mm [args]` | Create migration files |
| `csu` | Create a superuser |
| `manage <args>` | Any manage.py command |
| `serve` | Start dev server |
| `run <cmd> [args]` | Run in project env |
| `dep` | Install Python deps |
| `pull` | Git pull |
| `docker build` | Compose up --build |
| `docker clean` | Compose down |
| `docker shell [svc]` | Exec into container |
| `docker logs [svc]` | Tail container logs |
| `docker migrate` | Migrate inside container |
| `beat` | Celery Beat scheduler |
| `worker [queue]` | Celery worker |
| `secret` | Generate SECRET_KEY |
| `update` | pull+dep+migrate+serve |
| `deploy` | pull+docker build+docker migrate |

## Commands

### Django management

| Command | Description |
|---|---|
| `dj migrate [args...]` | Run database migrations |
| `dj mm [args...]` | Make database migrations |
| `dj csu` | Create a superuser |
| `dj subscribe-all-users` | Subscribe all users to notification topics |
| `dj manage <args...>` | Run any manage.py command (escape hatch) |

### Development

| Command | Description |
|---|---|
| `dj serve` | Start the uvicorn dev server |
| `dj run <command> [args...]` | Run any command inside the project env (`uv run --env-file .env ...`) |
| `dj dep` | Install Python dependencies |
| `dj pull` | Pull latest changes from git |

### Docker

| Command | Description |
|---|---|
| `dj docker build` | Build and start Docker containers |
| `dj docker clean` | Stop and remove Docker containers |
| `dj docker shell [service]` | Open a bash shell in a container |
| `dj docker logs [service]` | Follow container logs |
| `dj docker migrate` | Run migrations inside a container |

### Celery

| Command | Description |
|---|---|
| `dj beat` | Start Celery Beat scheduler |
| `dj worker [queue]` | Start a Celery worker (default queue: `celery`) |

### Utilities

| Command | Description |
|---|---|
| `dj secret` | Generate a random 50-char SECRET_KEY and write it to `.env` |

### Composite workflows

| Command | Description |
|---|---|
| `dj update` | `pull` → `dep` → `migrate` → `serve` |
| `dj deploy` | `pull` → `docker build` → `docker migrate` |

## Configuration

`dj` reads three sources (later overrides earlier):

1. **`config.yaml`** — app defaults
2. **`.env`** — environment variables (loaded automatically)
3. **OS environment** — actual env vars

### config.yaml

```yaml
manager: "uv run --env-file .env manage.py"
serve:
  command: "uv run --env-file .env uvicorn --host 0.0.0.0 --port 8000 core.asgi:application --reload"
docker:
  image: ""
  compose_file: docker-compose.yml
celery:
  app: core
git:
  branch: develop
dep:
  command: uv sync
```

Override any value by setting the equivalent environment variable or `.env` key:

| Config path | Env var |
|---|---|
| `manager` | `MANAGER` |
| `serve.command` | `SERVE_COMMAND` |
| `docker.image` | `DOCKER_IMAGE` |
| `docker.compose_file` | `DOCKER_COMPOSE_FILE` |
| `celery.app` | `CELERY_APP` |
| `git.branch` | `GIT_BRANCH` |
| `dep.command` | `DEP_COMMAND` |

## Project structure

```
.
├── main.go
├── config.yaml
├── cmd/
│   ├── root.go                  # root command, Viper init
│   ├── serve.go                 # dj serve
│   ├── run.go                   # dj run
│   ├── manage.go                # dj manage
│   ├── helpers.go               # shared runManage/runCommand helpers
│   ├── migrate.go               # dj migrate
│   ├── mm.go                    # dj mm
│   ├── csu.go                   # dj csu
│   ├── subscribe_all_users.go   # dj subscribe-all-users
│   ├── dep.go                   # dj dep
│   ├── pull.go                  # dj pull
│   ├── docker.go                # dj docker (parent)
│   ├── build.go                 # dj docker build
│   ├── clean.go                 # dj docker clean
│   ├── shell.go                 # dj docker shell
│   ├── logs.go                  # dj docker logs
│   ├── migrate_docker.go        # dj docker migrate
│   ├── beat.go                  # dj beat
│   ├── worker.go                # dj worker
│   ├── secret.go                # dj secret
│   ├── update.go                # dj update
│   └── deploy.go                # dj deploy
└── internal/
    ├── config/config.go         # Viper setup, .env loader
    └── runner/runner.go         # exec.Cmd helpers
```

## Development

```bash
go build -o dj .
./dj --help
```
