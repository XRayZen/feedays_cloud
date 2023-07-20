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
	GetUserConfig(user_unique_Id string) (Data.UserConfig, error)
	RegisterUser(userInfo Data.UserConfig) error
	DeleteUser(user_unique_Id string) error
	AddApiActivity(user_unique_Id string, activityInfo Data.Activity) error
	AddReadActivity(user_unique_Id string, activityInfo Data.ReadActivity) error
	UpdateClientConfig(user_unique_Id string, configInfo Data.UserConfig) error
	// 検索履歴を変更したら履歴を返す
	ModifySearchHistory(user_unique_Id string, text string, isAddOrRemove bool) ([]string, error)
	ModifyFavoriteSite(user_unique_Id string, siteInfo Data.WebSite, isAddOrRemove bool) error
	ModifyFavoriteArticle(user_unique_Id string, articleInfo Data.Article, isAddOrRemove bool) error
	GetAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error)
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

func (repo DBRepoImpl) GetUserConfig(user_unique_Id string) (Data.UserConfig, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ClientConfig").Preload("ApiActivity").Preload("FavoriteSite").Preload("SubscriptionSite").Preload("ReadHistory").Preload("SearchHistory").First(&user).Error; err != nil {
		return Data.UserConfig{}, err
	}
	return ConvertToApiUserConfig(user), nil
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

func (repo DBRepoImpl) DeleteUser(user_unique_Id string) error {
	err := DBMS.Where("user_unique_id = ?", user_unique_Id).Delete(&User{})
	if err != nil {
		return err.Error
	}
	return nil
}

func (repo DBRepoImpl) UpdateClientConfig(user_unique_Id string, configInfo Data.UserConfig) error {
	// ユーザー設定からクライアント設定IDを取得する
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ClientConfig").First(&user).Error; err != nil {
		return err
	}
	// クライアント設定を更新する為だからそこら辺だけ更新する
	user = ConvertToDbUserConfig(configInfo)
	// クライアント設定テーブルを更新する
	if err := DBMS.Model(&ClientConfig{}).Where("id = ?", user.ClientConfig.ID).Updates(user.ClientConfig).Error; err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) AddApiActivity(user_unique_Id string, activityInfo Data.Activity) error {
	// FirstOrCreateで存在しなかったら作成してインサートする
	return nil
}

func (repo DBRepoImpl) AddReadActivity(user_unique_Id string, activityInfo Data.ReadActivity) error {
	return nil
}

func (repo DBRepoImpl) ModifySearchHistory(user_unique_Id string, text string, isAddOrRemove bool) ([]string, error) {
	return []string{}, nil
}

func (repo DBRepoImpl) ModifyFavoriteSite(user_unique_Id string, siteInfo Data.WebSite, isAddOrRemove bool) error {
	return nil
}

func (repo DBRepoImpl) ModifyFavoriteArticle(user_unique_Id string, articleInfo Data.Article, isAddOrRemove bool) error {
	return nil
}

func (repo DBRepoImpl) GetAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error) {
	return Data.ApiConfig{}, nil
}
