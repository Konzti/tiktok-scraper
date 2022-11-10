package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	mathRand "math/rand"
	"net/http"
	"os"
	"regexp"
	"tiktok-scraper/constants"
	"time"
)

//headers := map[string]string{
//	"user-agent": "Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36 Edg/87.0.664.66",
//}
//tiktok_headers := map[string]string{
//	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
//	"authority":       "www.tiktok.com",
//	"Accept-Encoding": "gzip, deflate",
//	"Connection":      "keep-alive",
//	"Host":            "www.tiktok.com",
//	"User-Agent":      "Mozilla/5.0  (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) coc_coc_browser/86.0.170 Chrome/80.0.3987.170 Safari/537.36",
//}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
func ReturnError(err error) error {
	fmt.Println(err)
	return err
}

func createFolderIfNotExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
func CreateVolumeFolders() {
	for _, folder := range constants.Folders {
		createFolderIfNotExist(folder)
	}
}

func LoadEnv() {
	err := godotenv.Load(constants.EnvFileLocation)
	if err != nil {
		fmt.Println("Couldn't find .env file")
	}
}

func RandomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(bytes)
}

func RandomId() string {
	chars := "01234567890123456"
	var seededRand *mathRand.Rand = mathRand.New(
		mathRand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 16)
	for i := range b {
		b[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(b)
}

func SearchForCanonical(url string) (foundUrl string, err error) {
	c := http.Client{Timeout: time.Duration(constants.TimeOutInSeconds) * time.Second}
	// Request the HTML page.
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		err = ReturnError(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36")
	resp, err := c.Do(req)

	if err != nil {
		err = ReturnError(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(err)
		err = errors.New("no 200 response")
	}
	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		err = errors.New("error reading body")
	}

	//imageRegExp := regexp.MustCompile(`<meta[^>]+\bcontent=["']([^"']+)["']`)
	imageRegExp := regexp.MustCompile(`https://www.tiktok.com/@[^/]+/video/[^"/]+`)

	subMatchSlice := imageRegExp.FindAllStringSubmatch(string(htmlData), -1)
	foundUrl = subMatchSlice[len(subMatchSlice)-1][0]
	fmt.Println(foundUrl)
	return foundUrl, nil
}

// Not needed bc of chromedp
func SearchForImageLinks(url string) {

	log.Println("Parsing : ", url)
	c := http.Client{Timeout: time.Duration(3) * time.Second}
	// Request the HTML page.
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.112 Safari/537.36")
	resp, err := c.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("Unable to get URL with status code error: %d %s", resp.StatusCode, resp.Status)
	}
	//fmt.Println(resp.Cookies())

	htmlData, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(htmlData))

	if err != nil {
		log.Fatal(err)
	}

	imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)

	subMatchSlice := imageRegExp.FindAllStringSubmatch(string(htmlData), -1)
	imageUrl := subMatchSlice[0][1]
	userAvatar := subMatchSlice[1][1]
	fmt.Println(imageUrl)
	fmt.Println("Avatar: ", userAvatar)
	//for i, item := range subMatchSlice {
	//	log.Printf("Item found at index %d : %s ", i, item[1])
	//}

}

func FindVideoStrings(body string) (videoStrings []string) {
	imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	subMatchSlice := imageRegExp.FindAllStringSubmatch(body, -1)
	imageUrl := subMatchSlice[0][1]
	userAvatar := subMatchSlice[1][1]
	fmt.Println(imageUrl)
	fmt.Println("Avatar: ", userAvatar)
	//for i, item := range subMatchSlice {
	//	log.Printf("Item found at index %d : %s ", i, item[1])
	//}
	videoStrings = append(videoStrings, imageUrl, userAvatar)
	return videoStrings
}
func FindBackGroundImage(body string) {
	//imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imageRegExp := regexp.MustCompile(`<img mode="4"[^>]+\bsrc=["']([^"']+)["']`)
	subMatchSlice := imageRegExp.FindAllStringSubmatch(body, -1)
	//imageUrl := subMatchSlice[0][1]
	//userAvatar := subMatchSlice[1][1]
	//fmt.Println(imageUrl)
	//fmt.Println("Avatar: ", userAvatar)
	for i, item := range subMatchSlice {
		log.Printf("Item found at index %d : %s ", i, item[1])
	}
	//videoStrings = append(videoStrings, imageUrl, userAvatar)
	//return videoStrings
}

func FindVideo(body string) string {
	imageRegExp := regexp.MustCompile(`<video[^>]+\bsrc=["']([^"']+)["']`)
	subMatchSlice := imageRegExp.FindAllStringSubmatch(body, -1)
	videoUrl := subMatchSlice[0][1]
	return videoUrl
}
