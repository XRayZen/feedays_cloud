package Repo

import (
	"fmt"
	"read/Data"
	"strconv"
	"time"
)

// テストを容易にするためDependency Injection（依存性の注入）を採用
// DBを呼び出す層はインターフェースを定義する
type DBRepository interface {
	// Readで使う
	GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error)
	FetchRanking(useeID string) (resRanking Data.Ranking, err error)
	FetchExploreCategories(country string) (resExp Data.ExploreCategories, err error)

	// heavyで使う
	// サイトURLをキーにサイトテーブルに該当するサイトがあるか確認する
	IsExistSite(site_url string) bool
	// ユーザーはサイトを購読しているか確認する
	IsSubscribeSite(user_id string, site_url string) bool
	// サイトURLをキーにDBに該当するサイトを検索して返す
	FetchSite(site_url string) (Data.WebSite, error)
	// サイトURLをキーに記事更新チェック日時を取得する
	FetchSiteLastModified(site_url string) (time.Time, error)
	// 新規サイトをDB(サイトテーブル)に登録する
	RegisterSite(site Data.WebSite, articles []Data.Article) error
	// キーワード検索でDBに該当する記事を返す
	SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
	// サイトの記事を指定した日時より新しい記事を返す
	SearchArticlesByTime(siteUrl string, lastModified time.Time) ([]Data.Article, error)
	// サイト名をキーにサイトを検索
	SearchSiteByName(siteName string) ([]Data.WebSite, error)
	// サイトの記事を更新する
	// サイトの更新日時より新しい記事があればDBに登録する
	UpdateArticles(siteUrl string, articles []Data.Article) error
	// サイトを購読登録する
	SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error

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

// DBRepoのリアルを実装
type DBRepoImpl struct {
}

func (s DBRepoImpl) GetUserInfo(userID string) (resUserInfo Data.UserInfo, err error) {
	// これを使うかは疑問
	return Data.UserInfo{}, nil
}

func (s DBRepoImpl) FetchRanking(userID string) (resRanking Data.Ranking, err error) {
	// ユーザーIDからユーザーの国を取得してDBから国ごとのランキングを取得する
	db, err := Connect()
	if err != nil {
		return Data.Ranking{}, err
	}
	// ユーザーIDをintに変換する
	num, err := strconv.Atoi(userID)
	if err != nil {
		return Data.Ranking{}, err
	}
	// ユーザーIDからユーザーの国を取得する
	var user User
	db.Where(&User{UserUniqueID: num}).Select("country").Find(&user)
	// ユーザーの国をキーにDBからランキングを取得する
	// TODO:これ以降の作業を進むにはDBユーザーに国を追加定義する必要がある
	// var ranking SiteRanking
	// db.Where(&SiteRanking{country: user.}).Find(&ranking)

	return Data.Ranking{}, nil
}

func (s DBRepoImpl) FetchExploreCategories(country string) (res []Data.ExploreCategories, err error) {
	// ユーザーIDと国をキーにDBからカテゴリーを取得する
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	// ExploreCategoriesテーブルから国をキーにカテゴリーを全件取得する
	var expCats []ExploreCategory
	result := db.Where(&ExploreCategory{Country: country}).Find(&expCats)
	if result.Error != nil {
		return nil, result.Error
	}
	// カテゴリーをExploreCategories型に変換する
	var categories []Data.ExploreCategories
	for _, expCat := range expCats {
		categories = append(categories, Data.ExploreCategories{
			CategoryName:        expCat.CategoryName,
			CategoryDescription: expCat.Description,
			CategoryID:          fmt.Sprint(expCat.ID),
		})
	}
	return categories, nil
}

// heavyで使う
func (r DBRepoImpl) IsExistSite(site_url string) bool {
	db, err := Connect()
	if err != nil {
		return false
	}
	var site Site
	result := db.Where(&Site{SiteUrl: site_url}).Find(&site)
	if result.Error != nil {
		return false
	}
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func (r DBRepoImpl) IsSubscribeSite(user_id string, site_url string) bool {
	db, err := Connect()
	if err != nil {
		return false
	}
	user_id_int, err := strconv.Atoi(user_id)
	if err != nil {
		return false
	}
	// Userのuser_idとUserの中のSubscriptionSiteのsite_urlをキーにSubscriptionSiteを検索する
	var res SubscriptionSite
	result := db.Where(&User{UserUniqueID: user_id_int}).Where(&SubscriptionSite{site_url: site_url}).Find(&res)
	// result := db.Model(&User{}).Where("user_id = ?", user_id_int).Where(&SubscriptionSite{site_url: site_url}).Find(&res)
	if result.Error != nil {
		return false
	}
	if result.RowsAffected > 0 {
		return true
	}
	return false
}

func (r DBRepoImpl) FetchSite(site_url string) (Data.WebSite, error) {
	return Data.WebSite{}, nil
}

func (r DBRepoImpl) FetchSiteLastModified(site_url string) (time.Time, error) {
	return time.Now(), nil
}

func (r DBRepoImpl) RegisterSite(site Data.WebSite, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
	return nil, nil
}

func (r DBRepoImpl) SearchArticlesByTime(siteUrl string, lastModified time.Time) ([]Data.Article, error) {
	return nil, nil
}

func (r DBRepoImpl) SearchSiteByName(siteName string) ([]Data.WebSite, error) {
	return []Data.WebSite{}, nil
}

func (r DBRepoImpl) UpdateArticles(siteUrl string, articles []Data.Article) error {
	return nil
}

func (r DBRepoImpl) SubscribeSite(user_id string, siteUrl string, is_subscribe bool) error {
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
