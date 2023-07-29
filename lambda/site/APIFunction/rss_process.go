package APIFunction

import (
	"site/Data"
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

// 指定されたサイトのRSS_URLからRSSフィードを取得して記事リストとして返す
func fetchRSSArticles(rssUrl string) ([]Data.Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssUrl)
	if err != nil {
		return nil, err
	}
	articles := []Data.Article{}
	for _, v := range feed.Items {
		// Feedのカテゴリはタグにしておく
		category := ""
		if len(v.Categories) > 0 {
			category = v.Categories[0]
		}
		article := Data.Article{
			Title:       v.Title,
			Link:        v.Link,
			Description: v.Description,
			Category:    category,
			Site:        feed.Title,
			PublishedAt: v.PublishedParsed.Format(time.RFC3339),
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// 並列処理で記事のイメージURLを取得する
func getArticleImageURLs(articles []Data.Article) ([]Data.Article, error) {
	// 並列処理で記事のイメージURLを取得する
	// 1. og:imageを取得する
	// 3. それでもなければfavicon.icoを取得する
	ch := make(chan Data.Article)
	for _, article := range articles {
		go func(article Data.Article) {
			// 1. og:imageを取得する
			doc, err := getHtmlGoQueryDoc(article.Link)
			if err != nil {
				ch <- article
				return
			}
			imageUrl, err := getArticleImageURL(doc, article.Link)
			if err != nil {
				ch <- article
				return
			}
			article.Image = Data.RssFeedImage{
				Link: imageUrl,
			}
			ch <- article
		}(article)
	}
	for i := 0; i < len(articles); i++ {
		articles[i] = <-ch
	}
	// articleを日時でソートする
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt > articles[j].PublishedAt
	})
	return articles, nil
}

// 記事のイメージURLを取得する
func getArticleImageURL(doc *goquery.Document, articleUrl string) (string, error) {
	// 記事のイメージURLを取得する
	// 1. og:imageを取得する
	// 3. それでもなければfavicon.icoを取得する
	imageUrl := ""
	// 1. og:imageを取得する
	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		property, exists := s.Attr("property")
		if exists {
			if property == "og:image" {
				imageUrl = s.AttrOr("content", "")
				return
			}
		}
	})
	// 2. それでもなければfavicon.icoを取得する
	if imageUrl == "" {
		imageUrl = articleUrl + "/favicon.ico"
	}
	return imageUrl, nil
}
