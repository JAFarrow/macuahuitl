package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func setupLogging() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 {
				switch a.Key {
				case slog.TimeKey:
					a.Key = "dt"
				case slog.MessageKey:
					a.Key = "message"
				}
			}
			return a
		},
	})))
}

type requestMeta struct {
	typ string
}

type metaCtxKey struct{}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(p)
	r.bytes += n
	return n, err
}

func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}
		meta := &requestMeta{}
		r = r.WithContext(context.WithValue(r.Context(), metaCtxKey{}, meta))
		rec := &statusRecorder{ResponseWriter: w}
		start := time.Now()
		defer func() {
			panicked := recover()
			dur := time.Since(start)

			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}
			if panicked != nil && rec.status == 0 {
				status = http.StatusInternalServerError
			}
			typ := meta.typ
			if typ == "" {
				typ = "html"
			}
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				ip = strings.TrimSpace(strings.Split(fwd, ",")[0])
			}
			level := slog.LevelInfo
			if status >= http.StatusInternalServerError {
				level = slog.LevelError
			}
			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("route", r.Pattern),
				slog.String("host", r.Host),
				slog.String("path", r.URL.Path),
			}
			if r.URL.RawQuery != "" {
				attrs = append(attrs, slog.String("query", r.URL.RawQuery))
			}
			attrs = append(attrs,
				slog.Int("status", status),
				slog.Float64("duration_ms", float64(dur.Microseconds())/1e3),
				slog.Int("bytes", rec.bytes),
				slog.String("type", typ),
				slog.String("content_type", rec.Header().Get("Content-Type")),
				slog.String("user_agent", r.UserAgent()),
				slog.String("referer", r.Referer()),
				slog.String("ip", ip),
			)
			slog.LogAttrs(r.Context(), level, fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, status), attrs...)
			if panicked != nil {
				panic(panicked)
			}
		}()
		next.ServeHTTP(rec, r)
	})
}
