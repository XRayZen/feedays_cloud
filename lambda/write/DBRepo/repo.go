package DBRepo

import (
	"read/Data"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBRepo はDBにアクセスするためのインターフェース
type DBRepo interface {
	// 全てのDBRepoに共通する処理であるDB接続を行う
	ConnectDB(isMock bool) error
	AutoMigrate() error
	SearchUserConfig(user_unique_Id string) (Data.UserConfig, error)
	RegisterUser(userInfo Data.UserConfig) error
	DeleteUser(user_unique_Id string) error
	AddReadActivity(user_unique_Id string, activityInfo Data.ReadActivity) error
	UpdateUser(user_unique_Id string, dataUserCfg Data.UserConfig) error
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
		DBMS.AutoMigrate(&ApiActivity{})
		DBMS.AutoMigrate(&FavoriteSite{})
		DBMS.AutoMigrate(&FavoriteArticle{})
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

func (repo DBRepoImpl) RegisterUser(userInfo Data.UserConfig) error {
	// API構造体からDB構造体に変換する
	dbApiCfg, dbUiCfg := ConvertToDbApiCfgAndUiCfg(userInfo, 0)
	user := User{
		UserName:     userInfo.UserName,
		UserUniqueID: userInfo.UserUniqueID,
		AccountType:  userInfo.AccountType,
		Country:      userInfo.Country,
		ApiConfig:    dbApiCfg,
		UiConfig:     dbUiCfg,
	}
	// DBに保存する
	if err := DBMS.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) SearchUserConfig(user_unique_Id string) (Data.UserConfig, error) {
	var user User
	// if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ClientConfig").Preload("ApiActivity").Preload("FavoriteSite").Preload("SubscriptionSite").Preload("ReadHistory").Preload("SearchHistory").First(&user).Error; err != nil {
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ApiConfig").Preload("UiConfig").First(&user).Error; err != nil {
		return Data.UserConfig{}, err
	}
	return ConvertToApiUserConfig(user), nil
}

func (repo DBRepoImpl) DeleteUser(user_unique_Id string) error {
	err := DBMS.Where("user_unique_id = ?", user_unique_Id).Delete(&User{})
	if err != nil {
		return err.Error
	}
	return nil
}

func (repo DBRepoImpl) UpdateUser(user_unique_Id string, dataUserCfg Data.UserConfig) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ApiConfig").Preload("UiConfig").First(&user).Error; err != nil {
		return err
	}
	// クライアント設定を更新する為だからそこら辺だけ更新する
	dbApiCfg, dbUiCfg := ConvertToDbApiCfgAndUiCfg(dataUserCfg, user.ID)
	// これだと更新出来ていないから一旦削除してから追加する
	// 古いUserを物理削除してから追加する
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Unscoped().Delete(&User{}).Error; err != nil {
		return err
	}
	// 変更を加えたUserを追加する事で更新とする
	user.UpdatedAt = time.Now()
	user.UserName = dataUserCfg.UserName
	user.UserUniqueID = dataUserCfg.UserUniqueID
	user.AccountType = dataUserCfg.AccountType
	user.Country = dataUserCfg.Country
	user.ApiConfig = dbApiCfg
	user.UiConfig = dbUiCfg
	if err := DBMS.Create(&user).Error; err != nil {
		return err
	}
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
