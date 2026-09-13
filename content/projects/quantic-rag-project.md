---
title: Quantic RAG Project
status: fruited
tech:
  - Python
  - Langchain
  - Pinecone
summary: A chatbot that answers HR policy questions with citations
repo: https://github.com/JAFarrow/quantic_rag_project
link: https://quantic-rag-project.onrender.com
draft: false
created: 2026-04-01
modified: 2026-09-13
---
## What
A RAG chatbot over a fictional company's HR policy PDFs, with answers grounded in citations. It's FastAPI + LangChain + Pinecone, deployed as a single Render service that serves both the API and a built-in dark-themed chat UI.

## Why
A Quantic capstone-style build. I wanted a production-shaped RAG pipeline end to end: ingestion, retrieval, grounded generation, and a measured evaluation instead of vibes.

## How
- An ingestion CLI chunks the policy PDFs with `RecursiveCharacterTextSplitter`, using separators tuned for policy documents (`\n• `, `\n- `), embeds them with `text-embedding-3-small`, then wipes and upserts the Pinecone namespace so re-runs are idempotent.
- The LCEL chain runs a retriever (top_k=12) into a numbered-citation context builder and then `gpt-5.4-mini`. Only citations actually referenced as `[n]` in the answer get returned, and an insufficient-context fallback is enforced in the prompt and post-checked in code.

## Status / Learnings
Feature-complete, though the live URL currently 503s thanks to Render's free-tier spin-down. A 20-question eval drove groundedness/citation accuracy from 80% to 90% after I shrank the chunks and tightened the prompts; similarity-score filtering turned out counterproductive (min_score=0.0 won). The AI-assisted dev process is documented openly in `ai-tooling.md`.
