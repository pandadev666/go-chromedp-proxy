package main

import (
	"context"
    "log"
    "os"
	"fmt"
    "github.com/chromedp/chromedp"
    "github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("http://zproxy.lum-superproxy.io:22225"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	lctx, lcancel := context.WithCancel(ctx)
	chromedp.ListenTarget(lctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *fetch.EventRequestPaused:
			go func() {
				_ = chromedp.Run(ctx, fetch.ContinueRequest(ev.RequestID))
			}()
		case *fetch.EventAuthRequired:
			if ev.AuthChallenge.Source == fetch.AuthChallengeSourceProxy {
				go func() {
					_ = chromedp.Run(ctx,
						fetch.ContinueWithAuth(ev.RequestID, &fetch.AuthChallengeResponse{
							Response: fetch.AuthChallengeResponseResponseProvideCredentials,
							Username: "brd-customer-hl_90faa028-zone-data_center-country-gb",
							Password: "i7t589h750mh",
						}),
						fetch.Disable(),
					)
					lcancel()
				}()
			}
		}
	})

	if err := chromedp.Run(ctx,
		// fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate("https://direct.asda.com/george/outdoor-garden/all-garden-sheds-outdoor-storage/forest-garden-pressure-treated-pent-garden-store/051326489,default,pd.html?cgid=D33M13G01C01"),
	); err != nil {
		log.Fatal(err)
	}

	var data string
	headers := map[string]interface{}{
		"Sec-Ch-Ua": `"Google Chrome";v="113", "Chromium";v="113", "Not-A.Brand";v="24"`,
	}

	urls := []string{
		"https://direct.asda.com/george/outdoor-garden/all-garden-sheds-outdoor-storage/forest-garden-pressure-treated-pent-garden-store/051326489,default,pd.html?cgid=D33M13G01C01",
		"https://direct.asda.com/george/outdoor-garden/all-outdoor-buildings-storage/keter-brightwood-454l-storage-box-brown/051332779,default,pd.html?cgid=D33M13G01C01",
		"https://direct.asda.com/george/outdoor-garden/all-garden-sheds-outdoor-storage/forest-garden-pressure-treated-pent-tall-garden-store/051326491,default,pd.html?cgid=D33M13G01C01",
		"https://direct.asda.com/george/home/home-storage/wham-box-and-lid-96-litres/051172411,default,pd.html?cgid=D26M21G01C03",
		"https://direct.asda.com/george/outdoor-garden/all-garden-sheds-outdoor-storage/forest-garden-pressure-treated-tall-apex-garden-store/051326490,default,pd.html?cgid=D33M13G01C01",
		"https://direct.asda.com/george/home/curtains-blinds/blinds/aluminium-venetian-blind-white/GEM793317,default,pd.html?cgid=D26M01G01C04",
		"https://direct.asda.com/george/home/curtains-blinds/blinds/aluminium-venetian-blind-silver/GEM793251,default,pd.html?cgid=D26M01G01C04",
		"https://direct.asda.com/george/home/curtains-blinds/homemaker-white-25mm-wood-blind/GEM997981,default,pd.html?cgid=D26M01G01C04",
		"https://direct.asda.com/george/home/curtains-blinds/homemaker-grey-25mm-wood-blind/GEM997747,default,pd.html?cgid=D26M01G01C04",
		"https://direct.asda.com/george/home/vacuums-steam-mops/tower-t513004blg-vl40-pro-pet-3-in-1-cordless-vacuum-cleaner-with-cyclonic-suction/051269325,default,pd.html?cgid=D26M10G10C09",
		"https://direct.asda.com/george/home/vacuums-steam-mops/tower-t548003-tcw-aquajet-plus-carpet-washer/051269327,default,pd.html?cgid=D26M10G10C09",
		"https://direct.asda.com/george/home/vacuums-steam-mops/numatic-henry-hvr160-cylinder-vacuum-cleaner/050483938,default,pd.html?cgid=D26M10G10C09",
		"https://direct.asda.com/george/home/vacuums-steam-mops/hetty-het160-vacuum-cleaner/050367836,default,pd.html?cgid=D26M10G10C09",
		"https://direct.asda.com/george/home/vacuums-steam-mops/tower-hh77-cordless-74v-handheld-vacuum/051196544,default,pd.html?cgid=D26M10G10C09",
		"https://direct.asda.com/george/home/vacuums-steam-mops/tower-txp10pet-cylinder-vacuum/051196531,default,pd.html?cgid=D26M10G10C09",
	}

	tctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	for idx, url := range urls {
		if err := chromedp.Run(tctx,
			network.Enable(),
			network.SetExtraHTTPHeaders(network.Headers(headers)),
			chromedp.Navigate(url),
			chromedp.OuterHTML ("html", &data, chromedp.ByQuery),
		); err != nil {
			log.Fatal(err)
		}
	
		file, _ := os.Create(fmt.Sprintf("%v.html", idx))
		file.WriteString(data)
		defer file.Close()

		fmt.Printf("%v scraped!", idx)
	}
}
