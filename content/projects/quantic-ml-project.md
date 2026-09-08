---
title: Quantic ML Project
status: fruited
tech:
  - Python
  - scikit-learn
summary: Machine learning playground for malware / clearware identification
repo: https://github.com/JAFarrow/quantic_ml_project
link: https://quantic-ml-frontend.vercel.app/
draft: false
created: 2026-02-15
modified: 2026-09-07
---
## What
Classifies Windows PE (executable) metadata as malware or clearware. A Flask API serves a tuned XGBoost model (batch JSON + CSV-upload endpoints); a small React frontend handles manual row entry and CSV upload with per-row verdicts.

## Why
Quantic coursework, taken as an excuse to run the full ML lifecycle properly: dataset → feature engineering → model selection → tuning → deployed, test-covered inference API.

## How
- Brazilian Malware Dataset (47,701 rows) → sklearn `ColumnTransformer` pipelines (median imputation, date decomposition, missingness flags, DLL token counts) — all feature engineering inside the pipeline, so no leakage.
- 7-model CV leaderboard (incl. a torch MLP); XGBoost won at test AUC 0.9996 / 99.55% accuracy, then RandomizedSearchCV (n_iter=100, 5-fold).
- marshmallow-validated API (5,000-row cap); CI runs pytest then fires a Render deploy hook. Frontend (Vercel) built against a written API contract (`.opencode/api_contract.md`) that mirrors the backend schema exactly.

## Status / Learnings
Complete and deployed. Learnings: an explicit written contract keeps decoupled repos in sync; keeping transforms inside sklearn Pipelines is the cleanest leakage defense; the committed `dist/` bundle has `localhost:5000` baked in — env config must happen at Vercel build time.
