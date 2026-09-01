---
title: Quantic Rag Project
status: fruited
tech:
  - Python
  - Langchain
  - Pinecone
summary: RAG enabled chatbot
repo: https://github.com/JAFarrow/quantic_rag_project
link: https://quantic-rag-project.onrender.com
draft: false
created: 2026-04-01
modified: 2026-09-01
---
## What
RAG chatbot over a fictional company's HR policy PDFs, with citation-grounded answers. FastAPI + LangChain + Pinecone, single Render service serving both the API and a built-in dark-themed chat UI.

![Screenshot demonstrating the frontend view of the rag chatbot in use](rag_project.webp)

## Why
Quantic capstone-style build: a production-shaped RAG pipeline end to end — ingestion, retrieval, grounded generation, and measured evaluation rather than vibes.

## How
- Ingestion CLI chunks policy PDFs with `RecursiveCharacterTextSplitter` using separators tuned for policy docs (`\n• `, `\n- `), embeds with `text-embedding-3-small`, then wipes + upserts the Pinecone namespace so re-runs are idempotent.
- LCEL chain: retriever (top_k=12) → numbered-citation context builder → `gpt-5.4-mini` → only citations actually referenced as `[n]` in the answer are returned; an insufficient-context fallback is enforced in prompt and post-checked in code.

## Status / Learnings
Feature-complete; the live URL currently 503s (Render free-tier spin-down). A 20-question eval drove groundedness/citation accuracy from 80% → 90% after shrinking chunks and tightening prompts — and similarity-score filtering turned out counterproductive (min_score=0.0 won). `ai-tooling.md` frankly documents the AI-assisted dev process.
