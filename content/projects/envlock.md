---
title: EnvLock
status: fruited
tech:
  - TypeScript
  - NestJs
  - PostgreSQL
summary: A CLI that injects secrets into processes at runtime, keeping .env files off disk
repo: https://github.com/JAFarrow/EnvLock
link: https://envlock-api.onrender.com
draft: false
created: 2026-06-08
modified: 2026-09-13
---
## What
I keep secrets in a small NestJS + PostgreSQL API, scoped by project and environment, with a plain HTML dashboard on top. A zero-dependency npm CLI (`envlock run <cmd>`) fetches them at runtime and injects them into the subprocess's environment, so no `.env` files ever touch disk.

## Why
Given free reign for my MSc capstone, I decided to tackle a challenge I'd seen consistently within my professional life: managing environmental variables and secrets across a team of developers (and AI agents), preferably with a single source of truth and as few potential leakages as possible. 

## How
- It's an npm workspaces monorepo: a NestJS 11 API (TypeORM, argon2, JWT) plus the published CLI, `@jfarrow777/envlock`.
- Secrets are encrypted with AES-256-GCM, with AAD binding each ciphertext to its `secretId` and `environmentId`, so a ciphertext can't be swapped between secrets or environments. Keys and format are versioned, which leaves room for rotation.
- PATs are 32 random bytes and only the SHA-256 hash is stored. The CLI prefers the `ENVLOCK_PAT` env var over `--pat`, because flags end up in shell history.
- `envlock doctor` diffs the keys in `.env.example` against the stored keys without ever fetching values.

## Status / Learnings
A DB was provisioned during the marking phase of this project, however I've since let it spin down and so this project is no longer live. The encryption component was the most interesting aspect of this project that I enjoyed working on. Were this to enter production, quite a bit of work would still be required to ensure secure operation, most outstanding of which would be a proper sign-up mechanism. 
