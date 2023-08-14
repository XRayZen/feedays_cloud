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
	SearchReadActivityByTime(user_unique_id string,from time.Time, to time.Time) ([]Data.ReadHistory, error)

	// テスト用
	// サイトURLをキーにサイトの最新記事を取得する
	SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error)
}

// DBRepoを実装
type DBRepoImpl struct {
}

// DB接続
func (s DBRepoImpl) ConnectDB(is_mock_mode bool) error {
	// DB接続
	if is_mock_mode {
		// もしモックモードが有効ならSqliteに接続する
		in_memory_str := "file::memory:"
		// fileSqlStr := "test_1.db"
		db, err := gorm.Open(sqlite.Open(in_memory_str))
		if err != nil {
			panic("failed to connect database")
		}
		isDbConnected = true
		DBMS = db
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
		DBMS.AutoMigrate(&ApiLimitConfig{})
		DBMS.AutoMigrate(&UiConfig{})

		DBMS.AutoMigrate(&Site{})
		DBMS.AutoMigrate(&Article{})
		DBMS.AutoMigrate(&Tag{})
		DBMS.AutoMigrate(&ExploreCategory{})
	}

	return nil
}

// Readで使う
func (s DBRepoImpl) SearchUserConfig(user_unique_Id string, is_preload_related_tables bool) (Data.UserConfig, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("UiConfig").First(&user).Error; err != nil {
		return Data.UserConfig{}, err
	}
	if is_preload_related_tables {
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
	var exp_catagories []ExploreCategory
	result := DBMS.Where(&ExploreCategory{Country: country}).Find(&exp_catagories)
	if result.Error != nil {
		return nil, result.Error
	}
	// カテゴリーをExploreCategories型に変換する
	var categories []Data.ExploreCategory
	for _, exp_category := range exp_catagories {
		categories = append(categories, Data.ExploreCategory{
			CategoryName:        exp_category.CategoryName,
			CategoryDescription: exp_category.Description,
			CategoryID:          fmt.Sprint(exp_category.ID),
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
	var web_sites []Data.WebSite
	for _, site := range sites {
		api_site, _ := convertDbSiteToApi(site)
		web_sites = append(web_sites, api_site)
	}
	return web_sites, nil
}

// 読んだ履歴テーブルから全ての履歴を取得する
func (r DBRepoImpl) FetchAllReadHistories() ([]ReadHistory, error) {
	var histories []ReadHistory
	result := DBMS.Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil
}

// サイトごとのリレーションとして記事を探してそれを含めて更新
func (r DBRepoImpl) UpdateSiteAndArticle(site Data.WebSite, articles []Data.Article) error {
	var site_model Site
	if err := DBMS.Where(&Site{SiteUrl: site.SiteURL}).Find(&site_model); err.Error != nil {
		log.Println("サイトの検索に失敗しました :", err)
		return err.Error
	}
	// サイトの更新日時を更新する
	last_modified, err := time.Parse(time.RFC3339, site.LastModified)
	if err != nil {
		log.Println("更新日時の変換に失敗しました :", err.Error())
		return err
	}
	DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&site_model).Update("LastModified", last_modified).Error; err != nil {
			log.Println("サイトの更新に失敗しました :", err.Error())
			// 失敗したらロールバック
			tx.Rollback()
			return err
		}
		return nil
	})
	// 記事の更新
	// そのサイトの記事テーブルの最新の更新日時を取得する
	var latest_article Article
	if err := DBMS.Where(&Article{SiteID: site_model.ID}).Order("published_at desc").Find(&latest_article); err.Error != nil {
		return err.Error
	}
	// 記事をAPI型からDB型に変換する
	db_articles := convertApiArticleToDb(articles)
	// 最新記事よりも新しい記事をピックアップする
	var new_articles []Article
	for _, article := range db_articles {
		if article.PublishedAt.After(latest_article.PublishedAt) {
			new_articles = append(new_articles, article)
		}
	}
	// トランザクションで記事をDBに保存する
	DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&site_model).Association("SiteArticles").Append(new_articles); err != nil {
			log.Println("記事の保存に失敗しました : ", err)
			// 失敗したらロールバック
			tx.Rollback()
			return err
		}
		return nil
	})
	return nil
}

// 時間（From・To）を指定してリードアクテビティを検索する
func (r DBRepoImpl) SearchReadActivityByTime(user_unique_id string,from time.Time, to time.Time) ([]Data.ReadHistory, error) {
	var histories []ReadHistory
	result := DBMS.Where("access_at BETWEEN ? AND ?", from, to).Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	var api_histories []Data.ReadHistory
	for _, history := range histories {
		api_history := ConvertToApiReadHistory(history,user_unique_id)
		api_histories = append(api_histories, api_history)
	}
	return api_histories, nil
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
	_, result_articles := convertDbSiteToApi(site)
	return result_articles, nil
}
