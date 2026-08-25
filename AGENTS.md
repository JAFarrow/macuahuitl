# AGENTS.md

**macuahuitl**: Obsidian vault in, single Go binary out. HTML for humans, markdown for agents.

## Layout

- `main.go` — the entire Go server. Deliberately one file; no `cmd/`, `internal/`, or `pkg/` — the codebase is too small for layout ceremony.
- `content/` — the Obsidian vault (`posts/`, `projects/`, `attachments/`, `templates/`). Markdown with YAML frontmatter, standard markdown links.
- `frontend/` — SvelteKit app (Svelte 5, mdsvex, `adapter-static`). Build output goes to `frontend/build/`.
- `docs/update.md` — running repository-state log; keep it current when you change things.
- `justfile`, `.github/workflows/go.yml`, `Dockerfile`, `render.yaml` — tooling, CI, and Render deployment.

## Commands

- `just dev` — `go run .` (run `cd frontend && npm run dev` alongside for frontend work; vite proxies `/api` and `/attachments` to :8080).
- `just build` — frontend build, then `go build -o macuahuitl .`
- `just lint` — `go vet ./...` + `golangci-lint run` (v2 defaults; must stay clean)
- `docker build -t macuahuitl .` — the exact build Render runs.

## Hard rules

- **Go: stdlib only, zero third-party deps.** No router, no YAML parser, no logging libs.
- **`//go:embed all:frontend build all:content` needs `frontend/build/` on disk**, but it's gitignored — always build the frontend before `go build` (CI and the Dockerfile already do).
- **Serve markdown verbatim.** No frontmatter stripping, no link rewriting.
- **Drafts are private.** `draft: true` frontmatter → 404 everywhere (the `.md` route, the `Accept: text/markdown` path, and HTML — the frontend never prerenders drafts).
- **Don't expose the vault beyond `posts/`, `projects/`, `attachments/`.** `.obsidian/` and `templates/` are embedded but must stay unreachable. Attachment serving stays confined to its `fs.Sub` (encoded `..%2F` traversal must 404), and directory listings are 404ed (`noListings`).

## Conventions & gotchas

- **`trailingSlash: 'always'` is load-bearing.** Routes prerender as `route/index.html` so `http.FileServerFS` resolves everything; slash-less URLs 301 to the slash form. Don't revert without adding Go-side `.html` mapping.
- **No mid-segment wildcards in `http.ServeMux`** (`"/posts/{slug}.md"` panics at startup). That's why one `sectionHandler` per section dispatches on `.md` suffix and `Accept` header.
- **Markdown negotiation works on both URL forms** — slash-less (`/posts/x`, also the only place a `.md` suffix is recognized) and canonical trailing-slash (`/posts/x/` via `{$}`-anchored patterns). Agents should prefer the explicit `.md` URLs anyway.
- **Draft check is a frontmatter line-prefix scan**, not a YAML parse — keep it that way.
- **The vite build reads `content/` directly** (`import.meta.glob` over `../../../../content/**/*.md`), so any build environment (CI, Docker) needs both directories side by side.
- No Go tests yet; verify behavior with the curl matrix (HTML, `.md` verbatim + content type, `Accept` negotiation, draft 404, attachment bytes, `/api/health`).

## Deployment

Render, via `Dockerfile` (node → golang → scratch, ~27 MB image) + `render.yaml` blueprint (free plan, frankfurt, health check `/api/health`, auto-deploy on `main`). Free tier spins down when idle. Publish flow: commit vault changes → push → auto-deploy.
