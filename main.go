package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"tiktok-scraper/constants"
	"tiktok-scraper/structs"
	"tiktok-scraper/utils"
	"time"
)

func main() {
	utils.LoadEnv()
	utils.CreateVolumeFolders()

	if mode := os.Getenv("GIN_MODE"); mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	fmt.Println("Mode is in: ", gin.Mode())

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type", "Content-Disposition", "Access-Control-Allow-Origin"}
	r.Use(cors.New(config))

	// PROD: SEND VITE APP TO CLIENT
	r.Use(static.Serve("/", static.LocalFile("./dist", true)))

	r.GET("/video/:id", VideoHandler)
	r.POST("/url", ApiHandler)

	utils.HandleError(r.Run())
}

func ApiHandler(c *gin.Context) {
	var requestBody structs.URL
	// Checking for valid JSON
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Checking for valid URL ...
	parsedUrl, err := url.ParseRequestURI(requestBody.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if parsedUrl.Host != "www.tiktok.com" {
		if parsedUrl.Host != "vm.tiktok.com" && parsedUrl.Host != "vt.tiktok.com" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		} else {
			secondPart := strings.Split(parsedUrl.Path, "/")[1]
			if len(secondPart) < 4 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
				return
			}
			//utils.SearchForImageLinks(requestBody.URL)

			foundUrl, err := utils.SearchForCanonical(requestBody.URL)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
				return
			}
			newParsedUrl, err := url.ParseRequestURI(foundUrl)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// SCRAPE VIDEO!!!
			ScrapeForVideo(newParsedUrl, c)
			return
		}
	}
	ScrapeForVideo(parsedUrl, c)

}

func ScrapeForVideo(parsedUrl *url.URL, c *gin.Context) {
	if !strings.Contains(parsedUrl.Path, "/video/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid video path"})
		return
	}
	videoId := strings.Split(parsedUrl.Path, "/")[3]
	fmt.Println("Video ID: ", videoId)
	//utils.SearchForImageLinks(requestBody.URL)

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
				go StoreFile(videoURL, videoId)
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
	//for _, n := range nodes {
	//	fmt.Println(n)
	//}
	//fmt.Printf("h1 contains: '%s'\n", res)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
	c.JSON(http.StatusOK, gin.H{"data": scrapeResponse})
}

func VideoHandler(c *gin.Context) {
	id := c.Param("id")
	fileName := fmt.Sprintf("tt_%s.mp4", id)
	path := constants.VideoFolder + fileName
	_, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist anymore"})
	}

	cd := mime.FormatMediaType("attachment", map[string]string{"filename": fileName})
	c.Header("Content-Disposition", cd)
	c.Header("Content-Type", "application/octet-stream")
	fmt.Println("Sending file.....", fileName)
	c.FileAttachment(path, fileName)
}

func StoreFile(url string, id string) {
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
