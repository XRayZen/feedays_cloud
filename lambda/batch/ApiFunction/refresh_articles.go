package ApiFunction

import (
	"batch/Repo"
	"batch/Data"
	"errors"
	"log"
	"time"
)

func RefreshArticles(repo Repo.DBRepository,RefreshInterval int) (string, error) {
	// バッチ処理を実行する
	// サイトテーブルから読み込んで記事を更新する
	// それによりサイトのフィード鮮度を維持する
	sites, err := repo.FetchAllSites()
	if err != nil {
		log.Println("BATCH RefreshArticles DB Read ERROR! :", err)
		return "DB Error", err
	}
	type temp struct {
		site     Data.WebSite
		articles []Data.Article
		err      error
		isNotExp bool
	}
	// 並列処理でサイトの記事を更新する
	ch := make(chan temp, len(sites))
	for _, site := range sites {
		go func(site Data.WebSite) {
			time, err := time.Parse(time.RFC3339, site.LastModified)
			if err != nil {
				log.Println("BATCH RefreshArticles ERROR! :", err)
			}
			if isUpdateExpired(time, RefreshInterval) {
				articles, err := fetchRSSArticles(site.SiteRssURL)
				if err != nil {
					ch <- temp{site: site, articles: nil, err: err}
				}
				ch <- temp{site: site, articles: articles, err: nil}
			}
			ch <- temp{site: site, articles: nil, err: errors.New("not Expired"), isNotExp: true}
		}(site)
	}
	// 並列処理の結果を受け取る
	fetch_articles_result := []temp{}
	for i := 0; i < len(sites); i++ {
		t := <-ch
		if t.isNotExp {
			continue
		}
		if t.err != nil {
			log.Println("BATCH RefreshArticles ERROR! :", t.err)
			log.Printf("BATCH RefreshArticles ERROR! name:%s Site:%s ", t.site.SiteName, t.site.SiteURL)
		}
		fetch_articles_result = append(fetch_articles_result, t)
	}
	// 更新する記事がない場合は終了する
	if len(fetch_articles_result) == 0 {
		return "No Update", nil
	}
	// 同期でループして記事イメージURLは並列で取得する
	for _, t := range fetch_articles_result {
		articles, err := getArticleImageURLs(t.articles)
		if err != nil {
			log.Println("BATCH RefreshArticles get_Images ERROR! :", err)
			log.Printf("BATCH RefreshArticles get_Images ERROR! name:%s Site:%s ", t.site.SiteName, t.site.SiteURL)
		}
		// サイトの更新日時を更新する
		t.site.LastModified = time.Now().Format(time.RFC3339)
		// ここでDBに保存するサイトごと記事を更新する
		if err := repo.UpdateSiteAndArticle(t.site, articles); err != nil {
			log.Println("BATCH RefreshArticles DB Update ERROR! :", err)
			log.Printf("BATCH RefreshArticles DB Update ERROR! name:%s Site:%s ", t.site.SiteName, t.site.SiteURL)
		}
	}
	result_msg := "BATCH RefreshArticles SUCCESS!"
	return result_msg, nil
}

// 記事更新日時にIntervalMinutesを足した更新期限日時を現時間が過ぎていたらtrueを返す
func isUpdateExpired(lastModified time.Time, intervalMinutes int) bool {
	// 現時間を取得する
	now_time := time.Now()
	// 記事更新日時にIntervalMinutesを足した更新期限日時を取得する
	update_expired_time := lastModified.Add(time.Minute * time.Duration(intervalMinutes))
	// 更新期限日時が現時間より過ぎていたらtrueを返す
	return update_expired_time.Before(now_time)
}
