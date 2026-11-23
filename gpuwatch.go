package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	gpuFeedURL      = "https://carbonite.co.za/index.php?forums/amd_gpu/index.rss"
	gpuMaxPriceZAR  = 5000
	gpuDefaultEvery = 30 * time.Minute
	gpuMaxFeedBytes = 4 << 20
	ntfyBaseURL     = "https://ntfy.sh/"
)

type gpuFeed struct {
	Items []gpuItem `xml:"channel>item"`
}

type gpuItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
	Body    string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

type gpuListing struct {
	GUID     string
	Title    string
	Link     string
	Posted   string
	Price    int
	Location string
}

var (
	htmlBreak  = regexp.MustCompile(`(?i)<br\s*/?>`)
	htmlTag    = regexp.MustCompile(`<[^>]*>`)
	priceToken = regexp.MustCompile(`(?i)r?\s*(\d{1,3}(?:[ ,]\d{3})+|\d+)`)
	capacity8  = regexp.MustCompile(`(?i)\b8\s*(gb|gig)\b`)
	excluded   = regexp.MustCompile(`(?i)\b(wanted|wtb|sold|nvidia|geforce|gtx|rtx)\b`)
	capeArea   = regexp.MustCompile(`(?i)\b(cape\s?town|cpt|western\s+cape|bellville|durbanville|stellenbosch|somerset\s+west|strand|gordon'?s\s+bay|milnerton|table\s+view|blouberg|claremont|newlands|rondebosch|observatory|woodstock|sea\s+point|century\s+city|parow|goodwood|brackenfell|kuils\s+river|kraaifontein|fish\s+hoek|muizenberg|constantia|tokai|plumstead|wynberg|mitchells\s+plain|khayelitsha|paarl|hout\s+bay|bothasig|edgemead|plattekloof|tygervalley|northern\s+suburbs|southern\s+suburbs)\b`)
)

func bodyText(raw string) string {
	s := htmlBreak.ReplaceAllString(raw, "\n")
	s = htmlTag.ReplaceAllString(s, "")
	return html.UnescapeString(s)
}

func field(text, key string) string {
	prefix := strings.ToLower(key) + ":"
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if len(line) >= len(prefix) && strings.ToLower(line[:len(prefix)]) == prefix {
			return strings.TrimSpace(line[len(prefix):])
		}
	}
	return ""
}

func parsePrice(value string) int {
	value = strings.TrimSpace(value)
	if value == "" || (value[0] != 'R' && value[0] != 'r' && (value[0] < '0' || value[0] > '9')) {
		return -1
	}
	matches := priceToken.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return -1
	}
	digits := strings.NewReplacer(" ", "", ",", "").Replace(matches[len(matches)-1][1])
	n, err := strconv.Atoi(digits)
	if err != nil || n < 100 {
		return -1
	}
	return n
}

func matchListing(it gpuItem) (gpuListing, bool) {
	text := bodyText(it.Body)
	if excluded.MatchString(it.Title) || excluded.MatchString(field(text, "Item")) {
		return gpuListing{}, false
	}
	if !capacity8.MatchString(it.Title) && !capacity8.MatchString(field(text, "Item")) {
		return gpuListing{}, false
	}
	location := field(text, "Location")
	if !capeArea.MatchString(location) {
		return gpuListing{}, false
	}
	price := parsePrice(field(text, "Price"))
	if price > gpuMaxPriceZAR {
		return gpuListing{}, false
	}
	return gpuListing{
		GUID:     it.GUID,
		Title:    strings.TrimSpace(it.Title),
		Link:     strings.TrimSpace(it.Link),
		Posted:   strings.TrimSpace(it.PubDate),
		Price:    price,
		Location: location,
	}, true
}

func parseGPUFeed(r io.Reader) ([]gpuItem, error) {
	var feed gpuFeed
	if err := xml.NewDecoder(io.LimitReader(r, gpuMaxFeedBytes)).Decode(&feed); err != nil {
		return nil, fmt.Errorf("decode feed: %w", err)
	}
	return feed.Items, nil
}

func fetchGPUFeed(ctx context.Context, client *http.Client) ([]gpuItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gpuFeedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "macuahuitl-gpuwatch")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Debug("close gpu feed body", "err", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("feed status %d", resp.StatusCode)
	}
	return parseGPUFeed(resp.Body)
}

func headerSafe(s string) string {
	s = strings.Map(func(r rune) rune {
		if r < ' ' || r == 0x7f {
			return ' '
		}
		return r
	}, s)
	return mime.QEncoding.Encode("utf-8", strings.TrimSpace(s))
}

func notifyNTFY(ctx context.Context, client *http.Client, topic string, l gpuListing, price string) error {
	body := price + " · " + l.Location
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ntfyBaseURL+topic, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "macuahuitl-gpuwatch")
	req.Header.Set("Title", headerSafe(l.Title))
	req.Header.Set("Click", headerSafe(l.Link))
	req.Header.Set("Priority", "high")
	req.Header.Set("Tags", "video_camera")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Debug("close ntfy body", "err", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ntfy status %d", resp.StatusCode)
	}
	return nil
}

func watchGPUs(ctx context.Context, every time.Duration) {
	client := &http.Client{Timeout: 20 * time.Second}
	topic := os.Getenv("NTFY_TOPIC")
	notified := map[string]bool{}
	poll := func() {
		items, err := fetchGPUFeed(ctx, client)
		if err != nil {
			slog.Warn("gpu feed fetch failed", "err", err)
			return
		}
		fresh := 0
		for _, it := range items {
			l, ok := matchListing(it)
			if !ok {
				continue
			}
			price := "unknown"
			if l.Price >= 0 {
				price = "R" + strconv.Itoa(l.Price)
			}
			if _, known := notified[l.GUID]; !known {
				notified[l.GUID] = false
				fresh++
				slog.Warn("gpu listing",
					"event", "gpu_listing",
					"title", l.Title,
					"price", price,
					"location", l.Location,
					"link", l.Link,
					"posted", l.Posted,
				)
			}
			if topic != "" && !notified[l.GUID] {
				if err := notifyNTFY(ctx, client, topic, l, price); err != nil {
					slog.Warn("gpu notify failed", "err", err, "link", l.Link)
					continue
				}
				notified[l.GUID] = true
			}
		}
		slog.Info("gpu feed polled", "items", len(items), "new_matches", fresh)
	}
	poll()
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			poll()
		}
	}
}
