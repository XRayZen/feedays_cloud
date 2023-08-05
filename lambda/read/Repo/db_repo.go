package Repo

import (
	"fmt"
	"read/Data"

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
	SearchUserConfig(user_unique_Id string,isPreloadRelatedTables bool) (Data.UserConfig, error)
	FetchExploreCategories(country string) (resExp []Data.ExploreCategory, err error)

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
