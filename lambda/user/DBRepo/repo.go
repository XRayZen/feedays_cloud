package DBRepo

import (
	"log"
	"sort"
	"time"
	"user/Data"

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
	ModifyFavoriteSite(user_unique_Id string, siteUrl string, isAddOrRemove bool) error
	ModifyFavoriteArticle(user_unique_Id string, articleUrl string, isAddOrRemove bool) error
	FetchAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error)
	UpdateAPIRequestLimit(user_unique_Id string, apiConfig Data.ApiConfig) error
	DeleteUserData(user_unique_Id string) error
	DeletesUnscopedUserData(user_unique_Id string) error
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
		DBMS.AutoMigrate(&Article{})
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
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("ApiConfig").Preload("UiConfig").First(&user).Error; err != nil {
		return Data.UserConfig{}, err
	}
	// AssociationでReadHistoryなどの配列を取得する
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
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return []string{}, err
	}
	if isAddOrRemove {
		// 追加
		dbSearchHist := SearchHistory{
			SearchWord: text,
			searchAt:   time.Now(),
		}
		err := DBMS.Model(&user).Association("SearchHistories").Append(&dbSearchHist)
		if err != nil {
			return []string{}, err
		}
	} else {
		// 今のUserテーブルとの参照を削除
		// 参照先のテーブルからは削除されない
		err := DBMS.Model(&user).Association("SearchHistories").Delete(&SearchHistory{SearchWord: text})
		if err != nil {
			return []string{}, err
		}
		// ReadHistoriesテーブルからも削除する
		err = DBMS.Where("search_word = ?", text).Delete(&SearchHistory{}).Error
		if err != nil {
			return []string{}, err
		}
	}
	// 再度取得する
	var dbSearchHist []SearchHistory
	if err := DBMS.Model(&user).Association("SearchHistories").Find(&dbSearchHist); err != nil {
		return []string{}, err
	}
	// SearchAtでdescソートしてdbSearchHistに入れる
	sort.Slice(dbSearchHist, func(i, j int) bool {
		return dbSearchHist[i].searchAt.After(dbSearchHist[j].searchAt)
	})
	var apiSearchHists []string
	for _, dbSearchHist := range dbSearchHist {
		apiSearchHists = append(apiSearchHists, dbSearchHist.SearchWord)
	}
	return apiSearchHists, nil
}

func (repo DBRepoImpl) ModifyFavoriteSite(user_unique_Id string, siteUrl string, isAddOrRemove bool) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	// Siteを取得する
	var site Site
	if err := DBMS.Where("site_url = ?", siteUrl).First(&site).Error; err != nil {
		return err
	}
	if isAddOrRemove {
		// 追加
		if err := DBMS.Model(&user).Association("FavoriteSites").Append(&FavoriteSite{SiteID: site.ID}); err != nil {
			return err
		}
	} else {
		// 削除
		if err := DBMS.Model(&user).Association("FavoriteSites").Delete(&FavoriteSite{SiteID: site.ID}); err != nil {
			return err
		}
		// FavoriteSitesテーブルからも削除する
		if err := DBMS.Where("site_id = ?", site.ID).Delete(&FavoriteSite{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repo DBRepoImpl) ModifyFavoriteArticle(user_unique_Id string, articleUrl string, isAddOrRemove bool) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	// Articleを取得する
	var article Article
	if err := DBMS.Where("url = ?", articleUrl).First(&article).Error; err != nil {
		return err
	}
	if isAddOrRemove {
		// 追加
		if err := DBMS.Model(&user).Association("FavoriteArticles").Append(&FavoriteArticle{ArticleID: article.ID}); err != nil {
			return err
		}
	} else {
		// 削除
		if err := DBMS.Model(&user).Association("FavoriteArticles").Delete(&FavoriteArticle{ArticleID: article.ID}); err != nil {
			return err
		}
	}
	return nil
}

func (repo DBRepoImpl) FetchAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return Data.ApiConfig{}, err
	}
	var apiConfig ApiConfig
	if err := DBMS.Model(&user).Association("ApiConfig").Find(&apiConfig); err != nil {
		return Data.ApiConfig{}, err
	}
	user.ApiConfig = apiConfig
	// DB型からAPI型に変換する
	return Data.ApiConfig{
		RefreshArticleInterval:      user.ApiConfig.RefreshArticleInterval,
		FetchArticleRequestInterval: user.ApiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    user.ApiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   user.ApiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      user.ApiConfig.FetchTrendRequestLimit,
	}, nil
}

// UpdateAPIRequestLimit(user_unique_Id string, apiConfig Data.ApiConfig) error
func (repo DBRepoImpl) UpdateAPIRequestLimit(user_unique_Id string, apiConfig Data.ApiConfig) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	// 更新する
	apiCfg := ApiConfig{
		RefreshArticleInterval:      apiConfig.RefreshArticleInterval,
		FetchArticleRequestInterval: apiConfig.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    apiConfig.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   apiConfig.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      apiConfig.FetchTrendRequestLimit,
	}
	if err := DBMS.Model(&user).Association("ApiConfig").Replace(&apiCfg); err != nil {
		return err
	}
	return nil
}

// DeleteUserData(user_unique_Id string) error
func (repo DBRepoImpl) DeleteUserData(user_unique_Id string) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	// Userテーブルから削除
	if err := DBMS.Delete(&user).Error; err != nil {
		return err
	}
	// FavoriteSitesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Delete(&FavoriteSite{}).Error; err != nil {
		return err
	}
	// FavoriteArticlesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Delete(&FavoriteArticle{}).Error; err != nil {
		return err
	}
	// SearchHistoriesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Delete(&SearchHistory{}).Error; err != nil {
		return err
	}
	return nil
}

// DeletesUnscopedUserData(user_unique_Id string) error
func (repo DBRepoImpl) DeletesUnscopedUserData(user_unique_Id string) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	// ApiConfigテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&ApiConfig{}).Error; err != nil {
		return err
	}
	// UiConfigテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&UiConfig{}).Error; err != nil {
		return err
	}
	// Userテーブルから削除
	if err := DBMS.Unscoped().Delete(&user).Error; err != nil {
		return err
	}
	// FavoriteSitesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&FavoriteSite{}).Error; err != nil {
		return err
	}
	// FavoriteArticlesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&FavoriteArticle{}).Error; err != nil {
		return err
	}
	// SearchHistoriesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&SearchHistory{}).Error; err != nil {
		return err
	}
	// ReadHistoriesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&ReadHistory{}).Error; err != nil {
		return err
	}
	// SubscriptionSitesテーブルから削除
	if err := DBMS.Where("user_id = ?", user.ID).Unscoped().Delete(&SubscriptionSite{}).Error; err != nil {
		return err
	}
	return nil
}
