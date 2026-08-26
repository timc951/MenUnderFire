# Men Under Fire

A group accountability platform. Members file periodic check-in reports to a small
private group; reports can be anonymous to the group while remaining attributable
to the group's leader. Built as a full-stack application with a Go API and a React
front end.

---

## Why it looks the way it does

A few design decisions carried most of the weight:

**Anonymity is a read-time property, not a storage property.** Reports always store
their author. `ReportService.List` takes a `requesterID` and decides per-caller
whether the author is revealed. Storing reports anonymously would have made
moderation and leader follow-up impossible; deciding at read time keeps one record
and one code path.

**Roles are hierarchical, and group authorization lives in the service layer.**
`OWNER > LEADER > MODERATOR > MEMBER`, plus organization admins and site admins.
A member can only act on someone strictly below them — so role changes and removals
share one comparison rather than a per-endpoint permission list.

This is not uniform, and I would rather name it than let you find it: site-admin
gating for admin-only endpoints (`site_page_handler`, `organization_handler`,
`dashboard_handler`) is still done in the handler. Group-scoped decisions are in the
services; the site-admin check never moved down. It is the first thing I would
refactor.

**Interfaces at the boundaries, so the services are testable without a database.**
Each service depends on repository interfaces, not on `*sql.DB`. That is what lets
the service layer sit at ~73% statement coverage with no test containers and a
sub-4-second suite.

## Architecture

```
cmd/api            entrypoint, dependency wiring
internal/
  routes           mux routing, CORS, public/protected split
  middleware       JWT auth (JWKS), rate limiting
  handlers         HTTP concerns only: decode, call service, encode
  services         business logic and authorization
  repositories     interfaces
    postgres         SQL implementations
  models           request/response and domain types
  database         connection and migrations
```

Dependencies point inward. Handlers know about services; services know about
repository interfaces; only the `postgres` package knows SQL.

## Stack

| Layer | Choice |
|---|---|
| API | Go 1.26, `gorilla/mux`, `zerolog` |
| Database | PostgreSQL, `golang-migrate` (21 migrations) |
| Auth | OIDC / JWT with JWKS verification and key caching |
| Front end | React 19, TypeScript, Vite 8, Tailwind |
| Testing | Go stdlib `testing`; Vitest + Testing Library + MSW |
| Deploy | Docker Compose, Cloudflare Tunnel, Traefik |
| CI | GitHub Actions and GitLab CI |

## Running it

Requires Docker, Go 1.26+, and a current Node LTS.

```bash
cp .env.example .env      # then fill in the AUTH_* values
docker compose up -d      # API, Postgres, proxy
```

Migrations run automatically on API start.

Backend directly:

```bash
cd backend_go
go test ./...             # service layer ~73% coverage
go build ./cmd/api
```

Front end:

```bash
cd frontend
npm ci
npm test                  # 240 tests
npm run dev
```

## Testing and analysis

The service layer carries the tests, because that is where the authorization
decisions live. Repositories are covered indirectly; handlers are thin enough to
be low-risk.

```bash
cd backend_go
go vet ./...
golangci-lint run ./...
gosec ./...
govulncheck ./...
```

All four are clean at the current commit.

## Notes on authorship

This project was built with substantial AI assistance. I wrote the specifications,
made the architecture and data-model decisions, and reviewed, corrected, and tested
the generated code; a large share of the initial implementation was produced by an
LLM working from those specs. The security hardening, the RBAC hierarchy design, the
CI/CD setup, and the operational work are mine.

I mention it because it is true and because the distinction matters when you are
reading someone's code to judge how they think. I am happy to walk through any
decision in here.

## Status

Personal project, not accepting contributions. Published as a work sample.
