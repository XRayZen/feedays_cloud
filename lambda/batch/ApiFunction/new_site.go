package ApiFunction

import (
	"fmt"
	"net/http"
	"batch/Data"

	"github.com/PuerkitoBio/goquery"
)

// 新規サイトを調べて、サイト情報と記事情報を返す
func newSite(siteUrl string) (Data.WebSite, []Data.Article, error) {
	doc, err := getHtmlGoQueryDoc(siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getHtmlGoQueryDoc error: %v", err)
	}
	// RSSのURLを取得する
	rss_urls, err := getRSSUrls(doc, siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getRSSUrl error: %v", err)
	}
	// サイトメタデータを取得する
	site_meta, err := getSiteMeta(doc, siteUrl)
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("getSiteMeta error: %v", err)
	}
	// RSSをパースする
	articles, err := fetchRSSArticles(rss_urls[0])
	if err != nil {
		return Data.WebSite{}, nil, fmt.Errorf("parseRssFeed error: %v", err)
	}
	site_meta.SiteRssURL = rss_urls[0]
	return site_meta, articles, nil
}

func getHtmlGoQueryDoc(url string) (*goquery.Document, error) {
	// /を消す
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status error: %d", response.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery error: %v", err)
	}
	return doc, nil
}

func getRSSUrls(doc *goquery.Document, siteUrl string) ([]string, error) {
	rss_url := []string{}
	doc.Find("link[type='application/rss+xml']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			rss_url = append(rss_url, href)
			return
		}
	})
	// atomの場合
	doc.Find("link[type='application/atom+xml']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			rss_url = append(rss_url, href)
			return
		}
	})
	if rss_url == nil {
		return nil, fmt.Errorf("RSS URL not found")
	}
	//RSSのURLが相対パスの場合の処理
	for i, v := range rss_url {
		if v[0] == '/' {
			rss_url[i] = siteUrl + v
		}
	}
	for i, v := range rss_url {
		// 最後のURLに/がある場合は消す
		if v[len(v)-1] == '/' {
			rss_url[i] = v[:len(v)-1]
		}
	}
	return rss_url, nil
}

// サイトのHTMLを取得してパースしてメタ情報を取得する
func getSiteMeta(doc *goquery.Document, siteUrl string) (Data.WebSite, error) {
	// メタ情報を取得する
	site_name := ""
	site_image := ""
	site_description := ""
	site_tags := []string{}
	doc.Find("head meta").Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if exists {
			switch name {
			case "og:title":
				site_name = s.AttrOr("content", "")
			case "og:image":
				site_image = s.AttrOr("content", "")
			case "og:description":
				site_description = s.AttrOr("content", "")
			case "keywords":
				site_tags = append(site_tags, s.AttrOr("content", ""))
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
				site_name = s.AttrOr("content", "")
			case "og:image":
				site_image = s.AttrOr("content", "")
			case "og:description":
				site_description = s.AttrOr("content", "")
			}
			return
		}
	})
	return Data.WebSite{
		SiteName:        site_name,
		SiteImage:       site_image,
		SiteDescription: site_description,
		SiteURL:         siteUrl,
		SiteTags:        site_tags,
	}, nil
}
