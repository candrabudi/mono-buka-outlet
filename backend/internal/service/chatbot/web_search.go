package chatbot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ──────────────────────────────────────────────────────────────
// WEB SEARCH SERVICE
// Searches the web for franchise & business-related information
// Uses DuckDuckGo Instant Answer API (free, no API key needed)
// + Google Custom Search JSON API as fallback (optional)
// ──────────────────────────────────────────────────────────────

type WebSearchResult struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

type WebSearcher struct {
	httpClient *http.Client
}

func NewWebSearcher() *WebSearcher {
	return &WebSearcher{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Search performs a franchise/business-focused web search
// It combines multiple search strategies for best results
func (ws *WebSearcher) Search(query string) ([]WebSearchResult, error) {
	// Ensure query is franchise/business related by appending context
	enrichedQuery := ws.enrichQuery(query)

	var allResults []WebSearchResult

	// Strategy 1: DuckDuckGo Instant Answer API
	ddgResults, err := ws.searchDuckDuckGo(enrichedQuery)
	if err == nil {
		allResults = append(allResults, ddgResults...)
	}

	// Strategy 2: DuckDuckGo HTML search (scrape-lite)
	if len(allResults) < 3 {
		ddgHTMLResults, err := ws.searchDuckDuckGoHTML(enrichedQuery)
		if err == nil {
			allResults = append(allResults, ddgHTMLResults...)
		}
	}

	// Deduplicate and limit
	allResults = ws.deduplicateResults(allResults)
	if len(allResults) > 8 {
		allResults = allResults[:8]
	}

	return allResults, nil
}

// enrichQuery adds franchise/business context to search queries
func (ws *WebSearcher) enrichQuery(query string) string {
	lower := strings.ToLower(query)

	// If query already has business context, use as-is
	businessTerms := []string{"franchise", "bisnis", "outlet", "kemitraan", "usaha", "investasi", "modal", "profit", "waralaba"}
	for _, term := range businessTerms {
		if strings.Contains(lower, term) {
			return query
		}
	}

	// Add franchise context
	return query + " franchise bisnis Indonesia"
}

// searchDuckDuckGo uses DuckDuckGo Instant Answer API
func (ws *WebSearcher) searchDuckDuckGo(query string) ([]WebSearchResult, error) {
	apiURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1",
		url.QueryEscape(query))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "BukaOutlet-AI-Konsultan/1.0")

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ddgResp struct {
		Abstract       string `json:"Abstract"`
		AbstractSource string `json:"AbstractSource"`
		AbstractURL    string `json:"AbstractURL"`
		RelatedTopics  []struct {
			Text     string `json:"Text"`
			FirstURL string `json:"FirstURL"`
			Result   string `json:"Result"`
			Topics   []struct {
				Text     string `json:"Text"`
				FirstURL string `json:"FirstURL"`
			} `json:"Topics"`
		} `json:"RelatedTopics"`
		Results []struct {
			Text     string `json:"Text"`
			FirstURL string `json:"FirstURL"`
		} `json:"Results"`
	}

	if err := json.Unmarshal(body, &ddgResp); err != nil {
		return nil, err
	}

	var results []WebSearchResult

	// Abstract (primary answer)
	if ddgResp.Abstract != "" {
		results = append(results, WebSearchResult{
			Title:   ddgResp.AbstractSource,
			Snippet: ddgResp.Abstract,
			URL:     ddgResp.AbstractURL,
		})
	}

	// Related Topics
	for _, topic := range ddgResp.RelatedTopics {
		if topic.Text != "" && topic.FirstURL != "" {
			results = append(results, WebSearchResult{
				Title:   extractTitle(topic.Text),
				Snippet: topic.Text,
				URL:     topic.FirstURL,
			})
		}
		// Sub-topics
		for _, sub := range topic.Topics {
			if sub.Text != "" && sub.FirstURL != "" {
				results = append(results, WebSearchResult{
					Title:   extractTitle(sub.Text),
					Snippet: sub.Text,
					URL:     sub.FirstURL,
				})
			}
		}
	}

	// Direct results
	for _, r := range ddgResp.Results {
		if r.Text != "" {
			results = append(results, WebSearchResult{
				Title:   extractTitle(r.Text),
				Snippet: r.Text,
				URL:     r.FirstURL,
			})
		}
	}

	return results, nil
}

// searchDuckDuckGoHTML uses DuckDuckGo lite/HTML endpoint for more results
func (ws *WebSearcher) searchDuckDuckGoHTML(query string) ([]WebSearchResult, error) {
	apiURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; BukaOutlet-AI/1.0)")

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlStr := string(body)
	var results []WebSearchResult

	// Simple parsing of DuckDuckGo HTML results
	// Look for result snippets between known markers
	parts := strings.Split(htmlStr, "result__snippet")
	for i, part := range parts {
		if i == 0 {
			continue // skip before first result
		}
		if i > 6 {
			break // limit results
		}

		snippet := extractBetween(part, ">", "</")
		snippet = stripHTML(snippet)
		snippet = strings.TrimSpace(snippet)

		// Extract URL
		urlPart := ""
		if idx := strings.Index(part, "result__url"); idx != -1 {
			urlPart = extractBetween(part[idx:], "href=\"", "\"")
		}

		// Extract title
		titlePart := ""
		if idx := strings.Index(part, "result__a"); idx != -1 {
			titlePart = extractBetween(part[idx:], ">", "</")
			titlePart = stripHTML(titlePart)
		}

		if snippet != "" {
			results = append(results, WebSearchResult{
				Title:   titlePart,
				Snippet: snippet,
				URL:     urlPart,
			})
		}
	}

	return results, nil
}

// deduplicateResults removes duplicate results based on URL
func (ws *WebSearcher) deduplicateResults(results []WebSearchResult) []WebSearchResult {
	seen := map[string]bool{}
	var unique []WebSearchResult
	for _, r := range results {
		key := r.URL
		if key == "" {
			key = r.Snippet
		}
		if !seen[key] && r.Snippet != "" {
			seen[key] = true
			unique = append(unique, r)
		}
	}
	return unique
}

// FormatForContext formats search results as context for GPT
func FormatSearchResults(results []WebSearchResult) string {
	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## HASIL PENCARIAN WEB (Referensi Tambahan)\n\n")
	sb.WriteString("Berikut informasi terbaru yang ditemukan dari internet:\n\n")

	for i, r := range results {
		sb.WriteString(fmt.Sprintf("**[%d] %s**\n", i+1, r.Title))
		sb.WriteString(fmt.Sprintf("%s\n", r.Snippet))
		if r.URL != "" {
			sb.WriteString(fmt.Sprintf("Sumber: %s\n", r.URL))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n")
	sb.WriteString("*Gunakan informasi di atas sebagai referensi tambahan. Selalu pastikan informasi sesuai konteks BukaOutlet.*\n")

	return sb.String()
}

// ──────────────────────────────────────────────────────────────
// HELPERS
// ──────────────────────────────────────────────────────────────

func extractTitle(text string) string {
	// Take first sentence or first 80 chars
	if idx := strings.Index(text, " - "); idx > 0 && idx < 80 {
		return text[:idx]
	}
	if len(text) > 80 {
		return text[:80] + "..."
	}
	return text
}

func extractBetween(s, start, end string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return ""
	}
	s = s[i+len(start):]
	j := strings.Index(s, end)
	if j < 0 {
		return s
	}
	return s[:j]
}

func stripHTML(s string) string {
	var result strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}
