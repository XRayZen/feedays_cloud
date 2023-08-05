package Repo

import (
	"batch/Data"
	"errors"
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
	SearchUserConfig(user_unique_Id string,isPreloadRelatedTables bool) (Data.UserConfig, error)
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
func (s DBRepoImpl) SearchUserConfig(user_unique_Id string,isPreloadRelatedTables bool) (Data.UserConfig, error) {
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

// heavyで使う

// サイト系処理
// 登録・存在確認・購読・購読確認・URL検索
func (r DBRepoImpl) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	// サイトを登録する
	// もし、サイトが存在していたら登録しない
	if !r.IsExistSite(site.SiteURL) {
		dbSite := convertApiSiteToDb(site, articles)
		DBMS.Transaction(func(tx *gorm.DB) error {
			result := tx.Create(&dbSite)
			if result.Error != nil {
				log.Println("サイト登録失敗 : " + result.Error.Error())
				tx.Rollback()
				return result.Error
			}
			return nil
		})
	} else {
		return errors.New("サイトが既に存在しています")
	}
	return nil
}

func (r DBRepoImpl) SearchSiteByUrl(site_url string) (Data.WebSite, error) {
	var site Site
	result := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return Data.WebSite{}, result.Error
	}
	resultSite, _ := convertDbSiteToApi(site)
	return resultSite, nil
}

func (r DBRepoImpl) SearchSiteByName(siteName string) ([]Data.WebSite, error) {
	var sites []Site
	result := DBMS.Where(&Site{SiteName: siteName}).Find(&sites)
	if result.Error != nil {
		return []Data.WebSite{}, result.Error
	}
	var resultSites []Data.WebSite
	for _, site := range sites {
		resultSite, _ := convertDbSiteToApi(site)
		resultSites = append(resultSites, resultSite)
	}
	return resultSites, nil
}

func (r DBRepoImpl) IsExistSite(site_url string) bool {
	var count int64
	result := DBMS.Model(&Site{}).Where(&Site{SiteUrl: site_url}).Count(&count)
	if result.Error != nil {
		return false
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

func (r DBRepoImpl) SubscribeSite(user_unique_id string, site_url string, is_subscribe bool) error {
	// まずはサイトURLをキーにサイトテーブルからサイトを検索する
	var site Site
	result := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return result.Error
	}
	// 次に対象のUserを検索する
	var user User
	result = DBMS.Where(&User{UserUniqueID: user_unique_id}).Find(&user)
	if result.Error != nil {
		return result.Error
	}
	// サブスクリブされていなかったらサブスクリプションサイトに追加する
	if !r.IsSubscribeSite(user_unique_id, site_url) && is_subscribe {
		subscriptionSite := SubscriptionSite{
			UserID: user.ID,
			SiteID: site.ID,
		}
		// トランザクション内で処理する
		DBMS.Transaction(func(tx *gorm.DB) error {
			result = tx.Create(&subscriptionSite)
			if result.Error != nil {
				log.Println("サブスクリプションサイト登録失敗 : " + result.Error.Error())
				// エラーが発生したらロールバックする
				tx.Rollback()
				return result.Error
			}
			return nil
		})
	} else if r.IsSubscribeSite(user_unique_id, site_url) && !is_subscribe {
		// サブスクリブされていたらサブスクリプションサイトから削除する
		DBMS.Transaction(func(tx *gorm.DB) error {
			result = tx.Where(&SubscriptionSite{UserID: user.ID, SiteID: site.ID}).Delete(&SubscriptionSite{})
			if result.Error != nil {
				log.Println("サブスクリプションサイト削除失敗 : " + result.Error.Error())
				// エラーが発生したらロールバックする
				tx.Rollback()
				return result.Error
			}
			return nil
		})
	}
	return nil
}

func (r DBRepoImpl) IsSubscribeSite(user_unique_id string, site_url string) bool {
	// まずはサイトURLをキーにサイトテーブルからサイトを検索する
	var site Site
	result := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return false
	}
	// Userのuser_idとUserの中のSubscriptionSiteのsite_urlをキーにSubscriptionSiteを検索する
	var res int64
	result = DBMS.Model(&User{}).Where(&User{UserUniqueID: user_unique_id}).Preload("SubscriptionSites", &SubscriptionSite{SiteID: site.ID}).Count(&res)
	if result.Error != nil || res == 0 {
		log.Println(result.Error)
		return false
	}
	if res > 0 {
		return true
	}
	return false
}

func (r DBRepoImpl) FetchSiteLastModified(site_url string) (time.Time, error) {
	// サイトURLをキーにサイトの最終更新日時だけを取得する
	var site Site
	result := DBMS.Model(&Site{}).Where(&Site{SiteUrl: site_url}).Select("LastModified").Find(&site)
	if result.Error != nil {
		return time.Now(), result.Error
	}
	return site.LastModified, nil
}

// 記事系処理
// 検索・最新記事取得・時間指定記事取得・記事更新

func (r DBRepoImpl) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	var articles []Article
	result := DBMS.Where("title LIKE ?", "%"+keyword+"%").Find(&articles)
	if result.Error != nil {
		return nil, result.Error
	}
	_, resultArticles := convertDbSiteToApi(Site{
		SiteArticles: articles,
	})
	return resultArticles, nil
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

func (r DBRepoImpl) SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isNew bool) ([]Data.Article, error) {
	var site Site
	result := DBMS.Where(&Site{SiteUrl: siteUrl}).Find(&site)
	if result.Error != nil {
		return nil, result.Error
	}
	var articles []Article
	if isNew {
		result = DBMS.Where(&Article{SiteID: site.ID}).Where("published_at BETWEEN ? AND ?", lastModified, time.Now().UTC()).Order("published_at desc").Limit(get_count).Find(&articles)
	} else {
		// 指定した時間より前の記事を取得する
		result = DBMS.Where(&Article{SiteID: site.ID}).Where("published_at BETWEEN ? AND ?", time.Time{}, lastModified).Order("published_at desc").Limit(get_count).Find(&articles)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	site.SiteArticles = articles
	_, resultArticles := convertDbSiteToApi(site)
	return resultArticles, nil
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
		if err := tx.Model(&siteModel).Update("LastModified",lastModified).Error; err != nil {
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
