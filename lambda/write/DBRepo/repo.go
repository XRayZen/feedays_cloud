package DBRepo

import (
	"log"
	"read/Data"
	// "time"

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
	UpdateUser(user_unique_Id string, dataUserCfg Data.UserConfig) error

	AddReadHistory(user_unique_Id string, activityInfo Data.ReadHistory) error
	SearchReadHistory(user_unique_Id string, limit int) ([]Data.ReadHistory, error)
	// 検索履歴を変更したら履歴を返す
	ModifySearchHistory(user_unique_Id string, word string, isAddOrRemove bool) ([]string, error)
	ModifyFavoriteSite(user_unique_Id string, siteInfo Data.WebSite, isAddOrRemove bool) error
	ModifyFavoriteArticle(user_unique_Id string, articleInfo Data.Article, isAddOrRemove bool) error
	GetAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error)
}

type DBRepoImpl struct {
}

func (repo DBRepoImpl) ConnectDB(isMock bool) error {
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
	} else {
		// もしモックモードが無効ならLambdaの環境変数が指定するDBに接続する
		if err := DataBaseConnect(); err != nil {
			return err
		}
	}
	return nil
}

func (repo DBRepoImpl) AutoMigrate() error {
	// DB接続
	if DBMS != nil {
		DBMS.AutoMigrate(&User{})
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
	dbApiCfg, dbUiCfg := ConvertToDbApiCfgAndUiCfg(userInfo)
	user := User{
		UserName:      userInfo.UserName,
		UserUniqueID:  userInfo.UserUniqueID,
		AccountType:   userInfo.AccountType,
		Country:       userInfo.Country,
		ApiConfig:     dbApiCfg,
		UiConfig:      dbUiCfg,
		ReadHistories: []ReadHistory{},
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
	// AssociationでReadHistoryなどの配列を取得する
	if err := DBMS.Model(&user).Association("ReadHistories").Find(&user.ReadHistories); err != nil {
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
	dbApiCfg, dbUiCfg := ConvertToDbApiCfgAndUiCfg(dataUserCfg)
	// Associationでリプレースする
	if err := DBMS.Model(&user).Association("ApiConfig").Replace(&dbApiCfg); err != nil {
		return err
	}
	if err := DBMS.Model(&user).Association("UiConfig").Replace(&dbUiCfg); err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) AddReadHistory(user_unique_Id string, readHst Data.ReadHistory) error {
	// DB型に変換する
	dbReadHist, err := ConvertToDbReadHistory(readHst)
	if err != nil {
		return err
	}
	// Userを取得する
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	err = DBMS.Model(&user).Association("ReadHistories").Append(&dbReadHist)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (repo DBRepoImpl) SearchReadHistory(user_unique_Id string, limit int) ([]Data.ReadHistory, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return []Data.ReadHistory{}, err
	}
	var dbReadHist []ReadHistory
	if err := DBMS.Where("user_id = ?", user.ID).Order("created_at desc").Limit(limit).Find(&dbReadHist).Error; err != nil {
		return []Data.ReadHistory{}, err
	}
	var apiHists []Data.ReadHistory
	for _, dbReadHist := range dbReadHist {
		apiHists = append(apiHists, ConvertToApiReadHistory(dbReadHist, user.UserUniqueID))
	}
	return apiHists, nil
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
