// Command macuahuitl serves the prerendered SvelteKit site (HTML for humans)
// and the raw Obsidian vault markdown (for agents) from a single binary.
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
	attachments, err := fs.Sub(content, "attachments")
	if err != nil {
		slog.Error("embed content/attachments", "err", err)
		os.Exit(1)
	}

	html := http.FileServerFS(build)

	// Markdown sources: vault documents by section slug, and fixed files
	// (build-generated index twins plus the vault about page).
	project := slugMarkdown(content, "projects")
	cvEntry := slugMarkdown(content, "cv")
	projectsIdx := fixedMarkdown(build, "projects.md")
	cvIdx := fixedMarkdown(build, "cv.md")
	about := fixedMarkdown(content, "about.md")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	// Vault attachments (images etc.); the frontend rewrites Obsidian bare
	// image embeds to this prefix. StripPrefix + the sub-FS confines requests
	// to the attachments directory (no ../ escapes into the rest of the vault),
	// and directory listings are 404ed rather than rendered.
	mux.Handle("GET /attachments/", noListings(http.StripPrefix("/attachments/", http.FileServerFS(attachments))))
	mux.Handle("GET /projects/{slug}", sectionHandler(project, html))
	mux.Handle("GET /projects/{slug}/{$}", onAccept(project, html))
	mux.Handle("GET /cv/{slug}", sectionHandler(cvEntry, html))
	mux.Handle("GET /cv/{slug}/{$}", onAccept(cvEntry, html))
	mux.HandleFunc("GET /projects.md", projectsIdx)
	mux.HandleFunc("GET /cv.md", cvIdx)
	mux.Handle("GET /projects/{$}", onAccept(projectsIdx, html))
	mux.Handle("GET /cv/{$}", onAccept(cvIdx, html))
	mux.HandleFunc("GET /about.md", about)
	mux.Handle("GET /{$}", onAccept(about, html))
	mux.Handle("GET /", html)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	exporter := newOTLPExporterFromEnv()
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           accessLog(exporter, mux),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:      60 * time.Second,
	}

	// Render sends SIGTERM on deploy; shut down gracefully so buffered log
	// records get flushed rather than dropped.
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
		if exporter != nil {
			exporter.shutdown()
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

// sectionHandler serves a vault section at slash-less URLs: a slug ending in
// .md, or an Accept header containing text/markdown, gets the raw markdown;
// anything else falls through to the prerendered HTML.
func sectionHandler(md http.Handler, html http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.PathValue("slug"), ".md") || acceptsMarkdown(r) {
			md.ServeHTTP(w, r)
			return
		}
		html.ServeHTTP(w, r)
	}
}

// onAccept serves md when the request's Accept header contains text/markdown,
// otherwise html.
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

// slugMarkdown returns a handler that serves one vault document:
// content/<section>/<slug>.md. It validates the slug (defense in depth against
// traversal) and runs the draft check.
func slugMarkdown(content fs.FS, section string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Mark markdown intent for the access log even when the outcome is a
		// 404 (draft, missing, or invalid slug); the status field carries the outcome.
		if meta, ok := r.Context().Value(metaCtxKey{}).(*requestMeta); ok {
			meta.typ = "markdown"
		}
		slug := strings.TrimSuffix(r.PathValue("slug"), ".md")
		if slug == "" || strings.Contains(slug, "/") || strings.Contains(slug, "..") {
			http.NotFound(w, r)
			return
		}
		serveMarkdownFile(w, r, content, section+"/"+slug+".md")
	}
}

// fixedMarkdown returns a handler that serves one markdown file from fsys.
// The file is draft-checked too, but that is a no-op for frontmatter-less
// generated files and keeps behavior uniform when fsys points back at the vault.
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

// noListings 404s directory requests instead of rendering the stdlib file
// server's directory listing; individual files pass through unchanged.
func noListings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// isDraft reports whether the document's YAML frontmatter sets draft: true.
// It is a line-prefix scan of the frontmatter block only, not a YAML parser.
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
