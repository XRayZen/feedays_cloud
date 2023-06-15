package APIFunction

import "heavy/Data"

// サイトのRSSを取得してWebSiteを返す
func NewSite(siteUrl string) (Data.WebSite, []Data.Article, error) {
	// サイトのHTMLを取得してパースしてRSS_URLを取得する

	return Data.WebSite{}, nil, nil
}

// SiteのRSSを取得してsliceの記事を返す
func GetFeedArticle(siteUrl string) ([]Data.Article, error) {

	return nil, nil
}
