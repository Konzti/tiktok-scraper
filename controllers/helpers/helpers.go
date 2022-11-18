package helpers

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tiktok-scraper/constants"
	"tiktok-scraper/utils"
	"time"
)

func ScrapeForVideo(parsedUrl *url.URL, c *gin.Context) {
	if !strings.Contains(parsedUrl.Path, "/video/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid video path"})
		return
	}
	videoId := strings.Split(parsedUrl.Path, "/")[3]
	fmt.Println("Video ID: ", videoId)

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	var nodes []*cdp.Node
	var scrapeResponse interface{}
	selector := "video"
	start := time.Now()

	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("Mozilla/5.0  (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) coc_coc_browser/86.0.170 Chrome/80.0.3987.170 Safari/537.36"),
		chromedp.Navigate(parsedUrl.String()),
		chromedp.WaitReady(selector, chromedp.ByQuery),
		chromedp.Nodes(selector, &nodes),
		chromedp.ActionFunc(
			func(ctx context.Context) error {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					panic(err)
				}
				res, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				if err != nil {
					panic(err)
				}
				videoStrings := utils.FindVideoStrings(res)
				videoURL := utils.FindVideo(res)
				go storeFile(videoURL, videoId)
				scrapeResponse = map[string]interface{}{
					"avatar":   videoStrings[1],
					"img":      videoStrings[0],
					"videoURL": videoURL,
					"videoId":  videoId,
				}
				utils.FindBackGroundImage(res)
				return nil
			}),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
	c.JSON(http.StatusOK, gin.H{"data": scrapeResponse})
}

func storeFile(url string, id string) {
	tiktokCdnHeaders := map[string]string{
		"user-agent": "com.ss.android.ugc.trill/2613 (Linux; U; Android 10; en_US; Pixel 4; Build/QQ3A.200805.001; Cronet/58.0.2991.0)",
	}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", tiktokCdnHeaders["user-agent"])
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := client.Do(request)
	if err != nil {
		log.Panic("REQUEST ERROR", err)
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	fileName := fmt.Sprintf("tt_%s.mp4", id)
	path := constants.VideoFolder + fileName
	fileHandle, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()
	_, err = io.Copy(fileHandle, res.Body)
}
