package ApiFunction

import (
	"batch/Data"
	"batch/Repo"
	"testing"
	"time"
)

func TestBatch(t *testing.T) {
	// ギガジンのRSSフィードを取得する
	// 記事最新から10件までをリストに入れる
	// 正解の記事リストは最新から10件まで
	// それをDB型に変換してからDBに登録する
	db_repo := Repo.DBRepoImpl{}
	answer_articles := setup(t, db_repo)
	// その後、DBから記事を取得してテストする
	// 最新から10件までが正しくDBに登録してたらテスト成功
	t.Run("Batch", func(t *testing.T) {
		// Batch処理を実行
		result, err := Batch(db_repo, 15)
		if err != nil {
			t.Fatal(err)
		}
		// テスト結果判定
		// DBから記事を取得して正解の記事リストと比較する
		if result != true {
			t.Fatal("Batch処理が失敗しています")
		}
		// ちゃんとDBにサイトが登録されているか確認
		var sites []Repo.Site
		err = Repo.DBMS.Find(&sites).Error
		if err != nil {
			t.Fatal(err)
		}
		result_articles := []Repo.Article{}
		if err = Repo.DBMS.Model(&sites[0]).Association("SiteArticles").Find(&result_articles); err != nil {
			t.Fatal(err)
		}
		// 結果の記事リストをPublishedAtでソートする
		for i := 0; i < len(result_articles); i++ {
			for j := i + 1; j < len(result_articles); j++ {
				if result_articles[i].PublishedAt.Before(result_articles[j].PublishedAt) {
					tmp := result_articles[i]
					result_articles[i] = result_articles[j]
					result_articles[j] = tmp
				}
			}
		}
		// 結果の記事リストと正解の記事リストを比較してちゃんと最新記事が登録されているか確認する
		if result_articles[0].Title != answer_articles[0].Title {
			t.Fatal("記事が10件登録されていません")
		}
	})
}

func setup(t *testing.T, db_repo Repo.DBRepository) []Data.Article {
	db_repo.ConnectDB(true)
	db_repo.AutoMigrate()

	articles, err := fetchRSSArticles("https://gigazine.net/news/rss_2.0/")
	if err != nil {
		t.Fatal(err)
	}
	insert_articles := []Data.Article{}
	// 最後から10件までを取得
	for i := len(articles) - 1; i > len(articles)-11; i-- {
		insert_articles = append(insert_articles, articles[i])
	}

	answer_articles := articles[0:10]
	db_insert_articles := []Repo.Article{}
	for _, article := range insert_articles {
		publication_date, _ := time.Parse(time.RFC3339, article.PublishedAt)
		db_article := Repo.Article{
			Title:       article.Title,
			Url:         article.Link,
			Description: article.Description,
			PublishedAt: publication_date,
		}
		db_insert_articles = append(db_insert_articles, db_article)
	}

	dbSite := Repo.Site{
		SiteName:     "GIGAZINE",
		RssUrl:       "https://gigazine.net/news/rss_2.0/",
		SiteUrl:      "https://gigazine.net/",
		Description:  "ギガジンのRSSフィード",
		SiteArticles: db_insert_articles,
	}
	if err := Repo.DBMS.Create(&dbSite).Error; err != nil {
		t.Fatal(err)
	}
	return answer_articles
}
