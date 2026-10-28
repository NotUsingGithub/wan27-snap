package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://wan2-7.io"
	WebsiteURL     = "https://wan2-7.io/"
	MCPURL         = "https://wan2-7.io/mcp/"
	MCPGuideURL    = "https://wan2-7.io/mcp/guide/"
	CLIURL         = "https://wan2-7.io/cli/"
	PricingURL     = "https://wan2-7.io/pricing/"
	Version        = "0.1.0"
)

var validIDRe = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(token, baseURL string) *Client {
	if token == "" {
		token = os.Getenv("WAN27_API_TOKEN")
	}
	if baseURL == "" {
		baseURL = os.Getenv("WAN27_BASE_URL")
	}
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Token:      strings.TrimSpace(token),
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func ValidateID(value, field string) error {
	if value == "" || !validIDRe.MatchString(value) {
		return fmt.Errorf("invalid %s: %q", field, value)
	}
	return nil
}

func (c *Client) Request(method, path string, query map[string]string, body any) (map[string]any, error) {
	if c.Token == "" {
		return nil, errors.New("missing WAN27_API_TOKEN; obtain an API token at https://wan2-7.io/mcp/guide/")
	}
	u := c.BaseURL + path
	if len(query) > 0 {
		vals := url.Values{}
		for k, v := range query {
			if v != "" {
				vals.Set(k, v)
			}
		}
		if enc := vals.Encode(); enc != "" {
			u += "?" + enc
		}
	}
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, u, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "wan27-snap/"+Version+" (+https://wan2-7.io/cli/)")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var parsed map[string]any
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &parsed)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parsed, fmt.Errorf("wan27 API HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	return parsed, nil
}

func main() {
	var (
		model       = flag.String("model", "wan2.7-t2v", "Model ID (wan2.7-t2v, wan2.7-i2v, wan2.7-image, wan2.7-image-pro)")
		resolution  = flag.String("resolution", "1080p", "Output resolution (720p, 1080p, 2K)")
		aspectRatio = flag.String("aspect-ratio", "16:9", "Aspect ratio (16:9, 9:16, 1:1)")
		duration    = flag.Int("duration", 5, "Video duration in seconds (5 or 10)")
		imageURL    = flag.String("image", "", "Input image URL for image-to-video")
		showCredits = flag.Bool("credits", false, "Check remaining Wan 2.7 account credits")
		estimate    = flag.Bool("estimate", false, "Estimate credit cost without starting generation")
		showURLs    = flag.Bool("urls", false, "Print official Wan 2.7 documentation and MCP URLs")
		showVersion = flag.Bool("version", false, "Print CLI version")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("wan27 %s (%s)\n", Version, CLIURL)
		return
	}
	if *showURLs {
		out, _ := json.MarshalIndent(map[string]string{
			"website":   WebsiteURL,
			"mcp":       MCPURL,
			"mcp_guide": MCPGuideURL,
			"cli":       CLIURL,
			"pricing":   PricingURL,
		}, "", "  ")
		fmt.Println(string(out))
		return
	}

	client := NewClient("", "")
	if *showCredits {
		res, err := client.Request("GET", "/api/v1/account/balance", nil, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(out))
		return
	}

	prompt := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if prompt == "" {
		fmt.Fprintf(os.Stderr, "Usage: wan27 [flags] \"<prompt>\"\nDocs: %s | MCP: %s | Pricing: %s\n", CLIURL, MCPURL, PricingURL)
		os.Exit(2)
	}

	payload := map[string]any{
		"model":        *model,
		"prompt":       prompt,
		"resolution":   *resolution,
		"aspect_ratio": *aspectRatio,
		"duration":     *duration,
	}
	if *imageURL != "" {
		payload["image_url"] = *imageURL
	}

	endpoint := "/api/v1/generations"
	if *estimate {
		endpoint = "/api/v1/credits/estimate"
	}
	res, err := client.Request("POST", endpoint, nil, payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	out, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(out))
}
