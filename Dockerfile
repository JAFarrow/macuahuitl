# syntax=docker/dockerfile:1

FROM node:26-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json frontend/
RUN cd frontend && npm ci
COPY frontend/ frontend/
COPY content/ content/
RUN cd frontend && npm run build

FROM golang:1.26-alpine AS server
WORKDIR /app
COPY go.mod main.go logging.go ./
COPY content/ content/
COPY --from=frontend /app/frontend/build frontend/build
RUN CGO_ENABLED=0 go build -o macuahuitl .

FROM scratch
COPY --from=server /app/macuahuitl /macuahuitl
ENTRYPOINT ["/macuahuitl"]
