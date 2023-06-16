package APIFunction

import (
	"fmt"
	"heavy/Data"
	"io/ioutil"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

// サイトのRSSを取得してWebSiteを返す
func NewSite(siteUrl string) (Data.WebSite, []Data.Article, error) {
	// サイトのHTMLを取得してパースしてRSS_URLを取得する

	return Data.WebSite{}, nil, nil
}

// SiteのRSSを取得してsliceの記事を返す
func GetFeedArticle(siteUrl string) ([]Data.Article, error) {

	return nil, nil
}


func getHTML(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("HTTP status error: %d", resp.StatusCode)
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(bodyBytes), nil
}

func getRSSUrl(siteUrl string) ([]string, error) {
	// /を消す
	if siteUrl[len(siteUrl)-1] == '/' {
		siteUrl = siteUrl[:len(siteUrl)-1]
	}
	resp, err := http.Get(siteUrl)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
	// RSS_URLが複数ある場合は別のサイトとして扱う
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP status error: %d", resp.StatusCode)
    }
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil,  err
    }
    rssUrl :=[]string{}
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
	return rssUrl, nil
}

