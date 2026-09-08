// Access logging: one structured slog line per request locally, plus a
// hand-rolled OTLP/HTTP logs exporter (stdlib only) that ships the same
// records to Grafana Cloud when the standard OTEL_EXPORTER_OTLP_* env vars
// are set. Local slog output always happens; Render's log retention is the
// fallback when export is disabled or a batch is dropped.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	otlpFlushInterval = 5 * time.Second
	otlpBatchSize     = 64
	otlpQueueSize     = 256
)

// logField is one typed access-log attribute, rendered identically into the
// local slog line and the OTLP record attributes.
type logField struct {
	key string
	val any // string or int64
}

// requestMeta carries facts from the handler that knows them back to the
// access-log middleware: request contexts only flow downward, so the
// middleware plants a pointer that handlers mutate.
type requestMeta struct {
	typ string // "markdown" | "html"
}

type metaCtxKey struct{}

// statusRecorder captures the response status and byte count for logging.
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

// Unwrap lets http.ResponseController reach the underlying ResponseWriter.
func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// accessLog emits one structured log line per request. /api/health is
// skipped: Render polls it constantly, so it is pure noise in both local
// logs and Grafana. A nil exp means local logging only.
func accessLog(exp *otlpExporter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/health" {
			next.ServeHTTP(w, r)
			return
		}
		meta := &requestMeta{}
		r = r.WithContext(context.WithValue(r.Context(), metaCtxKey{}, meta))
		rec := &statusRecorder{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(rec, r)
		dur := time.Since(start)

		status := rec.status
		if status == 0 {
			status = http.StatusOK
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
		fields := []logField{
			{"method", r.Method},
			{"route", r.Pattern},
			{"path", r.URL.Path},
			{"status", int64(status)},
			{"duration_ms", dur.Milliseconds()},
			{"bytes", int64(rec.bytes)},
			{"type", typ},
			{"user_agent", r.UserAgent()},
			{"referer", r.Referer()},
			{"ip", ip},
		}
		level := slog.LevelInfo
		if status >= http.StatusInternalServerError {
			level = slog.LevelError
		}
		body := fmt.Sprintf("%s %d %dms", r.Pattern, status, dur.Milliseconds())
		attrs := make([]slog.Attr, len(fields))
		for i, f := range fields {
			attrs[i] = slog.Any(f.key, f.val)
		}
		slog.LogAttrs(r.Context(), level, body, attrs...)
		if exp != nil {
			exp.enqueue(start, level, body, fields)
		}
	})
}

// OTLP/JSON payload types. int64s (timeUnixNano, intValue) serialize as
// strings per the protobuf JSON mapping.
type otlpValue struct {
	Str *string `json:"stringValue,omitempty"`
	Int *string `json:"intValue,omitempty"`
}

type otlpAttr struct {
	Key   string    `json:"key"`
	Value otlpValue `json:"value"`
}

type otlpLogRecord struct {
	TimeUnixNano   string     `json:"timeUnixNano"`
	SeverityNumber int        `json:"severityNumber"`
	SeverityText   string     `json:"severityText"`
	Body           otlpValue  `json:"body"`
	Attributes     []otlpAttr `json:"attributes"`
}

type otlpResource struct {
	Attributes []otlpAttr `json:"attributes"`
}

type otlpScope struct {
	Name string `json:"name"`
}

type otlpScopeLogs struct {
	Scope      otlpScope       `json:"scope"`
	LogRecords []otlpLogRecord `json:"logRecords"`
}

type otlpResourceLogs struct {
	Resource  otlpResource    `json:"resource"`
	ScopeLogs []otlpScopeLogs `json:"scopeLogs"`
}

type otlpPayload struct {
	ResourceLogs []otlpResourceLogs `json:"resourceLogs"`
}

type otlpExporter struct {
	url           string
	headers       http.Header
	client        *http.Client
	resourceAttrs []otlpAttr
	queue         chan otlpLogRecord
	done          chan struct{}
	wg            sync.WaitGroup
}

// newOTLPExporterFromEnv builds an exporter from the standard OTel env vars
// (as issued by Grafana Cloud), or returns nil when no endpoint is set:
//
//	OTEL_EXPORTER_OTLP_ENDPOINT  base URL; records POST to <endpoint>/v1/logs
//	OTEL_EXPORTER_OTLP_HEADERS   "k=v,k=v" (e.g. Authorization=Basic ...)
//	OTEL_SERVICE_NAME            resource service.name (default macuahuitl)
func newOTLPExporterFromEnv() *otlpExporter {
	endpoint := strings.TrimRight(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "/")
	if endpoint == "" {
		return nil
	}
	headers := http.Header{}
	for _, kv := range strings.Split(os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"), ",") {
		k, v, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		// PathUnescape, not QueryUnescape: QueryUnescape decodes '+' as a
		// space, which would corrupt the base64 Basic token Grafana issues.
		if dec, err := url.PathUnescape(strings.TrimSpace(v)); err == nil {
			v = dec
		}
		headers.Set(strings.TrimSpace(k), v)
	}
	service := os.Getenv("OTEL_SERVICE_NAME")
	if service == "" {
		service = "macuahuitl"
	}
	e := &otlpExporter{
		url:     endpoint + "/v1/logs",
		headers: headers,
		client:  &http.Client{Timeout: 10 * time.Second},
		resourceAttrs: []otlpAttr{
			strAttr("service.name", service),
			strAttr("telemetry.sdk.name", "macuahuitl"),
			strAttr("telemetry.sdk.language", "go"),
		},
		queue: make(chan otlpLogRecord, otlpQueueSize),
		done:  make(chan struct{}),
	}
	e.wg.Add(1)
	go e.run()
	slog.Info("otlp log export enabled", "endpoint", endpoint, "service", service)
	return e
}

func strAttr(key, val string) otlpAttr { return otlpAttr{key, otlpValue{Str: &val}} }

func fieldAttr(f logField) otlpAttr {
	if v, ok := f.val.(int64); ok {
		s := strconv.FormatInt(v, 10)
		return otlpAttr{f.key, otlpValue{Int: &s}}
	}
	s, _ := f.val.(string)
	return otlpAttr{f.key, otlpValue{Str: &s}}
}

// enqueue buffers one access-log entry for export. A full queue drops the
// record with a local warning; the local slog line already persisted it.
func (e *otlpExporter) enqueue(ts time.Time, level slog.Level, body string, fields []logField) {
	sev, sevText := 9, "INFO" // OTLP severity numbers: 9=INFO, 17=ERROR
	if level >= slog.LevelError {
		sev, sevText = 17, "ERROR"
	}
	rec := otlpLogRecord{
		TimeUnixNano:   strconv.FormatInt(ts.UnixNano(), 10),
		SeverityNumber: sev,
		SeverityText:   sevText,
		Body:           otlpValue{Str: &body},
		Attributes:     make([]otlpAttr, len(fields)),
	}
	for i, f := range fields {
		rec.Attributes[i] = fieldAttr(f)
	}
	select {
	case e.queue <- rec:
	default:
		slog.Warn("otlp queue full, dropping record")
	}
}

func (e *otlpExporter) run() {
	defer e.wg.Done()
	ticker := time.NewTicker(otlpFlushInterval)
	defer ticker.Stop()
	var batch []otlpLogRecord
	flush := func() {
		if len(batch) > 0 {
			e.export(batch)
			batch = batch[:0]
		}
	}
	for {
		select {
		case rec := <-e.queue:
			if batch = append(batch, rec); len(batch) >= otlpBatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-e.done:
			// Drain whatever is still queued, then one final flush.
			for {
				select {
				case rec := <-e.queue:
					batch = append(batch, rec)
				default:
					flush()
					return
				}
			}
		}
	}
}

// export POSTs one batch. Failure is log-and-drop: this is access telemetry
// and the same lines exist in local logs.
func (e *otlpExporter) export(batch []otlpLogRecord) {
	data, err := json.Marshal(otlpPayload{ResourceLogs: []otlpResourceLogs{{
		Resource:  otlpResource{Attributes: e.resourceAttrs},
		ScopeLogs: []otlpScopeLogs{{Scope: otlpScope{Name: "macuahuitl/access"}, LogRecords: batch}},
	}}})
	if err != nil {
		slog.Error("otlp marshal", "err", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, e.url, bytes.NewReader(data))
	if err != nil {
		slog.Error("otlp request", "err", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k := range e.headers {
		req.Header.Set(k, e.headers.Get(k))
	}
	resp, err := e.client.Do(req)
	if err != nil {
		slog.Warn("otlp export failed, dropping batch", "err", err, "records", len(batch))
		return
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		slog.Warn("otlp export rejected, dropping batch", "status", resp.StatusCode, "records", len(batch))
	}
}

// shutdown flushes any buffered records; bounded by the client timeout.
func (e *otlpExporter) shutdown() {
	close(e.done)
	e.wg.Wait()
}
