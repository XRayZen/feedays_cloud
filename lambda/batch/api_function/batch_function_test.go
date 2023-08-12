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
	dbRepo := Repo.DBRepoImpl{}
	answerArticles := setup(t, dbRepo)
	// その後、DBから記事を取得してテストする
	// 最新から10件までが正しくDBに登録してたらテスト成功
	t.Run("Batch", func(t *testing.T) {
		// Batch処理を実行
		result, err := Batch(dbRepo, 15)
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
		resultArticles := []Repo.Article{}
		if err = Repo.DBMS.Model(&sites[0]).Association("SiteArticles").Find(&resultArticles); err != nil {
			t.Fatal(err)
		}
		// 結果の記事リストをPublishedAtでソートする
		for i := 0; i < len(resultArticles); i++ {
			for j := i + 1; j < len(resultArticles); j++ {
				if resultArticles[i].PublishedAt.Before(resultArticles[j].PublishedAt) {
					tmp := resultArticles[i]
					resultArticles[i] = resultArticles[j]
					resultArticles[j] = tmp
				}
			}
		}
		// 結果の記事リストと正解の記事リストを比較してちゃんと最新記事が登録されているか確認する
		if resultArticles[0].Title != answerArticles[0].Title {
			t.Fatal("記事が10件登録されていません")
		}
	})
}

func setup(t *testing.T, dbRepo Repo.DBRepository) []Data.Article {
	dbRepo.ConnectDB(true)
	dbRepo.AutoMigrate()

	articles, err := fetchRSSArticles("https://gigazine.net/news/rss_2.0/")
	if err != nil {
		t.Fatal(err)
	}
	insertArticles := []Data.Article{}
	// 最後から10件までを取得
	for i := len(articles) - 1; i > len(articles)-11; i-- {
		insertArticles = append(insertArticles, articles[i])
	}

	answerArticles := articles[0:10]
	dbInsertArticles := []Repo.Article{}
	for _, article := range insertArticles {
		publicationDate, _ := time.Parse(time.RFC3339, article.PublishedAt)
		dbArticle := Repo.Article{
			Title:       article.Title,
			Url:         article.Link,
			Description: article.Description,
			PublishedAt: publicationDate,
		}
		dbInsertArticles = append(dbInsertArticles, dbArticle)
	}

	dbSite := Repo.Site{
		SiteName:     "GIGAZINE",
		RssUrl:       "https://gigazine.net/news/rss_2.0/",
		SiteUrl:      "https://gigazine.net/",
		Description:  "ギガジンのRSSフィード",
		SiteArticles: dbInsertArticles,
	}
	if err := Repo.DBMS.Create(&dbSite).Error; err != nil {
		t.Fatal(err)
	}
	return answerArticles
}
