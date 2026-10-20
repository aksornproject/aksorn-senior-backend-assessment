# Aksorn Senior Backend Engineer assessment

## Assignment overview

Build part of a fictional course-progress service in Go. The exercise assesses business-rule correctness, testing, and engineering judgment. Your work will not be incorporated into Aksorn products.

You have **four hours** to read the detailed specification, implement the required behavior, write tests, and prepare a short technical note. Verify setup before starting the timer. Stop at four hours and identify unfinished work. Optional extensions are not necessary to pass.

Use this README to prepare and submit your work. Read [ASSIGNMENT.md](ASSIGNMENT.md) for the detailed rules after setup, once you start the timer.

## What you need to implement

Complete the two TODOs in `internal/course/service.go`:

1. Complete `Service.RecordAttempt` to store accepted attempts with the supplied repository.
2. Complete `Service.Progress` to calculate lesson results and course progress according to [ASSIGNMENT.md](ASSIGNMENT.md#required-implementation).

Add focused tests beside `service_test.go` and a `TECHNICAL_NOTE.md` of at most half a page. The [submission instructions](#submission-instructions) list the note requirements.

## What is already provided

The starter includes the following:

- The HTTP server, Chi routes, JSON parsing, validation, and response encoding.
- The fixed course, input and output types, and PostgreSQL attempt repository.
- PostgreSQL connections through `pgxpool`, Redis connectivity, and connection cleanup.
- Docker Compose and versioned SQL migrations applied by golang-migrate.
- An example application test and a setup check.

Within `internal/course`, `model.go` defines the types, `catalog.go` holds the course and validation, `repository.go` stores attempts, and `http.go` handles transport. `cmd/server` starts the application. `cmd/checksetup` verifies setup.

The example test checks an existing validation rule. It does not test the unfinished business behavior. Valid attempt and progress requests return **501 Not Implemented** until you complete the TODOs.

## Technical constraints and implementation freedom

Keep Chi for routing and `pgxpool` for PostgreSQL connections. Use the supplied plain SQL, or choose SQLC, SQLX, or Bun. Bun is the only permitted ORM. GORM, Ent, and other ORMs are not allowed. Changing the query layer is optional and earns no extra credit.

You can complete the baseline without restructuring the starter. You may add small helpers or change implementation details. Preserve the documented HTTP contract and setup commands. Keep startup under ten seconds for the supplied tooling.

The required implementation assumes sequential requests, fixed configuration, and trusted internal callers. Authentication, prerequisites, idempotency, concurrent-request safety, a frontend, and additional infrastructure are outside that scope. Use distinct attempt IDs in baseline examples and tests. Retry semantics are unspecified unless you choose the optional extension.

PostgreSQL remains the source of truth. Redis is connected and checked at startup, but no caching, distributed locking, or other Redis behavior is required. See [ASSIGNMENT.md](ASSIGNMENT.md#optional-extensions) for the unchanged optional extensions.

## Prerequisites and setup

Install a supported stable release of **Go 1.26 or newer**, using the latest patch from [Go downloads](https://go.dev/dl/). Install Docker with Docker Compose. Make is optional. You also need Git and a personal GitHub account for submission.

Fork this repository to your personal GitHub account and clone your fork before starting work. If GitHub prevents you from forking or accessing the repository, contact the hiring team before starting the timer.

From your local checkout, run:

```sh
make setup
```

Without Make, run the equivalent commands:

```sh
docker compose up -d --wait postgres redis
docker compose run --rm migrate up
go test ./...
go vet ./...
go build -o bin/server ./cmd/server
go run ./cmd/checksetup
```

Setup downloads the dependencies pinned in `go.mod`, starts PostgreSQL and Redis, applies migrations, and checks the application. No local migration CLI installation is needed. The defaults work without an environment file.

The example test and setup check pass on the untouched starter. The setup check accepts the expected 501 responses. It verifies setup, not business correctness. If setup fails, contact the hiring team before starting the timer. Infrastructure setup is outside the four-hour timebox.

## Running the application

Start the server:

```sh
go run ./cmd/server
```

The server listens on `127.0.0.1:8080`. To change the port, use `go run ./cmd/server -addr 127.0.0.1:8081`. Stop the server with Ctrl+C. Attempts persist in PostgreSQL across server and container restarts.

In another terminal, try these requests:

```sh
curl -i http://127.0.0.1:8080/healthz
curl -i -X POST http://127.0.0.1:8080/attempts \
	-H 'Content-Type: application/json' \
	--data '{"attempt_id":"attempt-001","learner_id":"learner-001","course_id":"course-001","lesson_id":"introduction","score":60}'
curl -i http://127.0.0.1:8080/learners/learner-001/courses/course-001/progress
```

Health returns 200. The two business endpoints initially return 501 with a JSON explanation. Once implemented, the POST returns 201 with the accepted attempt, and the GET returns 200 with progress. See [examples/requests.http](examples/requests.http) for additional requests.

## Testing

Run the tests:

```sh
go test ./...
```

Write focused tests for the required behavior and boundaries in [ASSIGNMENT.md](ASSIGNMENT.md#automated-tests). Application tests are sufficient. HTTP tests are optional, and there is no coverage-percentage requirement.

For application tests that use the supplied repository, open a pool with the supplied configuration. Use unique learner and attempt IDs to keep test runs independent. Unit tests may use a small repository abstraction if useful. No particular testing architecture is required.

## Submission instructions

Follow this workflow:

1. Fork the assessment repository to your personal GitHub account.
2. Develop the assignment on your fork.
3. Commit changes progressively as you work.
4. Push your work to GitHub.
5. Ensure that the hiring team can access your repository, either publicly or through granted private access.
6. Use **Reply All** on the recruitment email and include your repository link.

Stop after four hours. Include your code, tests, and `TECHNICAL_NOTE.md`. Exclude binaries and local caches.

Keep `TECHNICAL_NOTE.md` to at most half a page. Include one design decision and its tradeoff, limitations or unfinished work, one concrete correctness risk involving concurrent or retried requests, and what you would verify before production use. State your time spent and any AI tools used. AI tools are permitted. You must understand and be able to explain and change your submitted code.

You may include additional notes in your reply about trade-offs, limitations, unfinished work, or special setup instructions. These notes do not replace `TECHNICAL_NOTE.md`.

## Additional technical reference

### Connections and query tools

`compose.yaml` supplies PostgreSQL on `127.0.0.1:15432` and Redis on `127.0.0.1:16379`. The PostgreSQL username, password, and database name are each `assessment`. These credentials are for this exercise only. `internal/infra/connections.go` supplies pool configuration, connection checks, and cleanup.

To use different local ports or infrastructure, override `DATABASE_URL` or `REDIS_ADDR`. Update the Compose port mappings to match if needed. Migrations run inside Docker and use `postgres:5432`. Set `MIGRATION_DATABASE_URL` to a Docker-reachable URL only when targeting another exercise database.

`/healthz` reports process health. Startup checks the dependencies and the supplied attempts table.

SQLC can target pgx/v5. SQLX or Bun can use a `database/sql` adapter from [`stdlib.OpenDBFromPool`](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib#OpenDBFromPool) while retaining the supplied pool. Close both the adapter and the underlying pool on shutdown.

### Migrations and local data

`db/migrations` contains versioned up and down SQL files, applied by the pinned [golang-migrate](https://github.com/golang-migrate/migrate) Docker image. Migrations are supplied setup work, not a candidate TODO.

Use these commands to manage migrations:

- `make migrate-up` applies pending migrations. If the database is current, it reports `no change`.
- `make migrate-version` shows the current version.
- `make migrate-down` rolls back one migration. **Rolling back version 1 deletes the attempts table and its data.**

`make infra-down` stops this exercise's containers and preserves PostgreSQL data. Redis data is disposable.

**The following reset deletes this exercise's database.** To erase its data and reapply migrations, run `docker compose down -v` followed by `make setup`.

If you used the earlier starter with `db/init.sql`, that volume has no migration history. For disposable exercise data, use the reset command above before using this version. Do not force a migration version onto a database with an unverified schema or data you need to keep.
# aksorn-senior-backend-assessment
