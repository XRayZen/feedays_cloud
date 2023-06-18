package ApiFunction

import (
	"log"
	"read/Data"
	"read/Repo"
	"time"
)

func RefreshArticles(repo Repo.DBRepository) (bool, error) {
	// バッチ処理を実行する
	// サイトテーブルから読み込んで記事を更新する
	// それにより購読サイトのフィードの鮮度を維持する
	sites, err := repo.FetchAllSites()
	if err != nil {
		return false, err
	}
	type temp struct {
		site     Data.WebSite
		articles []Data.Article
		err      error
	}
	// 並列処理でサイトの記事を更新する
	ch := make(chan temp, len(sites))
	for _, site := range sites {
		go func(site Data.WebSite) {
			// RSSフィードを取得する
			articles, err := fetchRSSArticles(site.SiteRssURL)
			if err != nil {
				ch <- temp{site: site, articles: nil, err: err}
			}
			ch <- temp{site: site, articles: articles, err: nil}
		}(site)
	}
	// 並列処理の結果を受け取る
	fetch_articles_result := []temp{}
	for i := 0; i < len(sites); i++ {
		t := <-ch
		if t.err != nil {
			log.Println("BATCH RefreshArticles ERROR! :", t.err)
			log.Printf("BATCH RefreshArticles ERROR! name:%s Site:%s ", t.site.SiteName, t.site.SiteURL)
		}
		fetch_articles_result = append(fetch_articles_result, t)
	}
	// 同期でループしてイメージURLは並列で取得する
	result_articles := []Data.Article{}
	result_sites := []Data.WebSite{}
	for _, t := range fetch_articles_result {
		articles, err := getArticleImageURLs(t.articles)
		if err != nil {
			log.Println("BATCH RefreshArticles get_Images ERROR! :", err)
			log.Printf("BATCH RefreshArticles get_Images ERROR! name:%s Site:%s ", t.site.SiteName, t.site.SiteURL)
		}
		// サイトの更新日時を更新する
		t.site.LastModified = time.Now().Format(time.RFC3339)
		result_articles = append(result_articles, articles...)
		result_sites = append(result_sites, t.site)
	}
	// DBに記事を保存する
	if err := repo.UpdateSitesAndArticles(result_sites,result_articles); err != nil {
		log.Println("BATCH RefreshArticles DB ERROR! :", err)
		return false, err
	}
	return true, nil
}
