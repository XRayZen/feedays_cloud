package DBRepo

import (
	"read/Data"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBRepo はDBにアクセスするためのインターフェース
type DBRepo interface {
	// 全てのDBRepoに共通する処理であるDB接続を行う
	ConnectDB(isMock bool) error
	AutoMigrate() error
	GetUserConfig(userId string) (Data.UserConfig, error)
	RegisterUser(userInfo Data.UserConfig) error
	DeleteUser(userId string) error
	AddApiActivity(userId string, activityInfo Data.Activity) error
	AddReadActivity(userId string, activityInfo Data.ReadActivity) error
	UpdateAppConfig(userId string, configInfo Data.UserConfig) error
	// 検索履歴を変更したら履歴を返す
	ModifySearchHistory(userId string, text string, isAddOrRemove bool) ([]string, error)
	ModifyFavoriteSite(userId string, siteInfo Data.WebSite, isAddOrRemove bool) error
	ModifyFavoriteArticle(userId string, articleInfo Data.Article, isAddOrRemove bool) error
	GetAPIRequestLimit(userId string) (Data.ApiRequestLimitConfig, error)
}

type DBRepoImpl struct {
}

func (repo DBRepoImpl) ConnectDB(isMock bool) error {
	// DB接続
	if isMock {
		// もしモックモードが有効ならSqliteに接続する
		InMemoryStr := "file::memory:"
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

func (repo DBRepoImpl) AutoMigrate() error {
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

func (repo DBRepoImpl) GetUserConfig(userId string) (Data.UserConfig, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", userId).Preload("ClientConfig").Preload("ApiActivity").Preload("FavoriteSite").Preload("SubscriptionSite").Preload("ReadHistory").Preload("SearchHistory").First(&user).Error; err != nil {
		return Data.UserConfig{}, err
	}
	return ConvertToUserConfig(user), nil
}

func (repo DBRepoImpl) RegisterUser(userInfo Data.UserConfig) error {
	// API構造体からDB構造体に変換する
	user := ConvertToDbUserConfig(userInfo)
	// DBに保存する
	if err := DBMS.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) DeleteUser(userId string) error {
	err := DBMS.Where("user_unique_id = ?", userId).Delete(&User{})
	if err != nil {
		return err.Error
	}
	return nil
}

func (repo DBRepoImpl) UpdateAppConfig(userId string, configInfo Data.UserConfig) error {
	// UI設定などを更新する為だからそこら辺だけ更新する
	user := ConvertToDbUserConfig(configInfo)
	if err := DBMS.Model(&User{}).Where("user_unique_id = ?", userId).Updates(&User{
		ClientConfig: user.ClientConfig,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) AddApiActivity(userId string, activityInfo Data.Activity) error {
	// FirstOrCreateで存在しなかったら作成してインサートする
	return nil
}

func (repo DBRepoImpl) AddReadActivity(userId string, activityInfo Data.ReadActivity) error {
	return nil
}

func (repo DBRepoImpl) ModifySearchHistory(userId string, text string, isAddOrRemove bool) ([]string, error) {
	return []string{}, nil
}

func (repo DBRepoImpl) ModifyFavoriteSite(userId string, siteInfo Data.WebSite, isAddOrRemove bool) error {
	return nil
}

func (repo DBRepoImpl) ModifyFavoriteArticle(userId string, articleInfo Data.Article, isAddOrRemove bool) error {
	return nil
}

func (repo DBRepoImpl) GetAPIRequestLimit(userId string) (Data.ApiRequestLimitConfig, error) {
	return Data.ApiRequestLimitConfig{}, nil
}
