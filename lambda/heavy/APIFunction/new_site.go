package APIFunction

import (
	"fmt"
	"read/Data"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// 新規サイトを調べて、サイト情報と記事情報を返す
func newSite(siteUrl string) (Data.WebSite, []Data.Article, error) {
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
	siteMeta, err := getSiteMeta(doc, siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getSiteMeta error: %v", err)
	}
	// RSSをパースする
	articles, err := parseRssFeed(rssUrls[0])
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("parseRssFeed error: %v", err)
	}
	siteMeta.SiteRssURL = rssUrls[0]
	return siteMeta, articles, nil
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
	for i, v := range rssUrl {
		// 最後のURLに/がある場合は消す
		if v[len(v)-1] == '/' {
			rssUrl[i] = v[:len(v)-1]
		}
	}
	return rssUrl, nil
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
			case "keywords":
				siteTags = append(siteTags, s.AttrOr("content", ""))
			}
			// // カテゴリーも取得する
			// if name == "keywords" {
			// 	siteTags = append(siteTags, s.AttrOr("content", ""))
			// }
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
