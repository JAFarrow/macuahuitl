# AGENTS.md

**macuahuitl**: Obsidian vault in, single Go binary out. HTML for humans, markdown for agents.

## Layout

- `main.go` — routing and server lifecycle. `logging.go` — access logging and the OTLP exporter. Both are `package main` at the repo root; no `cmd/`, `internal/`, or `pkg/` — the codebase is too small for layout ceremony, and `go:embed` patterns can't reach parent dirs anyway.
- `content/` — the Obsidian vault (`posts/`, `projects/`, `attachments/`, `templates/`). Markdown with YAML frontmatter, standard markdown links.
- `frontend/` — SvelteKit app (Svelte 5, mdsvex, `adapter-static`). Build output goes to `frontend/build/`.
- `frontend/static/robots.txt` — crawler policy (allow all), copied verbatim into `frontend/build/`.
- `frontend/src/routes/llms.txt/+server.ts` — prerendered endpoint that generates `build/llms.txt` from the vault data loaders.
- `frontend/src/routes/sitemap.xml/+server.ts` — prerendered endpoint that generates `build/sitemap.xml` (home, section indexes, all non-draft posts/projects at their canonical trailing-slash URLs). The `SITE_URL` constant in that file is the canonical origin.
- `justfile`, `.github/workflows/go.yml`, `Dockerfile`, `render.yaml` — tooling, CI, and Render deployment.

## Commands

- `just dev` — `go run .` (run `cd frontend && npm run dev` alongside for frontend work; vite proxies `/api` and `/attachments` to :8080).
- `just build` — frontend build, then `go build -o macuahuitl .`
- `just lint` — `go vet ./...` + `golangci-lint run` (v2 defaults; must stay clean)
- `docker build -t macuahuitl .` — the exact build Render runs.

## Hard rules

- **Go: stdlib only, zero third-party deps.** No router, no YAML parser, no logging libs.
- **`//go:embed all:frontend/build all:content` needs `frontend/build/` on disk**, but it's gitignored — always build the frontend before `go build` (CI and the Dockerfile already do).
- **Serve markdown verbatim.** No frontmatter stripping, no link rewriting.
- **Drafts are private.** `draft: true` frontmatter → 404 everywhere (the `.md` route, the `Accept: text/markdown` path, and HTML — the frontend never prerenders drafts).
- **Don't expose the vault beyond `posts/`, `projects/`, `attachments/`.** `.obsidian/` and `templates/` are embedded but must stay unreachable. Attachment serving stays confined to its `fs.Sub` (encoded `..%2F` traversal must 404), and directory listings are 404ed (`noListings`).

## Conventions & gotchas

- **Canonical origin is `https://www.justin-farrow-dev.com`** — Cloudflare 301s apex → www at the edge. There is no Go-side host redirect: `macuahuitl.onrender.com` and direct-origin apex hits are served as-is. Hardcoded in two places — the sitemap `SITE_URL` constant and the `robots.txt` `Sitemap:` line; change them together if it ever moves.
- **`trailingSlash: 'always'` is load-bearing.** Routes prerender as `route/index.html` so `http.FileServerFS` resolves everything; slash-less URLs 301 to the slash form. Don't revert without adding Go-side `.html` mapping.
- **No mid-segment wildcards in `http.ServeMux`** (`"/posts/{slug}.md"` panics at startup). That's why one `sectionHandler` per section dispatches on `.md` suffix and `Accept` header.
- **Markdown negotiation works on both URL forms** — slash-less (`/posts/x`, also the only place a `.md` suffix is recognized) and canonical trailing-slash (`/posts/x/` via `{$}`-anchored patterns). Agents should prefer the explicit `.md` URLs anyway.
- **Draft check is a frontmatter line-prefix scan**, not a YAML parse — keep it that way.
- **The vite build reads `content/` directly** (`import.meta.glob` over `../../../../content/posts/*.md` and `../../../../content/projects/*.md` in the `$lib/data` loaders), so any build environment (CI, Docker) needs both directories side by side.
- **`/llms.txt` and `/sitemap.xml` are build-generated, never hand-edited.** They come from the same draft-filtered loaders (`$lib/data/posts`, `$lib/data/projects`) as the site, so vault changes regenerate them on the next deploy. `robots.txt`, `llms.txt`, and `sitemap.xml` land in `frontend/build/` and are served by the Go file server — no Go routes involved.
- No Go tests yet; verify behavior with the curl matrix (HTML, `.md` verbatim + content type, `Accept` negotiation, draft 404, attachment bytes, `/robots.txt`, `/llms.txt`, `/sitemap.xml`, `/api/health`).

## Observability

- **Access logs, not traces.** The goal is usage insight (HTML vs markdown serving, agents vs humans), which is a per-request counting question — a log signal. One structured `slog` line per request always goes to local stdout (Render log retention is the fallback); when `OTEL_EXPORTER_OTLP_ENDPOINT` is set, the same records are batched to Grafana Cloud's OTLP gateway at `{endpoint}/v1/logs` (Loki).
- **Hand-rolled OTLP/JSON on purpose.** The official OTel Go SDK would add ~16 modules, breaking the zero-deps hard rule for marginal benefit at this scale. The exporter (`logging.go`) reads the standard env vars Grafana issues: `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_HEADERS` (`k=v,k=v`; values decoded with `url.PathUnescape` — `QueryUnescape` would corrupt `+` in the base64 Basic token), `OTEL_SERVICE_NAME` (default `macuahuitl`). Endpoint unset → exporter is nil, zero overhead.
- **Log fields:** `method`, `route` (`r.Pattern`, e.g. `GET /posts/{slug}`), `path`, `status`, `duration_ms`, `bytes`, `type`, `user_agent`, `referer`, `ip` (first `X-Forwarded-For` hop). `type` is `markdown` | `html` | `attachment`: `serveMarkdown` marks `markdown` via a `*requestMeta` planted in the request context (even for draft/missing 404s — intent, with outcome in `status`); `/attachments/` prefix → `attachment`; everything else `html`.
- **`/api/health` is never logged** — Render polls it constantly; pure noise.
- Example Loki query: `sum by (type) (count_over_time({service_name="macuahuitl"} [1h]))`.
- The Dockerfile copies `ca-certificates.crt` from the golang builder stage into scratch — Grafana's gateway is HTTPS and scratch has no roots. Graceful shutdown (SIGINT/SIGTERM) flushes buffered records on deploy.

## Deployment

Render, via `Dockerfile` (node → golang → scratch, ~16 MB image) + `render.yaml` blueprint (starter plan, frankfurt, health check `/api/health`, auto-deploy on `main`). Publish flow: commit vault changes → push → auto-deploy. Grafana credentials are `sync: false` env vars in the blueprint — set `OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_EXPORTER_OTLP_HEADERS` once in the Render dashboard.

`llms.txt` and `sitemap.xml` ride the existing frontend build (`npm run build`) in every environment, so CI, the Dockerfile, and `render.yaml` need no special-casing for them.
