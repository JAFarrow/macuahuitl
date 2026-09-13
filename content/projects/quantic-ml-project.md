---
title: Quantic ML Project
status: fruited
tech:
  - Python
  - scikit-learn
summary: Malware / clearware classification, from dataset to deployed API
repo: https://github.com/JAFarrow/quantic_ml_project
link: https://quantic-ml-frontend.vercel.app/
draft: false
created: 2026-02-15
modified: 2026-09-13
---
## What
It classifies Windows PE (executable) metadata as malware or clearware. A Flask API serves a tuned XGBoost model with batch JSON and CSV-upload endpoints, and a small React frontend handles manual row entry and CSV upload, handing back a verdict per row.

## Why
Quantic coursework, which I used as an excuse to run the full ML lifecycle properly: from dataset and feature engineering through model selection and tuning, all the way to a deployed, test-covered inference API.

## How
- I trained on the Brazilian Malware Dataset (47,701 rows) with sklearn `ColumnTransformer` pipelines doing median imputation, date decomposition, missingness flags, and DLL token counts. All feature engineering happens inside the pipeline, so nothing leaks between train and test.
- A 7-model cross-validation leaderboard, including a torch MLP, put XGBoost on top with a test AUC of 0.9996 and 99.55% accuracy; I then tuned it with RandomizedSearchCV (n_iter=100, 5-fold).
- The API is marshmallow-validated with a 5,000-row cap, and CI runs pytest before firing a Render deploy hook. The frontend (on Vercel) was built against a written API contract (`.opencode/api_contract.md`) that mirrors the backend schema exactly.

## Status / Learnings
Complete and deployed. A few things stuck with me: namely that writing the API contract down keeps two decoupled repos in sync and keeping transforms inside sklearn Pipelines is the cleanest defense against leakage.
