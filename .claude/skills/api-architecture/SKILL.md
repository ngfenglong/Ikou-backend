---
name: api-architecture
description: How the Ikou Go API is wired — the layer pattern from route to SQL, the Repository/MockDB contract, dto/model/mapper, JWT auth, and the shared HTTP response helpers, plus the known landmines not to trip or "fix" by accident. Use when adding or changing an endpoint, a query, a DTO or a model, anything under api/ or internal/, or writing tests that need MockDB.
---

# Ikou API architecture

Go 1.19, chi, MySQL over `database/sql`. Layers are packages; each domain repeats the same shape.

## Where to look

Read these for the current inventory — mounts, routes and queries change, so don't trust a list written here.

| Need | File |
|---|---|
| what's mounted | `api/server.go` → `setupRoutes` |
| routes for a domain | `api/routes/<domain>_route.go` |
| handlers | `api/controllers/<domain>_controller.go` |
| the query contract | `api/repositories/repository.go` |
| query implementations | `api/repositories/<domain>_repository.go` |
| test double | `api/testutil/mock_db.go` |
| JWT | `internal/helper/auth.go` |
| request/response plumbing | `internal/helper/helper.go` |

Shapes: `api/dto/` (wire), `api/models/` (domain), `api/mapper/` (translation).

## One request, end to end

`GET /api/places` is the worked example. Every domain follows it.

1. `cmd/main.go` — load config, build `*store.Store`, hand both to `api.Application`
2. `api/server.go` — `s.Router.Mount("/api/places", routes.PlaceRoutes(s.App.Store))`
3. `api/routes/place_route.go` — `mux.Get("/", middleware.ExtractTokenMiddleware(placeController.GetAllPlaces))`
4. `api/controllers/place_controller.go` — `userID := r.Context().Value(middleware.UserIDKey).(string)`, then `pc.store.DB.GetAllPlaces(userID)`
5. `api/repositories/place_repository.go` — `DBModel` runs the SQL
6. out via `helper.WriteJSONResponse(w, http.StatusOK, places)`

A new domain is those six steps plus one `Mount` line. Controllers are `New<Name>Controller(store)` holding a `store` field; handlers are methods with the `http.HandlerFunc` signature.

## Adding a query touches three files

Every time, all three:

1. `api/repositories/repository.go` — the interface method
2. the domain's `*_repository.go` — the `DBModel` implementation
3. `api/testutil/mock_db.go` — a `…Func` field plus its delegating method

`MockDB` delegates to one `…Func` field per method. Skip step 3 → build breaks. Add the field but leave it nil in a test that reaches it → nil-func panic, not a clear failure. Tests set only the fields they exercise.

Controllers depend on `repository.Repository`, never on `DBModel`. The package is `repository` while the directory is `repositories`, so imports alias it:

```go
repository "github.com/ngfenglong/ikou-backend/api/repositories"
```

## Landmines

Pre-existing, deliberate to leave alone unless the task is to fix them — and easy to copy by accident when following a neighbouring file.

- **`log.Fatalf` in handlers.** The place and codestable controllers call it on a query or JSON error, which `os.Exit`s the whole server on one bad request. Never copy this into a new handler — use the `helper` error writers below.
- **`VerifyAccessToken` type-asserts every claim unchecked** (`claims["id"].(string)`), so a token missing one panics instead of failing validation. Adding or renaming a claim breaks already-issued tokens.
- **JWT secret names are transposed.** Code reads `JTW_ACCESS_SECRET` / `JTW_REFRESH_SECRET` via `os.Getenv` at package init. Viper's `app.env` load does not populate `os.Getenv`, so they're empty under `make start` and only arrive under `make docker-up`. `app.env` defines `JWT_ACCESS_SECRET`, which nothing reads — access tokens are signed with an empty key. Fixing it means changing code and env together.
- **Dead controller fields.** `NewServer` sets `app.PlaceController` / `AuthController` / `CodestableController`, but each `routes.*Routes` builds its own controller from the store.
- **CORS is hardcoded** to `{"https://*", "http://*"}`; the `ALLOWED_ORIGINS` env var is read by nothing.
- **The dto/model split is partial.** Place queries return `[]*dto.PlaceDTO` straight from the repository with no model between. A new endpoint whose DTO differs from its model gets a mapper in `api/mapper/` — don't widen the shortcut or cast at the call site.

## Auth

`middleware.ExtractTokenMiddleware` is permissive by design: strips `Bearer `, verifies, puts `userID`/`userName` in the request context, falls back to empty strings when the token is missing or invalid. It **does not reject unauthenticated requests** — place routes are browsable anonymously and each handler decides what an empty `userID` means. Read them via `middleware.UserIDKey` / `UserNameKey`, never raw strings.

HS256; access tokens 3 days, refresh 7. Refresh tokens are persisted — `InsertToken` clears the user's expired rows first, and `FetchRefreshTokenFromDB` is what makes logout effective.

## HTTP helpers

`internal/helper/helper.go` is the response layer — use it, don't write to the `ResponseWriter` directly.

- `WriteJSONResponse(w, status, data, headers...)`
- `ReadJSON(w, r, &dst)` — 1 MiB cap, rejects a second JSON value in the body
- `BadRequest` / `Unauthorized` / `InternalServerError` / `InvalidCredential` / `ConflictErrorResponse`

Every error shares one envelope: `{"error": true, "message": "..."}`.

`PasswordMatches` returns `(false, nil)` for a mismatch, `(false, err)` for a real bcrypt failure. Keep them distinct — collapsing them reports an internal error as "wrong password".

## Tests

Table-driven `_test.go` beside the code, `httptest` + `testutil.MockDB`. The repository layer is untested (needs a live DB); coverage is in controllers, mapper, middleware, helper, util.
