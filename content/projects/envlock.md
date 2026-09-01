---
title: EnvLock
status: fruited
tech:
  - TypeScript
  - NestJs
  - PostgreSQL
summary: Lightweight, CLI based secrets management solution
repo: https://github.com/JAFarrow/EnvLock
link: https://envlock-api.onrender.com
draft: false
created: 2026-06-08
modified: 2026-09-01
---
## What
Secrets live in a NestJS + PostgreSQL API (project/environment-scoped, with a plain-HTML dashboard); a zero-dependency npm CLI (`envlock run <cmd>`) fetches them at runtime and injects them into a subprocess's environment. No `.env` files on disk.

## Why
`.env` files leak — into git history, shell history, screenshots. Wanted 1Password-style secret injection without standing up Vault-grade infrastructure.

## How
- npm workspaces monorepo: NestJS 11 API (TypeORM, argon2, JWT) + published CLI `@jfarrow777/envlock`.
- AES-256-GCM with AAD binding each ciphertext to its `secretId`/`environmentId` — blocks swapping ciphertexts between secrets/environments; versioned keys/format allow rotation.
- PATs are 32 random bytes, only the SHA-256 hash stored; CLI prefers `ENVLOCK_PAT` env var over `--pat` (shell-history leakage).
- `envlock doctor` diffs `.env.example` keys against stored keys without ever fetching values.

## Status / Learnings
Live on Render (free tier, managed Postgres, migrations at startup); CLI at v0.2.4 on npm; 31 API spec files + a 666-line design doc. Learnings: server-side RBAC independent of the UI is non-negotiable; the static-HTML dashboard kept the stack simple but is starting to creak.
