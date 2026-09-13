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
modified: 2026-09-13
---
## What
This site. An Obsidian vault gets embedded into a single Go binary that serves prerendered HTML to humans and the raw markdown, verbatim, to agents (via `.md` URLs or `Accept: text/markdown`).

## Why
The vault is where I already write, so it should be the source of truth rather than an export step. Agents read the web now too, so markdown gets served on equal footing with HTML, alongside `llms.txt` and a generated sitemap.

## How
- Go stdlib only, zero third-party dependencies. The SvelteKit (`adapter-static`) build and the vault are embedded with `go:embed`, and everything ships as a ~16 MB scratch image on Render.
- `draft: true` frontmatter means a 404 everywhere, and serving is confined to `projects/` and `about.md`.
- Each request writes one structured `slog` JSON line to stdout, with Better Stack-aligned `dt`/`message` keys and a markdown/html type. Render's log stream forwards stdout to Better Stack, so the binary itself carries no exporter, no dependencies, and no secrets.

## Status / Learnings
Live and evergreen. 
