# CLAUDE.md

Backend API for Ikou — a Go service (chi, MySQL) serving a separate frontend app.

Push to `releases/production` triggers `.github/workflows/main.yml`, which SSHes into the prod host and rebuilds the Docker image there.

Code layout — repository interface, `MockDB` contract, dto/model/mapper, auth, HTTP helpers — is in the **api-architecture** skill. Load it instead of rediscovering.

## Development

- **Never assume behavior.** Unclear expected behavior or business rule → stop and ask the owner before writing code.
- **No comment noise.** Rationale, edge cases and decisions go in the task doc, not inline. Comment only where the code cannot speak for itself.
- **Names carry the explanation.** A good type, variable or method name removes the need for the comment.
- **Record design corrections.** A design changed or corrected mid-task goes in a doc — the decision and the why, not a narrative.
- **Don't reformat files you aren't otherwise changing.** Standard gofmt, no custom style.
- **Ephemeral AI plans → `history/`.** Permanent documentation → `README.md`.

## Commits and PRs

- **Branch before committing.** Check the current branch first. On `master` or `releases/production`, cut a feature branch off `master` before staging anything — don't start work and sort the branch out afterwards.
- **Branch names: short, self-explanatory, prefixed.** `feat/` for features, `fix/` for bugs, `chore/` for maintenance, then what the work is in a few words — `feat/place-search-pagination`, `fix/refresh-token-expiry`. No ticket numbers alone, no dates, no initials.
- **One feature, one commit.** Keep each commit to a single objective so a reviewer reads it in one pass. Unrelated changes get their own commit — never bundle them in.
- **State the verification result** in the commit message: what was checked, what happened, in plain language — not raw tool output.
- **One branch per PR**, off `master`.
- **Sync with master before raising a PR to it.** `git fetch origin` then `git merge origin/master` into the feature branch, resolve anything that surfaces, and re-run `make check` before opening the PR. Never raise a PR from a branch that has not seen the current master.
- **Verify before raising.** Run `make check`, then read back the implementation and confirm the commit stands alone as working code.

## Safety (non-negotiable)

Always ask before:

- `git commit`, `git push`, any `git` command that changes history or remote state
- Any `gh` command (`pr create`, `issue`, `repo`, etc.)
- `make clean`, `rm`, `rm -rf`, or any destructive file operation
- Installing dependencies (`go get`, `brew install`, etc.)

Before running any gated command, present: the reason, the exact command, the expected effect, the risk level.

Never commit directly to `master` or `releases/production` — always a feature branch.

## Pre-commit security scan

Before every commit, read what is actually staged — `git diff --staged` — and do not commit until it is clear. Look for:

- Credentials and secrets: passwords, API keys, tokens, private keys, JWT secrets, connection strings with embedded credentials
- Config and env files that should stay local: `app.env`, `prod.env`, anything `*.env`
- Internal detail that shouldn't be public: real hostnames, database endpoints, internal IPs, account IDs, bucket names
- Personal data: real names, emails, phone numbers in seed data, fixtures or tests
- Local noise: `.DS_Store`, IDE config, build output, data dumps

Found something → stop, don't commit, tell the owner what and where. Remember this repo is **public**, so a bad commit is a disclosure, and removing it later needs a history rewrite.

Never commit `app.env` or `prod.env`; they hold live credentials. `app.env` is not in `.dockerignore` while the Dockerfile does `COPY . .`, so a locally built image bakes those credentials in — don't push one.

## Commands

Everything goes through `make`; run `make` alone for the full list.

```sh
make check          # fmt-check + vet + test — the gate before committing
make test           # go test ./...
make fmt            # gofmt -w .
make start          # build, then run locally in the foreground
make docker-up      # build and run in Docker on :9001
make docker-down    # stop and remove the containers
```

Single test: `go test ./internal/helper -run TestPasswordMatches -v`.

No linter configured, so `make check` is gofmt + vet + tests only.

## Environment

`app.env` at the repo root (gitignored). Viper hardcodes the config name `app`, so the filename must be exactly that. Keys:

```
DB_DRIVER, DB_SOURCE, SERVER_ADDRESS, PORT, FRONTEND_ADDRESS,
JTW_ACCESS_SECRET, JTW_REFRESH_SECRET
```

The JWT secrets are read via `os.Getenv`, which viper does not populate — they take effect under `make docker-up`, not `make start`. The api-architecture skill explains why.
