package screenshot

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/chromedp/chromedp"
)

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Use
// device.Reset to reset the emulation and viewport settings.
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}

func Screenshot(url string, quality int) string {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte

	if err := chromedp.Run(ctx, fullScreenshot(url, quality, &buf)); err != nil {
		log.Fatal(err)
	}

	// save to file within screenshots folder, create if exists, with fqdn as filename
	// extract domain from url
	domain := strings.Split(url, "/")[2]
	filename := "screenshots/" + domain + ".png"

	// check if screenshot folder exists, if not create it
	if _, err := os.Stat("screenshot"); os.IsNotExist(err) {
		os.Mkdir("screenshot", 0755)
	}

	if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
		log.Fatal(err)
	}

	return filename
}
