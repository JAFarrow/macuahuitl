# syntax=docker/dockerfile:1

# Stage 1: prerender the SvelteKit site. The vite build reads markdown from
# content/ at the repo root, so both directories are required here.
FROM node:26-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json frontend/
RUN cd frontend && npm ci
COPY frontend/ frontend/
COPY content/ content/
RUN cd frontend && npm run build

# Stage 2: build the static Go binary (stdlib only, no CGO).
FROM golang:1.26-alpine AS server
WORKDIR /app
COPY go.mod main.go ./
COPY content/ content/
COPY --from=frontend /app/frontend/build frontend/build
RUN CGO_ENABLED=0 go build -o macuahuitl .

# Stage 3: minimal runtime. The binary is fully static and makes no outbound
# calls, so no shell or CA certificates are needed. Render injects PORT.
FROM scratch
COPY --from=server /app/macuahuitl /macuahuitl
ENTRYPOINT ["/macuahuitl"]
