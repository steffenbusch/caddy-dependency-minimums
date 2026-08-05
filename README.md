# caddy-dependency-minimums

This repository is a helper module for custom Caddy builds with `xcaddy`.

It provides no Caddy functionality. Its only purpose is to participate in the
Go module graph so that patched minimum versions for selected dependencies can
be enforced centrally.

At the moment this repository is an experiment. Based on my current
understanding of Go module resolution, this approach appears to work as
intended, but I do not consider it a formally proven or universally recommended
solution.

> [!CAUTION]
> Every dependency upgrade can have negative side effects. For this
> experimental personal module, I try to assess that risk myself, but that
> assessment may be wrong. Some effects may only show up later in real usage,
> in areas that are not relevant to my own builds, or only in combination with
> other Caddy plugins. No warranty or guarantee is provided.

The final Caddy binary should always be verified after every build, for example
with `go version -m` and a vulnerability scan.

Some findings, especially in the Go standard library, cannot be fixed by module
overrides alone and require building with a sufficiently patched Go toolchain.
Use at least the Go version declared in [go.mod](./go.mod), and move forward
when `govulncheck` reports a standard-library fix that requires a newer Go
release.

## How it works

- [caddy_anchor.go](./caddy_anchor.go) blank-imports `github.com/caddyserver/caddy/v2`.
- [overrides.go](./overrides.go) contains temporary blank imports for override modules and records the related `govulncheck` references.
- [go.mod](./go.mod) records the effective minimum versions.
- [AGENTS.md](./AGENTS.md) contains the detailed maintenance workflow for future updates.

Default behavior is to follow Caddy. Add overrides only for temporary security
bumps that need to land before Caddy updates.

## Contribution note

This repository is maintained for my own Caddy builds and risk profile.

I may decline pull requests even if they are technically valid, especially when
they add dependency overrides for issues that do not affect my own use cases or
when the potential side effects are unclear.

## Building with xcaddy

Use the local workspace state when building. This module is added in addition
to the Caddy plugins you actually want in the production binary:

```bash
xcaddy build \
  --with github.com/steffenbusch/caddy-dependency-minimums=. \
  --with <plugin-1> \
  --with <plugin-2>
```

The `=.` part is important. Without it, `xcaddy` may use a different source
than your current local checkout.

Active overrides are documented directly in [overrides.go](./overrides.go),
including their `govulncheck` references and removal notes. For the update and
cleanup procedure, see [AGENTS.md](./AGENTS.md).

## Verifying what is really in the binary

```bash
go version -m ./caddy
```

This is the reliable way to confirm which versions actually ended up in the
final binary.

Overrides should be removed once Caddy itself or the normal plugin set already
pulls in a sufficient version. The cleanup procedure is documented in
[AGENTS.md](./AGENTS.md).
