package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ApiData struct {
	Has_more int    `json:"has_more"`
	Data     []Data `json:"data"`
}

type Data struct {
	Title       string `json:"title"`
	Article_url string `json:"article_url"`
}

type Img struct {
	Src string `json:"src"`
}

var (
	host    string = "http://www.toutiao.com/search_content/?format=json&keyword=%s&count=30&offset=%d"
	hasmore bool   = true
)

func main() {
	for _, tag := range os.Args[1:] {
		hasmore = true
		getByTag(tag)
	}
	log.Println("全部抓取完毕")
}

func getByTag(tag string) {
	i, offset := 1, 140
	for {
		if hasmore {
			log.Printf("标签: '%s'，第 '%d' 页, OFFSET: '%d' \n", tag, i, offset)
			tmpUrl := fmt.Sprintf(host, tag, offset)
			getResFromApi(tmpUrl)
			offset += 20
			i++

			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}
	log.Printf("标签: '%s', 共 %v 页，爬取完毕\n", tag, i)
}

func getResFromApi(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var res ApiData
	json.Unmarshal([]byte(string(body)), &res)

	for _, item := range res.Data {
		getImgByPage(item.Article_url)
	}

	if res.Has_more == 0 {
		hasmore = false
	}
}

func getImgByPage(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	title := doc.Find("#article-main .article-title").Text()
	title = strings.Replace(title, "/", "", -1)
	os.Mkdir(title, 0777)

	doc.Find("#J_content .article-content img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		log.Println(title, src)
		getImgAndSave(src+".jpg", title)
	})
}

func getImgAndSave(url string, dirname string) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}

	resp, err := http.Get(url)
	defer func() {
		if x := recover(); x != nil {
			return
		}
	}()
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile("./"+dirname+"/"+name, contents, 0644)
	if err != nil {
		log.Fatal("写入文件失败", err)
	}
}
