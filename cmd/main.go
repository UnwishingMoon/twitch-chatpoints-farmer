package main

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"twitch_chatpoints_farmer/pkg/app"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func main() {
	appl := app.NewApp()
	bOpen := false
	var ctx context.Context
	var cancel context.CancelFunc
	var res []byte

	appl.SetAppToken()
	appl.StartBrowser()

	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), appl.Browser.URL)
	defer cancel()

	for {
		if appl.GetStreams() {
			log.Println("[INFO] Stream is online.")

			if !bOpen {
				log.Println("[INFO] Opening browser tab..")
				bOpen = true

				ctx, cancel = chromedp.NewContext(allocatorContext)
				defer cancel()

				chromedp.Run(ctx,
					chromedp.Navigate(app.SiteURL+appl.Configuration.UserName),
					setCookie("www.twitch.tv", "twilight-user", appl.Configuration.AuthCookie),
					writeLocalStorage("video-muted", appl.Settings.VideoMuted),
					writeLocalStorage("video-quality", appl.Settings.VideoQuality),
					chromedp.Reload(),
				)

				time.Sleep(2 * time.Second)
			}

			log.Println("[INFO] Clicking the button..")
			chromedp.Run(ctx,
				chromedp.EvaluateAsDevTools("document.querySelector('.community-points-summary div.tw-full-height button.tw-button').click()", &res),
			)
		} else {
			log.Println("[INFO] Stream is offline.")

			if bOpen {
				log.Println("[INFO] Closing the browser")
				bOpen = false

				cancel() // Forcefully closing the browser..
			}
		}

		log.Println("[INFO] Waiting before next click..")
		time.Sleep(5 * time.Minute)
	}

}

func setCookie(host string, name string, value string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			exp := cdp.TimeSinceEpoch(time.Now().Add(90 * 24 * time.Hour))

			_, err := network.SetCookie(name, value).
				WithExpires(&exp).
				WithDomain(host).
				WithHTTPOnly(false).
				Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func writeLocalStorage(name string, value map[string]string) chromedp.Tasks {
	res, _ := json.Marshal(value)

	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate("window.localStorage.setItem('" + name + "', '" + string(res) + "')").Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
	}
}
