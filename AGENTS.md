# AGENTS.md

**macuahuitl**: Obsidian vault in, single Go binary out. HTML for humans, markdown for agents.

## Layout

- `main.go` — routing and server lifecycle. `logging.go` — access logging (slog JSON to stdout). Both are `package main` at the repo root; no `cmd/`, `internal/`, or `pkg/` — the codebase is too small for layout ceremony, and `go:embed` patterns can't reach parent dirs anyway.
- `content/` — the Obsidian vault (`projects/`, `templates/`, plus `about.md` for the home page). Markdown with YAML frontmatter, standard markdown links.
- `frontend/` — SvelteKit app (Svelte 5, mdsvex, `adapter-static`). Build output goes to `frontend/build/`.
- `frontend/static/robots.txt` — crawler policy (allow all), copied verbatim into `frontend/build/`.
- `frontend/src/lib/site.ts` — shared canonical origin constant (`SITE_URL`) used by the build-time generators.
- `frontend/src/routes/llms.txt/+server.ts` — prerendered endpoint that generates `build/llms.txt` from the vault data loaders.
- `frontend/src/routes/sitemap.xml/+server.ts` — prerendered endpoint that generates `build/sitemap.xml` (home plus all non-draft projects at their canonical trailing-slash URLs). `lastmod` is taken from a `modified` frontmatter key (falling back to `created`); the home page uses `about.md`'s `modified`.
- `frontend/src/routes/projects.md/+server.ts` — prerendered endpoint that generates `build/projects.md` (markdown twin of the project gallery).
- The home page (`frontend/src/routes/+page.svelte`) renders the `about.md` body followed by the project gallery (cards linking to `/projects/{slug}` detail pages, which carry a `← Projects` backlink home). There is no navbar — just the site-title header. Contact links live in `about.md` frontmatter as one flat key per link (`email:`, `linkedin:`, `github:`), selected and ordered by the `CONTACT_KEYS` allowlist in `$lib/data/about.ts` — adding or renaming a link means editing both files. Links are rendered site-wide by `+layout.svelte` as a sticky footer (flexbox: `.site` is a `min-height: 100vh` column, `main` grows). Like markdown-rendered external links, footer links open in a new tab.
- `justfile`, `.github/workflows/go.yml`, `Dockerfile`, `render.yaml` — tooling, CI, and Render deployment.

## Commands

- `just dev` — `go run .` (run `cd frontend && npm run dev` alongside for frontend work; vite proxies `/api` to :8080).
- `just build` — frontend build, then `go build -o macuahuitl .`
- `just lint` — `go vet ./...` + `golangci-lint run` (v2 defaults; must stay clean)
- `docker build -t macuahuitl .` — the exact build Render runs.

## Hard rules

- **Go: stdlib only, zero third-party deps.** No router, no YAML parser, no logging libs.
- **`//go:embed all:frontend/build all:content` needs `frontend/build/` on disk**, but it's gitignored — always build the frontend before `go build` (CI and the Dockerfile already do).
- **Serve markdown verbatim.** No frontmatter stripping, no link rewriting.
- **Drafts are private.** `draft: true` frontmatter → 404 everywhere (the `.md` route, the `Accept: text/markdown` path, and HTML — the frontend never prerenders drafts).
- **Don't expose the vault beyond `projects/` and `about.md`.** `.obsidian/` and `templates/` are embedded but must stay unreachable.

## Conventions & gotchas

- **Code is comment-free by convention.** Rationale and gotchas live in this file, not in source comments — don't add comments when editing; document here instead. (Functional directives like `//go:embed` and the Dockerfile `# syntax=` line are not comments and stay.)
- **Canonical origin is `https://www.justin-farrow-dev.com`** — Cloudflare 301s apex → www at the edge. There is no Go-side host redirect: `macuahuitl.onrender.com` and direct-origin apex hits are served as-is. The canonical origin lives in `frontend/src/lib/site.ts` and is imported by the sitemap, `llms.txt`, and the index `.md` generators. `robots.txt` separately hardcodes the same origin in its `Sitemap:` line; update both if it ever moves.
- **`trailingSlash: 'always'` is load-bearing.** Routes prerender as `route/index.html` so `http.FileServerFS` resolves everything; slash-less URLs 301 to the slash form. Don't revert without adding Go-side `.html` mapping.
- **No mid-segment wildcards in `http.ServeMux`** (`"/projects/{slug}.md"` panics at startup). That's why one `sectionHandler` per section dispatches on `.md` suffix and `Accept` header.
- **Markdown negotiation works on both URL forms** — slash-less (`/projects/x`, also the only place a `.md` suffix is recognized) and canonical trailing-slash (`/projects/x/` via `{$}`-anchored patterns). The home page follows the same pattern: `/about.md` serves `content/about.md`, and `Accept: text/markdown` on `/` does too. The project gallery also exposes a generated markdown twin at `/projects.md`. Agents should prefer the explicit `.md` URLs anyway.
- **Draft check is a frontmatter line-prefix scan**, not a YAML parse — keep it that way.
- **The vite build reads `content/` directly** (`import.meta.glob` over `../../../../content/projects/*.md` and `../../../../content/about.md` in the `$lib/data` loaders), so any build environment (CI, Docker) needs both directories side by side.
- **`/llms.txt` and `/sitemap.xml` are build-generated, never hand-edited.** They come from the same loaders (`$lib/data/projects`, `$lib/data/about`) as the site, so vault changes regenerate them on the next deploy. `/projects.md` is also build-generated from the same loaders (a markdown twin — the link index for projects). All three land in `frontend/build/`; `llms.txt` and `sitemap.xml` are served directly by the Go file server, while `projects.md` has an explicit Go route so its `Content-Type` and access-log `type` stay consistent with other markdown responses.
- No Go tests yet; verify behavior with the curl matrix (HTML, `.md` verbatim + content type, `Accept` negotiation on `/projects/x/` and `/`, index markdown at `/projects.md`, draft 404, `/robots.txt`, `/llms.txt`, `/sitemap.xml`, `/api/health`).

## Favicon / brand colour

- The favicon is a minimal, left-slanted macuahuitl glyph: `frontend/static/favicon.svg`.
- The obsidian blades use the brand accent colour from `frontend/src/app.css`: `#0645ad` in light mode and `#7cb3ff` in dark mode, switched via `prefers-color-scheme` inside the SVG.
- The club/paddle is neutral (`#1c1c1c` in light mode, `#e4e4e4` in dark mode), matching `--fg`.

## Observability

- **Access logs, not traces.** The goal is usage insight (HTML vs markdown serving, agents vs humans), which is a per-request counting question — a log signal. One structured `slog` JSON line per request goes to stdout, the only transport — there is no in-app exporter and no app-side logging config. Render's workspace log stream (syslog TLS, dashboard-only config) forwards stdout to Better Stack; Render's own log retention is the fallback when no stream is configured.
- **stdout JSON is aligned to Better Stack's canonical fields.** `setupLogging()` (in `logging.go`, called first in `main`) installs a `slog.JSONHandler` whose `ReplaceAttr` renames top-level `time` → `dt` and `msg` → `message`: Better Stack parses `dt` as the event time and shows `message` as the log line, with all other keys becoming searchable fields. Render additionally maps the JSON `level` field to the syslog priority it attaches when streaming.
- **Log fields:** `method`, `route` (`r.Pattern`, e.g. `GET /projects/{slug}`), `path`, `status`, `duration_ms` (float, µs precision — whole-ms truncation zeroes nearly every request on this site), `bytes`, `type`, `user_agent`, `referer`, `ip` (first `X-Forwarded-For` hop). `type` is `markdown` | `html`: markdown handlers mark `markdown` via a `*requestMeta` planted in the request context (even for draft/missing/invalid-slug 404s — intent, with outcome in `status`); everything else is `html`.
- **`/api/health` is never logged** — Render polls it constantly; pure noise.
- Example Better Stack query: filter on `type` in Live tail (e.g. `type:markdown`) for the html/markdown split.
- Graceful shutdown (SIGINT/SIGTERM) remains so in-flight requests finish on deploy; logs are unbuffered stdout, so there is nothing to flush. The scratch image carries no CA bundle — the binary makes no outbound calls.

## Deployment

Render, via `Dockerfile` (node → golang → scratch, ~16 MB image) + `render.yaml` blueprint (starter plan, frankfurt, health check `/api/health`, auto-deploy on `main`). Publish flow: commit vault changes → push → auto-deploy. Logs reach Better Stack via Render's log stream: create a Render-platform source in Better Stack, then set the workspace's default log-stream destination to `<ingesting-host>:6514` with the source token (Render dashboard → Integrations > Observability). Blueprints have no log-stream field, so this is a one-time dashboard step; the app itself needs no env vars or secrets.

`llms.txt` and `sitemap.xml` ride the existing frontend build (`npm run build`) in every environment, so CI, the Dockerfile, and `render.yaml` need no special-casing for them.
