package APIFunction

import (
	"fmt"
	"read/Data"
	// "heavy/Data"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

// サイトのRSSを取得してWebSiteを返す
func NewSite(siteUrl string) (Data.WebSite, []Data.Article, error) {
	doc, err := getHtmlGoQueryDoc(siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getHtmlGoQueryDoc error: %v", err)
	}
	// RSSのURLを取得する
	rssUrls, err := getRSSUrls(doc, siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getRSSUrl error: %v", err)
	}
	// サイトメタデータを取得する
	siteMeta, err := getSiteMeta(doc,siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getSiteMeta error: %v", err)
	}
	// RSSをパースする
	articles := []Data.Article{}
	

	return Data.WebSite{}, nil, nil
}

// SiteのRSSを取得してsliceの記事を返す
func GetFeedArticle(siteUrl string) ([]Data.Article, error) {

	return nil, nil
}

func getHtmlGoQueryDoc(url string) (*goquery.Document, error) {
	// /を消す
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery error: %v", err)
	}
	return doc, nil
}

func getRSSUrls(doc *goquery.Document, siteUrl string) ([]string, error) {
	rssUrl := []string{}
	doc.Find("link[type='application/rss+xml']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			rssUrl = append(rssUrl, href)
			return
		}
	})
	// atomの場合
	doc.Find("link[type='application/atom+xml']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			rssUrl = append(rssUrl, href)
			return
		}
	})
	if rssUrl == nil {
		return nil, fmt.Errorf("RSS URL not found")
	}
	//RSSのURLが相対パスの場合の処理
	for i, v := range rssUrl {
		if v[0] == '/' {
			rssUrl[i] = siteUrl + v
		}
	}
	// 最後のURLに/がある場合は消す
	for i, v := range rssUrl {
		if v[len(v)-1] == '/' {
			rssUrl[i] = v[:len(v)-1]
		}
	}
	return rssUrl, nil
}

// RssFeedをパースする
func parseRssFeed(rssUrl string) (Data.WebSite, []Data.Article, error) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(rssUrl)
	articles := []Data.Article{}
	for _, v := range feed.Items {
		article := Data.Article{
			Title:        v.Title,
			Link:         v.Link,
			Description:  v.Description,
			Category:     v.Categories,
			Site:         feed.Title,
			LastModified: v.PublishedParsed.Format(time.RFC3339),
		}
		articles = append(articles, article)
	}
	cate := ""
	if len(feed.Categories) != 0 {
		cate = feed.Categories[0]
	}
	imageUrl := ""
	if feed.Image != nil {
		imageUrl = feed.Image.URL
	}
	return Data.WebSite{
		SiteName:        feed.Title,
		SiteImage:       imageUrl,
		SiteDescription: feed.Description,
		SiteURL:         feed.Link,
		SiteRssURL:      rssUrl,
		SiteCategory:    cate,
	}, articles, nil
}

// サイトのHTMLを取得してパースしてメタ情報を取得する
func getSiteMeta(doc *goquery.Document, siteUrl string) (Data.WebSite, error) {
	// メタ情報を取得する
	siteName := ""
	siteImage := ""
	siteDescription := ""
	siteTags := []string{}
	doc.Find("head meta").Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if exists {
			switch name {
			case "og:title":
				siteName = s.AttrOr("content", "")
			case "og:image":
				siteImage = s.AttrOr("content", "")
			case "og:description":
				siteDescription = s.AttrOr("content", "")
			}
			// カテゴリーも取得する
			if name == "keywords" {
				siteTags = append(siteTags, s.AttrOr("content", ""))
			}
			return
		}

		property, exists := s.Attr("property")
		if exists {
			switch property {
			case "og:title":
				siteName = s.AttrOr("content", "")
			case "og:image":
				siteImage = s.AttrOr("content", "")
			case "og:description":
				siteDescription = s.AttrOr("content", "")
			}
			return
		}
	})
	return Data.WebSite{
		SiteName:        siteName,
		SiteImage:       siteImage,
		SiteDescription: siteDescription,
		SiteURL:         siteUrl,
		SiteTags:        siteTags,
	}, nil
}
