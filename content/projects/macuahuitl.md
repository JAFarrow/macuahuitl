---
title: macuahuitl
created: 2026-08-24
status: evergreen
tech:
  - Go
  - SvelteKit
summary: Obsidian vault in, single Go binary out. HTML for humans, markdown for agents.
repo: https://github.com/JAFarrow/macuahuitl
link: https://justin-farrow-dev.com
draft: false
modified: 2026-09-08
---
## What
This site. An Obsidian vault is embedded into a single Go binary that serves prerendered HTML to humans and the raw markdown verbatim to agents (`.md` URLs or `Accept: text/markdown`).

## Why
The vault is where I already write — it should be the source of truth, not an export step. And agents are first-class readers now, so markdown is served as a first-class format alongside `llms.txt` and a generated sitemap.

## How
- Go stdlib only, zero third-party deps; SvelteKit (`adapter-static`) build + vault embedded via `go:embed`; ~16 MB scratch image on Render.
- Draft frontmatter → 404 everywhere; serving confined to `projects/`, `cv/`, `attachments/`, `about.md`.
- Hand-rolled OTLP/JSON access-log exporter to Grafana Cloud (the official OTel SDK would cost ~16 modules); one structured `slog` line per request, typed markdown/html/attachment.

## Status / Learnings
Live and evergreen. Hard-won: `trailingSlash: 'always'` is load-bearing for `http.FileServerFS`; `http.ServeMux` panics on mid-segment wildcards (hence one `sectionHandler` per section dispatching on suffix + `Accept`); the draft check is a line-prefix scan, deliberately not a YAML parse.
