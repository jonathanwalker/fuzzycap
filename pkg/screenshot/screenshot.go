package screenshot

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chromedp/chromedp"
)

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

// Screenshot captures a full-page screenshot of the given URL at the specified quality.
// It returns the filename where the screenshot was saved.
func Screenshot(url string, quality int) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, fullScreenshot(url, quality, &buf)); err != nil {
		return "", fmt.Errorf("chromedp run failed for %s: %w", url, err)
	}

	domain := extractDomain(url)
	filename := "screenshots/" + domain + ".png"

	if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
		return "", fmt.Errorf("failed to write screenshot: %w", err)
	}

	return filename, nil
}

func extractDomain(url string) string {
	parts := strings.Split(url, "//")
	if len(parts) < 2 {
		return "unknown"
	}
	domainAndPath := parts[1]
	domain := strings.Split(domainAndPath, "/")[0]
	domain = strings.TrimSpace(domain)
	if domain == "" {
		domain = "unknown"
	}
	// Replace any colons, etc., that might be in domain (for example, ports)
	domain = strings.ReplaceAll(domain, ":", "_")
	return domain
}
