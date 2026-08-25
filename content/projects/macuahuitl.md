---
title: macuahuitl
created: 2026-08-24
status: evergreen
tech:
  - Go
  - SvelteKit
summary: Obsidian vault in, single Go binary out. HTML for humans, markdown for agents.
repo: https://github.com/justin/macuahuitl
link: https://justin.example.com
draft: false
---

## What

This site. An Obsidian vault of markdown notes compiled into a static site - prerendered HTML for humans, raw markdown served alongside for agents.

## Why

I write in Obsidian anyway. Publishing should be `git push`, not a CMS. And every page should be equally readable by a browser or an LLM.

## How

SvelteKit with adapter-static prerenders every route from `content/` via mdsvex. A small Go server serves the build output plus the raw `.md` files, and handles webmentions and search.

![test](40k.png)