# Repository State — 2026-08-25

**macuahuitl**: Obsidian vault in, single Go binary out. HTML for humans, markdown for agents.

## Git

- Branch `main`, **4 commits ahead of `origin/main`** (unpushed).
- Uncommitted: the entire Go server (`main.go`, `go.mod`), CI (`.github/`), `justfile`, `.gitignore` (`/macuahuitl`), frontend `trailingSlash` change, `content/attachments/40k.png`, and this doc (`docs/`).

## Content (`content/`)

Obsidian vault is configured (`.obsidian/`) and seeded:

- `posts/` — 1 post (`test-post.md`)
- `projects/` — 1 project (`macuahuitl.md`, self-describing; embeds `![](40k.png)`)
- `templates/` — `post.md` and `project.md` frontmatter templates
- `attachments/` — `40k.png`, served by the Go server at `/attachments/*`

## Frontend (`frontend/`)

SvelteKit (Svelte 5, Kit 2.70, TypeScript, mdsvex) with `adapter-static`:

- Routes: home (posts + projects), `/posts`, `/posts/[slug]`, `/projects`, `/projects/[slug]`.
- Markdown loaded eagerly via `import.meta.glob` over `content/{posts,projects}/*.md`; drafts filtered, posts sorted by date.
- Custom remark plugin rewrites Obsidian bare image embeds to `/attachments/*`.
- `trailingSlash: 'always'` (`src/routes/+layout.ts`): prerendered routes emit as `route/index.html`, so the Go server's plain file server resolves every page with no mapping logic. Canonical URLs end in `/`; slash-less requests 301 to the slash form.
- Prerender tolerates 404s on `*.md`, `*.md/` (trailing-slash form), and `/attachments/*` — served by the Go server, not the frontend.

## Go server (`main.go`, stdlib only, zero deps)

Single-file server; `//go:embed all:frontend/build all:content` (the `all:` prefix keeps `_app/` assets and vault dotfiles):

- `GET /posts/{slug}` and `GET /projects/{slug}`: slug ending in `.md` **or** `Accept: text/markdown` → raw markdown verbatim (frontmatter included, `text/markdown; charset=utf-8`); otherwise HTML. `draft: true` frontmatter (line-prefix scan, no YAML lib) and `/`/`..` slugs → 404.
- `GET /attachments/` → file server over `content/attachments/`, confined via `fs.Sub` + `http.StripPrefix` (encoded `..%2F` traversal 404s; can't reach the rest of the vault).
- `GET /api/health` → `{"status":"ok"}`.
- `GET /` → `http.FileServerFS` over the embedded build.
- Port `8080` (`PORT` env override); 10s read/read-header/write, 60s idle timeouts; `log/slog` startup log.
- Note: `GET /posts/{slug}/` (trailing slash) does not match the section pattern, so `Accept: text/markdown` there returns HTML — agents should use the `.md` or slash-less form.

## Tooling

- `justfile`: `dev` (`go run .`), `build` (frontend build → `go build -o macuahuitl .`), `lint` (`go vet` + `golangci-lint`).
- CI (`.github/workflows/go.yml`): npm build (required — `frontend/build/` is gitignored but `go:embed` needs it) → `go build`, `go vet`, golangci-lint v2 defaults (errcheck, govet, ineffassign, staticcheck, unused).
- Verified: full curl matrix, draft 404s on all paths, `PORT` override, and binary-only serving with `frontend/build` + `content` deleted from disk.

## Not yet done

- Generated `/posts.md` and `/projects.md` index documents.
- OTLP telemetry.
- `llms.txt` (no `frontend/static/` exists yet, despite earlier assumptions).
- Webmentions and search (promised by README/project page).
- Graceful shutdown beyond timeouts; `/attachments/` directory root currently renders a stdlib listing.
