# AGENTS.md

This repository is a dependency override module for custom Caddy builds.

It provides no Caddy runtime functionality. Its purpose is to participate in
the Go module graph so that patched minimum versions for selected dependencies
can be enforced centrally in `xcaddy` builds.

## Important context

- `caddy_anchor.go` blank-imports `github.com/caddyserver/caddy/v2`.
- `overrides.go` contains temporary blank imports for modules that need to be
  bumped ahead of Caddy, with `govulncheck` references for each active
  override.
- `go.mod` records the effective minimum versions.
- `README.md` is the short human-facing document.
- Some findings, especially standard-library findings, require a patched Go
  toolchain and cannot be solved by module overrides alone.

Default behavior is to follow Caddy's dependency versions. Only add entries to
`overrides.go` when a temporary security bump is needed before Caddy updates.

Keep `README.md` short. Put the operational detail in this file instead of
duplicating longer procedures in both places.

## When given a govulncheck finding

Assume the user wants the repository updated so the vulnerable module is bumped
to at least the reported fixed version.

For each vulnerability:

1. Read the reported module path, found version, fixed version, and vulnerable
   symbols.
2. Add a blank import in `overrides.go` from the affected module.
3. Add a short comment above that import with the `GO-...` ID or IDs, the found
   version, the fixed version, and a note to remove the override once it is no
   longer needed.
4. Prefer importing the package that contains a reported vulnerable symbol.
5. Run `go get <module>@<fixed-version>`.
6. Run `go mod tidy`.
7. Run `go test ./...`.
8. Confirm the new direct requirement is present in `go.mod`.
9. Update the short override list in `README.md` if the set of active overrides
   changed.

If the finding is in the Go standard library instead of a module:

1. Read the fixed Go version from the `govulncheck` output.
2. Update the Go version in `go.mod` if the repository still declares an older
   version.
3. Do not add an entry to `overrides.go` for standard-library findings.
4. Tell the user that the final Caddy binary must be rebuilt with that newer Go
   toolchain.

## Import selection rules

- Do not create a fake Caddy plugin or dummy handler.
- Do not add runtime code.
- Use the smallest useful blank import that anchors the target module in the
  package graph.
- If `govulncheck` reports symbols from a subpackage, prefer importing that
  subpackage.
- If only the module root is relevant, importing the module root is fine.

Examples:

- `go.opentelemetry.io/otel` vulnerability with symbols in `baggage` or
  `propagation`: `_ "go.opentelemetry.io/otel"`
- `github.com/go-chi/chi/v5` vulnerability in `middleware.RealIP`:
  `_ "github.com/go-chi/chi/v5/middleware"`
- `github.com/klauspost/compress` vulnerability in `s2.NewDict`:
  `_ "github.com/klauspost/compress/s2"`

The initial scan that led to the current overrides included:

- `GO-2026-6061` for `google.golang.org/grpc`, found in `v1.81.0`, fixed in `v1.82.1`
- `GO-2026-5841` for `github.com/klauspost/compress`, found in `v1.18.6`, fixed in `v1.18.7`
- `GO-2026-5774`, `GO-2026-5775`, and `GO-2026-5777` for `github.com/go-chi/chi/v5`, found in `v5.2.5`, fixed in `v5.3.0`
- `GO-2026-5158` for `go.opentelemetry.io/otel`, found in `v1.43.0`, fixed in `v1.44.0`

That same scan also reported:

- `GO-2026-5970` for `golang.org/x/text`, found in `v0.37.0`, fixed in `v0.39.0`
- `GO-2026-5942` for `golang.org/x/net`, found in `v0.55.0`, fixed in `v0.56.0`
- `GO-2026-5932` for `golang.org/x/crypto/openpgp`, with no fixed version
- `GO-2026-5856` in `crypto/tls`, fixed in Go `1.26.5`
- `GO-2026-4970` in `os`, fixed in Go `1.26.5`

Not every finding needs an explicit override in this repository. Some are
resolved transitively by other bumps, some have no module-level fix, and some
require a Go toolchain update instead.

## Build and verification rules

When telling the user how to rebuild Caddy, always use the local checkout:

```bash
xcaddy build \
  --with github.com/steffenbusch/caddy-dependency-minimums=. \
  --with <plugin-1> \
  --with <plugin-2>
```

The `=.` suffix is important. Without it, `xcaddy` may build from another
source and ignore local changes. This module is an addition to the real plugin
set, not a replacement for it.

To inspect what actually ended up in the binary, use:

```bash
go version -m ./caddy
```

If the user says a vulnerability is still present after a rebuild, first check
the embedded module versions in the built binary before assuming the override
failed.

## Override cleanup rules

Do not let `overrides.go` grow forever.

When an override may no longer be needed:

1. Remove the blank import from `overrides.go`.
2. Remove or update the matching entry in `README.md`.
3. Run `go mod tidy`.
4. Rebuild Caddy with the normal full set of `--with` plugins.
5. Check the binary with `go version -m ./caddy`.
6. Run `govulncheck` again.

If the fixed-or-newer version is still present and the vulnerability does not
reappear, keep the override removed.

## Editing rules for this repository

- Keep the repository minimal.
- Prefer updating `overrides.go` over creating more files unless there is a
  strong reason.
- Remove temporary overrides once Caddy itself requires a sufficient version.
- Keep the README concise and keep the detail here.
