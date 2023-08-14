package Repo

import (
	"errors"
	"fmt"
	"log"
	"site/Data"
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
	// Siteで使う
	IsExistSite(site_url string) bool
	RegisterSite(site Data.WebSite, articles []Data.Article) error
	SubscribeSite(user_unique_id string, site_url string, is_subscribe bool) error
	IsSubscribeSite(user_unique_id string, site_url string) bool
	FetchSiteLastModified(site_url string) (time.Time, error)
	SearchSiteByUrl(site_url string) (Data.WebSite, error)
	SearchSiteByName(siteName string) ([]Data.WebSite, error)
	SearchSiteByCategory(category string) ([]Data.WebSite, error)
	SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error)
	SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isNew bool) ([]Data.Article, error)
	SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
	ChangeSiteCategory(user_unique_id string, site_url string, category_name string) error
	DeleteSite(site_url string) error
	DeleteSiteByUnscoped(site_url string) error
	ModifyExploreCategory(category Data.ExploreCategory, is_add_or_remove bool) error
}

// DBRepoを実装
type DBRepoImpl struct {
}

// DB接続
func (s DBRepoImpl) ConnectDB(isMock bool) error {
	// DB接続
	if isMock {
		// もしモックモードが有効ならSqliteに接続する
		in_memory_str := "file::memory:"
		// fileSqlStr := "test_1.db"
		DB, err := gorm.Open(sqlite.Open(in_memory_str))
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
	var explore_categories []ExploreCategory
	result := DBMS.Where(&ExploreCategory{Country: country}).Find(&explore_categories)
	if result.Error != nil {
		return nil, result.Error
	}
	// カテゴリーをExploreCategories型に変換する
	var categories []Data.ExploreCategory
	for _, category := range explore_categories {
		categories = append(categories, Data.ExploreCategory{
			CategoryName:        category.CategoryName,
			CategoryDescription: category.Description,
			CategoryID:          fmt.Sprint(category.ID),
		})
	}
	return categories, nil
}

// サイト系処理
// 登録・存在確認・購読・購読確認・URL検索
func (r DBRepoImpl) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	// サイトを登録する
	// もし、サイトが存在していたら登録しない
	if !r.IsExistSite(site.SiteURL) {
		db_site := convertApiSiteToDb(site, articles)
		if err := DBMS.Transaction(func(tx *gorm.DB) error {
			result := tx.Create(&db_site)
			if result.Error != nil {
				log.Println("サイト登録失敗 : " + result.Error.Error())
				tx.Rollback()
				return result.Error
			}
			return nil
		}); err != nil {
			return err
		}
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
	result_site, _ := convertDbSiteToApi(site)
	return result_site, nil
}

func (r DBRepoImpl) SearchSiteByName(siteName string) ([]Data.WebSite, error) {
	var sites []Site
	result := DBMS.Where(&Site{SiteName: siteName}).Find(&sites)
	if result.Error != nil {
		return []Data.WebSite{}, result.Error
	}
	var result_sites []Data.WebSite
	for _, site := range sites {
		result_site, _ := convertDbSiteToApi(site)
		result_sites = append(result_sites, result_site)
	}
	return result_sites, nil
}

func (r DBRepoImpl) SearchSiteByCategory(category string) ([]Data.WebSite, error) {
	var sites []Site
	if err := DBMS.Where(&Site{Category: category}).Find(&sites).Error; err != nil {
		return []Data.WebSite{}, err
	}
	var result_sites []Data.WebSite
	for _, site := range sites {
		result_site, _ := convertDbSiteToApi(site)
		result_sites = append(result_sites, result_site)
	}
	return result_sites, nil
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
	if err := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site); err != nil {
		return err.Error
	}
	// 次に対象のUserを検索する
	var user User
	if err := DBMS.Where(&User{UserUniqueID: user_unique_id}).Find(&user); err != nil {
		return err.Error
	}
	// サブスクリブされていなかったらサブスクリプションサイトに追加する
	if !r.IsSubscribeSite(user_unique_id, site_url) && is_subscribe {
		subscription_site := SubscriptionSite{
			UserID: user.ID,
			SiteID: site.ID,
		}
		// トランザクション内で処理する
		if err := DBMS.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&subscription_site).Error; err != nil {
				log.Println("サブスクリプションサイト登録失敗 : " + err.Error())
				// エラーが発生したらロールバックする
				tx.Rollback()
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	} else if r.IsSubscribeSite(user_unique_id, site_url) && !is_subscribe {
		// サブスクリブされていたらサブスクリプションサイトから削除する
		if err := DBMS.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(&SubscriptionSite{UserID: user.ID, SiteID: site.ID}).Delete(&SubscriptionSite{}).Error; err != nil {
				log.Println("サブスクリプションサイト削除失敗 : " + err.Error())
				// エラーが発生したらロールバックする
				tx.Rollback()
				return err
			}
			return nil
		}); err != nil {
			return err
		}
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
	_, result_articles := convertDbSiteToApi(Site{
		SiteArticles: articles,
	})
	return result_articles, nil
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
	_, result_articles := convertDbSiteToApi(site)
	return result_articles, nil
}

// サイトのカテゴリを変更する
func (r DBRepoImpl) ChangeSiteCategory(user_unique_id string, site_url string, category_name string) error {
	// サイトURLをキーにサイトを取得する
	var site Site
	if err := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site).Error; err != nil {
		return err
	}
	// サイトのカテゴリを変更する
	if err := DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&site).Update("category", category_name).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// サイトを削除する
// DeleteSite(user_unique_id string, site_url string) error
func (r DBRepoImpl) DeleteSite(site_url string) error {
	// サイトURLをキーにサイトを取得する
	var site Site
	if err := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site).Error; err != nil {
		return err
	}
	// サイトを削除する
	if err := DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&site).Error; err != nil {
			tx.Rollback()
			return err
		}
		// SiteArticlesを削除する
		if err := tx.Where(&Article{SiteID: site.ID}).Delete(&Article{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// Tagを削除する
		if err := tx.Where(&Tag{SiteID: site.ID}).Delete(&Tag{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// DeleteSiteByUnscoped(site_url string) error
func (r DBRepoImpl) DeleteSiteByUnscoped(site_url string) error {
	// サイトURLをキーにサイトを取得する
	var site Site
	if err := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site).Error; err != nil {
		return err
	}
	// サイトを削除する
	if err := DBMS.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Delete(&site).Error; err != nil {
			tx.Rollback()
			return err
		}
		// SiteArticlesを削除する
		if err := tx.Where(&Article{SiteID: site.ID}).Unscoped().Delete(&Article{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		// Tagを削除する
		if err := tx.Where(&Tag{SiteID: site.ID}).Unscoped().Delete(&Tag{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// ModifyExploreCategory(category Data.ExploreCategory, is_add_or_remove bool) error
func (r DBRepoImpl) ModifyExploreCategory(category Data.ExploreCategory, is_add_or_remove bool) error {
	if err := DBMS.Transaction(func(tx *gorm.DB) error {
		explore_category := ExploreCategory{
			CategoryName: category.CategoryName,
			Description:  category.CategoryDescription,
			Country:      category.CategoryCountry,
			image_url:    category.CategoryImage,
		}
		if is_add_or_remove {
			// カテゴリを追加する
			if err := tx.Create(&explore_category).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// カテゴリを削除する
			if err := tx.Where(&explore_category).Delete(&explore_category).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
