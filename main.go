package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//go:embed all:frontend/build all:content
var embedded embed.FS

func main() {
	setupLogging()
	build, err := fs.Sub(embedded, "frontend/build")
	if err != nil {
		slog.Error("embed frontend/build", "err", err)
		os.Exit(1)
	}
	content, err := fs.Sub(embedded, "content")
	if err != nil {
		slog.Error("embed content", "err", err)
		os.Exit(1)
	}

	html := http.FileServerFS(build)

	project := slugMarkdown(content)
	projectsIdx := fixedMarkdown(build, "projects.md")
	about := fixedMarkdown(content, "about.md")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.Handle("GET /projects/{slug}", sectionHandler(project, html))
	mux.Handle("GET /projects/{slug}/{$}", onAccept(project, html))
	mux.HandleFunc("GET /projects.md", projectsIdx)
	mux.Handle("GET /projects/{$}", onAccept(projectsIdx, http.RedirectHandler("/", http.StatusMovedPermanently)))
	mux.HandleFunc("GET /about.md", about)
	mux.Handle("GET /{$}", onAccept(about, html))
	mux.Handle("GET /", html)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           accessLog(mux),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:      60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	shutdown := make(chan struct{})
	go func() {
		defer close(shutdown)
		<-ctx.Done()
		slog.Info("shutting down")
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(timeout); err != nil {
			slog.Error("server shutdown", "err", err)
		}
	}()

	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server exited", "err", err)
		os.Exit(1)
	}
	<-shutdown
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		slog.Debug("write health response", "err", err)
	}
}

func sectionHandler(md http.Handler, html http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.PathValue("slug"), ".md") || acceptsMarkdown(r) {
			md.ServeHTTP(w, r)
			return
		}
		html.ServeHTTP(w, r)
	}
}

func onAccept(md http.Handler, html http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if acceptsMarkdown(r) {
			md.ServeHTTP(w, r)
			return
		}
		html.ServeHTTP(w, r)
	}
}

func acceptsMarkdown(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/markdown")
}

func slugMarkdown(content fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if meta, ok := r.Context().Value(metaCtxKey{}).(*requestMeta); ok {
			meta.typ = "markdown"
		}
		slug := strings.TrimSuffix(r.PathValue("slug"), ".md")
		if slug == "" || strings.Contains(slug, "/") || strings.Contains(slug, "..") {
			http.NotFound(w, r)
			return
		}
		serveMarkdownFile(w, r, content, "projects/"+slug+".md")
	}
}

func fixedMarkdown(fsys fs.FS, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveMarkdownFile(w, r, fsys, path)
	}
}

func serveMarkdownFile(w http.ResponseWriter, r *http.Request, fsys fs.FS, path string) {
	if meta, ok := r.Context().Value(metaCtxKey{}).(*requestMeta); ok {
		meta.typ = "markdown"
	}
	data, err := fs.ReadFile(fsys, path)
	if err != nil || isDraft(data) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		slog.Debug("write markdown response", "err", err)
	}
}

func isDraft(data []byte) bool {
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return false
	}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "---" {
			return false
		}
		if value, ok := strings.CutPrefix(line, "draft:"); ok && strings.TrimSpace(value) == "true" {
			return true
		}
	}
	return false
}
