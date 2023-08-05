package Repo

import (
	"batch/Data"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBRepository interface {
	// 全てのDBRepoに共通する処理であるDB接続を行う
	ConnectDB(isMock bool) error
	AutoMigrate() error
	// Readで使う
	SearchUserConfig(user_unique_Id string, isPreloadRelatedTables bool) (Data.UserConfig, error)
	FetchExploreCategories(country string) (resExp []Data.ExploreCategory, err error)

	// バッチ処理用
	// サイトテーブルを全件取得する
	FetchAllSites() ([]Data.WebSite, error)
	// 閲覧履歴テーブルを全件取得する 統計を取る為だからData型に変換する必要はない
	FetchAllReadHistories() ([]ReadHistory, error)
	// サイトと記事を大量に更新する
	// 記事はサイトの更新日時より新しい記事があればDBにインサートする
	UpdateSiteAndArticle(site Data.WebSite, articles []Data.Article) error
	// 時間（From・To）を指定してリードアクテビティを検索する
	SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadHistory, error)

	// テスト用
	// サイトURLをキーにサイトの最新記事を取得する
	SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error)
}

// DBRepoを実装
type DBRepoImpl struct {
}

// DB接続
func (s DBRepoImpl) ConnectDB(isMock bool) error {
	// DB接続
	if isMock {
		// もしモックモードが有効ならSqliteに接続する
		InMemoryStr := "file::memory:"
		// fileSqlStr := "test_1.db"
		DB, err := gorm.Open(sqlite.Open(InMemoryStr))
		if err != nil {
			panic("failed to connect database")
		}
		isDbConnected = true
		DBMS = DB
		return nil
	}
	if err := DataBaseConnect(); err != nil {
		return err
	}
	return nil
}

// DBアートマイグレーション
func (s DBRepoImpl) AutoMigrate() error {
	// DB接続
	if DBMS != nil {
		DBMS.AutoMigrate(&User{})
		DBMS.AutoMigrate(&FavoriteSite{})
		DBMS.AutoMigrate(&SubscriptionSite{})
		DBMS.AutoMigrate(&SearchHistory{})
		DBMS.AutoMigrate(&ReadHistory{})
		DBMS.AutoMigrate(&ApiConfig{})
		DBMS.AutoMigrate(&UiConfig{})

		DBMS.AutoMigrate(&Site{})
		DBMS.AutoMigrate(&Article{})
		DBMS.AutoMigrate(&Tag{})
		DBMS.AutoMigrate(&ExploreCategory{})
	}

	return nil
}

// Readで使う
func (s DBRepoImpl) SearchUserConfig(user_unique_Id string, isPreloadRelatedTables bool) (Data.UserConfig, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ApiConfig").Preload("UiConfig").First(&user).Error; err != nil {
		return Data.UserConfig{}, err
	}
	if isPreloadRelatedTables {
		if err := DBMS.Model(&user).Association("ReadHistories").Find(&user.ReadHistories); err != nil {
			return Data.UserConfig{}, err
		}
		if err := DBMS.Model(&user).Association("FavoriteSites").Find(&user.FavoriteSites); err != nil {
			return Data.UserConfig{}, err
		}
		if err := DBMS.Model(&user).Association("FavoriteArticles").Find(&user.FavoriteArticles); err != nil {
			return Data.UserConfig{}, err
		}
		if err := DBMS.Model(&user).Association("SubscriptionSites").Find(&user.SubscriptionSites); err != nil {
			return Data.UserConfig{}, err
		}
		if err := DBMS.Model(&user).Association("SearchHistories").Find(&user.SearchHistories); err != nil {
			return Data.UserConfig{}, err
		}
	}
	return ConvertToApiUserConfig(user), nil
}

func (s DBRepoImpl) FetchExploreCategories(country string) (res []Data.ExploreCategory, err error) {
	// ExploreCategoriesテーブルから国をキーにカテゴリーを全件取得する
	var expCats []ExploreCategory
	result := DBMS.Where(&ExploreCategory{Country: country}).Find(&expCats)
	if result.Error != nil {
		return nil, result.Error
	}
	// カテゴリーをExploreCategories型に変換する
	var categories []Data.ExploreCategory
	for _, expCat := range expCats {
		categories = append(categories, Data.ExploreCategory{
			CategoryName:        expCat.CategoryName,
			CategoryDescription: expCat.Description,
			CategoryID:          fmt.Sprint(expCat.ID),
		})
	}
	return categories, nil
}

// バッチ処理用
func (r DBRepoImpl) FetchAllSites() ([]Data.WebSite, error) {
	// サイトテーブルから全てのサイトを取得する
	var sites []Site
	result := DBMS.Find(&sites)
	if result.Error != nil {
		return nil, result.Error
	}
	// 記事を取得する必要がないので、サイトのみを返す
	var webSites []Data.WebSite
	for _, site := range sites {
		apiSite, _ := convertDbSiteToApi(site)
		webSites = append(webSites, apiSite)
	}
	return webSites, nil
}

func (r DBRepoImpl) FetchAllReadHistories() ([]ReadHistory, error) {
	// 読んだ履歴テーブルから全ての履歴を取得する
	var histories []ReadHistory
	result := DBMS.Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil
}

func (r DBRepoImpl) UpdateSiteAndArticle(site Data.WebSite, articles []Data.Article) error {
	// サイトごとのリレーションとして記事を探してそれを含めて更新
	var siteModel Site
	if err := DBMS.Where(&Site{SiteUrl: site.SiteURL}).Find(&siteModel); err.Error != nil {
		log.Println("サイトの検索に失敗しました :", err)
		return err.Error
	}
	// サイトの更新日時を更新する
	lastModified, err := time.Parse(time.RFC3339, site.LastModified)
	if err != nil {
		log.Println("更新日時の変換に失敗しました :", err.Error())
		return err
	}
	DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&siteModel).Update("LastModified", lastModified).Error; err != nil {
			log.Println("サイトの更新に失敗しました :", err.Error())
			// 失敗したらロールバック
			tx.Rollback()
			return err
		}
		return nil
	})
	// 記事の更新
	// そのサイトの記事テーブルの最新の更新日時を取得する
	var latestArticle Article
	if err := DBMS.Where(&Article{SiteID: siteModel.ID}).Order("published_at desc").Find(&latestArticle); err.Error != nil {
		return err.Error
	}
	// 記事をAPI型からDB型に変換する
	dbArticles := convertApiArticleToDb(articles)
	// 最新記事よりも新しい記事をピックアップする
	var newArticles []Article
	for _, article := range dbArticles {
		if article.PublishedAt.After(latestArticle.PublishedAt) {
			newArticles = append(newArticles, article)
		}
	}
	// トランザクションで記事をDBに保存する
	DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&siteModel).Association("SiteArticles").Append(newArticles); err != nil {
			log.Println("記事の保存に失敗しました : ", err)
			// 失敗したらロールバック
			tx.Rollback()
			return err
		}
		return nil
	})
	return nil
}

func (r DBRepoImpl) SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadHistory, error) {
	// 時間（From・To）を指定してリードアクテビティを検索する
	var histories []ReadHistory
	result := DBMS.Where("access_at BETWEEN ? AND ?", from, to).Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	var apiHistories []Data.ReadHistory
	for _, history := range histories {
		apiHistory := ConvertToApiReadHistory(history)
		apiHistories = append(apiHistories, apiHistory)
	}
	return apiHistories, nil
}

// サイトURLをキーにサイトの最新記事を取得する
func (r DBRepoImpl) SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error) {
	var site Site
	result := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return nil, result.Error
	}
	var articles []Article
	result = DBMS.Where(&Article{SiteID: site.ID}).Order("published_at desc").Limit(get_count).Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	site.SiteArticles = articles
	_, resultArticles := convertDbSiteToApi(site)
	return resultArticles, nil
}

