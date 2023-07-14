package Repo

import (
	"errors"
	"fmt"
	"log"
	"read/Data"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テストを容易にするためDependency Injection（依存性の注入）を採用
// DBを呼び出す層はインターフェースを定義する
type DBRepository interface {
	// 全てのDBRepoに共通する処理であるDB接続を行う
	ConnectDB(isMock bool) error
	AutoMigrate() error
	// Readで使う
	GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error)
	FetchExploreCategories(country string) (resExp []Data.ExploreCategory, err error)

	// heavyで使う
	// 新規サイトをDB(サイトテーブル)に登録する
	RegisterSite(site Data.WebSite, articles []Data.Article) error
	// サイトURLをキーにサイトテーブルに該当するサイトがあるか確認する
	IsExistSite(site_url string) bool
	// ユーザーはサイトを購読しているか確認する
	IsSubscribeSite(user_unique_id string, site_url string) bool
	// サイトURLをキーにDBに該当するサイトを検索して返す
	SearchSiteByUrl(site_url string) (Data.WebSite, error)
	// サイト名をキーにサイトを検索
	SearchSiteByName(siteName string) ([]Data.WebSite, error)
	// キーワード検索でDBに該当する記事を返す
	SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
	// サイトURLをキーに記事更新日時を取得する
	FetchSiteLastModified(site_url string) (time.Time, error)
	// サイトURLをキーにサイトの最新記事を取得する
	SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error)
	// 指定された時間よりも新しいか古いを指定して記事を検索して配列を返す
	SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isNew bool) ([]Data.Article, error)
	// サイトの記事を更新する
	// サイトの更新日時より新しい記事があればDBに登録する
	UpdateArticles(siteUrl string, articles []Data.Article) error
	// サイトを購読登録する
	SubscribeSite(user_unique_id string, siteUrl string, is_subscribe bool) error

	// バッチ処理用
	// サイトテーブルを全件取得する
	FetchAllSites() ([]Data.WebSite, error)
	// 閲覧履歴テーブルを全件取得する
	FetchAllHistories() ([]Data.ReadActivity, error)
	// サイトと記事を大量に更新する
	UpdateSitesAndArticles(sites []Data.WebSite, articles []Data.Article) error
	// 時間（From・To）を指定してリードアクテビティを検索する
	SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadActivity, error)
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
		DBMS.AutoMigrate(&ClientConfig{})
		DBMS.AutoMigrate(&ApiActivity{})
		DBMS.AutoMigrate(&FavoriteSite{})
		DBMS.AutoMigrate(&SubscriptionSite{})
		DBMS.AutoMigrate(&SearchHistory{})
		DBMS.AutoMigrate(&ReadHistory{})
		DBMS.AutoMigrate(&ApiConfig{})
		DBMS.AutoMigrate(&UiConfig{})

		DBMS.AutoMigrate(&Site{})
		DBMS.AutoMigrate(&SiteArticle{})
		DBMS.AutoMigrate(&Tag{})
		DBMS.AutoMigrate(&ExploreCategory{})
	}

	return nil
}

func (s DBRepoImpl) GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error) {
	// これを使うかは疑問
	return Data.UserInfo{}, nil
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
		// サイトが存在しない場合は登録する
		dbSite := convertApiSiteToDb(site, articles)
		result := DBMS.Create(&dbSite)
		if result.Error != nil {
			return result.Error
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

	return Data.WebSite{}, nil
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
	user_id_int, err := strconv.Atoi(user_unique_id)
	if err != nil {
		return err
	}
	// まずはサイトURLをキーにサイトテーブルからサイトを検索する
	var site Site
	result := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return result.Error
	}
	// 次に対象のUserを検索する
	var user User
	result = DBMS.Where(&User{UserUniqueID: user_id_int}).Find(&user)
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
				return result.Error
			}
			return nil
		})
	} else if r.IsSubscribeSite(user_unique_id, site_url) && !is_subscribe {
		// サブスクリブされていたらサブスクリプションサイトから削除する
		result = DBMS.Where(&SubscriptionSite{UserID: user.ID, SiteID: site.ID}).Delete(&SubscriptionSite{})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (r DBRepoImpl) IsSubscribeSite(user_unique_id string, site_url string) bool {
	user_unique_id_int, err := strconv.Atoi(user_unique_id)
	if err != nil {
		return false
	}
	// まずはサイトURLをキーにサイトテーブルからサイトを検索する
	var site Site
	result := DBMS.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return false
	}
	// Userのuser_idとUserの中のSubscriptionSiteのsite_urlをキーにSubscriptionSiteを検索する
	var res int64
	result = DBMS.Model(&User{}).Where(&User{UserUniqueID: user_unique_id_int}).Preload("SubscriptionSites", &SubscriptionSite{SiteID: site.ID}).Count(&res)
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
	return time.Now(), nil
}

func (r DBRepoImpl) SearchSiteByName(siteName string) ([]Data.WebSite, error) {
	return []Data.WebSite{}, nil
}

// 記事系処理
// 検索・最新記事取得・時間指定記事取得・記事更新

func (r DBRepoImpl) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	return nil, nil
}

// サイトURLをキーにサイトの最新記事を取得する
func (r DBRepoImpl) SearchSiteLatestArticle(site_url string, get_count int) ([]Data.Article, error) {

	return nil, nil
}

func (r DBRepoImpl) SearchArticlesByTimeAndOrder(siteUrl string, lastModified time.Time, get_count int, isNew bool) ([]Data.Article, error) {
	return nil, nil
}

func (r DBRepoImpl) UpdateArticles(siteUrl string, articles []Data.Article) error {
	return nil
}

// バッチ処理用
func (r DBRepoImpl) FetchAllSites() ([]Data.WebSite, error) {
	return []Data.WebSite{}, nil
}

func (r DBRepoImpl) FetchAllHistories() ([]Data.ReadActivity, error) {
	return []Data.ReadActivity{}, nil
}

func (r DBRepoImpl) UpdateSitesAndArticles(sites []Data.WebSite, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SearchReadActivityByTime(from time.Time, to time.Time) ([]Data.ReadActivity, error) {
	return []Data.ReadActivity{}, nil
}
