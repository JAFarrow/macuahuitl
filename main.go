// Command macuahuitl serves the prerendered SvelteKit site (HTML for humans)
// and the raw Obsidian vault markdown (for agents) from a single binary.
package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handleHealth)
	// Vault attachments (images etc.); the frontend rewrites Obsidian bare
	// image embeds to this prefix. StripPrefix + the sub-FS confines requests
	// to the attachments directory (no ../ escapes into the rest of the vault).
	mux.Handle("GET /attachments/", http.StripPrefix("/attachments/", http.FileServerFS(attachments)))
	mux.Handle("GET /posts/{slug}", sectionHandler(content, "posts", html))
	mux.Handle("GET /projects/{slug}", sectionHandler(content, "projects", html))
	mux.Handle("GET /", html)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server exited", "err", err)
		os.Exit(1)
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		slog.Debug("write health response", "err", err)
	}
}

// sectionHandler serves one vault section (posts or projects). A slug ending
// in .md, or an Accept header containing text/markdown, gets the raw markdown;
// anything else falls through to the prerendered HTML.
func sectionHandler(content fs.FS, section string, html http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		markdown := strings.Contains(r.Header.Get("Accept"), "text/markdown")
		if s, ok := strings.CutSuffix(slug, ".md"); ok {
			slug, markdown = s, true
		}
		if markdown {
			serveMarkdown(w, r, content, section, slug)
			return
		}
		html.ServeHTTP(w, r)
	}
}

func serveMarkdown(w http.ResponseWriter, r *http.Request, content fs.FS, section, slug string) {
	if slug == "" || strings.Contains(slug, "/") || strings.Contains(slug, "..") {
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(content, section+"/"+slug+".md")
	if err != nil || isDraft(data) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	if _, err := w.Write(data); err != nil {
		slog.Debug("write markdown response", "err", err)
	}
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
