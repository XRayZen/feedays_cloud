package DbRepo

import (
	"errors"
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
	ConnectDB(is_mock bool) error
	AutoMigrate() error
	DropTable() error

	SearchUserConfig(user_unique_Id string) (Data.UserConfig, error)
	RegisterUser(user_info Data.UserConfig) error
	DeleteUser(user_unique_Id string) error
	UpdateUserUiConfig(user_unique_Id string, data_user_config Data.UserConfig) error

	AddReadHistory(user_unique_Id string, activity_info Data.ReadHistory) error
	SearchReadHistory(user_unique_Id string, limit int) ([]Data.ReadHistory, error)
	// 検索履歴を変更したら履歴を返す
	ModifySearchHistory(user_unique_Id string, word string, is_add_or_remove bool) ([]string, error)
	ModifyFavoriteSite(user_unique_Id string, siteUrl string, is_add_or_remove bool) error
	ModifyFavoriteArticle(user_unique_Id string, articleUrl string, is_add_or_remove bool) error
	FetchAPIRequestLimit(user_unique_Id string) (Data.ApiConfig, error)
	ModifyApiRequestLimit(modify_type string, api_config Data.ApiConfig) error
	DeleteUserData(user_unique_Id string) error
	DeletesUnscopedUserData(user_unique_Id string) error
	ModifyExploreCategory(modify_type string, category Data.ExploreCategory) error
}

type DBRepoImpl struct {
}

func (repo DBRepoImpl) ConnectDB(isMock bool) error {
	if isMock {
		// もしモックモードが有効ならSqliteに接続する
		in_memory_str := "file::memory:"
		DB, err := gorm.Open(sqlite.Open(in_memory_str))
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
		DBMS.AutoMigrate(&ApiLimitConfig{})
		DBMS.AutoMigrate(&UiConfig{})

		DBMS.AutoMigrate(&Site{})
		DBMS.AutoMigrate(&Article{})
		DBMS.AutoMigrate(&Tag{})
		DBMS.AutoMigrate(&ExploreCategory{})
	}
	return nil
}

func (repo DBRepoImpl) DropTable() error {
	if DBMS != nil {
		DBMS.Migrator().DropTable(&User{})
		DBMS.Migrator().DropTable(&FavoriteSite{})
		DBMS.Migrator().DropTable(&FavoriteArticle{})
		DBMS.Migrator().DropTable(&SubscriptionSite{})
		DBMS.Migrator().DropTable(&SearchHistory{})
		DBMS.Migrator().DropTable(&ReadHistory{})
		DBMS.Migrator().DropTable(&ApiLimitConfig{})
		DBMS.Migrator().DropTable(&UiConfig{})

		DBMS.Migrator().DropTable(&Site{})
		DBMS.Migrator().DropTable(&Article{})
		DBMS.Migrator().DropTable(&Tag{})
		DBMS.Migrator().DropTable(&ExploreCategory{})
	}
	return nil
}

func (repo DBRepoImpl) RegisterUser(userInfo Data.UserConfig) error {
	// API構造体からDB構造体に変換する
	db_ui_config := ConvertToDbUiCfg(userInfo)
	user := User{
		UserName:      userInfo.UserName,
		UserUniqueID:  userInfo.UserUniqueID,
		AccountType:   userInfo.AccountType,
		Country:       userInfo.Country,
		UiConfig:      db_ui_config,
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
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("UiConfig").First(&user).Error; err != nil {
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

func (repo DBRepoImpl) UpdateUserUiConfig(user_unique_Id string, dataUserCfg Data.UserConfig) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).Preload("UiConfig").First(&user).Error; err != nil {
		return err
	}
	// クライアント設定を更新する為だからそこら辺だけ更新する
	db_ui_config := ConvertToDbUiCfg(dataUserCfg)
	if err := DBMS.Model(&user).Association("UiConfig").Replace(&db_ui_config); err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) AddReadHistory(user_unique_Id string, readHst Data.ReadHistory) error {
	// DB型に変換する
	db_read_History, err := ConvertToDbReadHistory(readHst)
	if err != nil {
		return err
	}
	// Userを取得する
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	if err = DBMS.Model(&user).Association("ReadHistories").Append(&db_read_History); err != nil {
		return err
	}
	return nil
}

func (repo DBRepoImpl) SearchReadHistory(user_unique_Id string, limit int) ([]Data.ReadHistory, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return []Data.ReadHistory{}, err
	}
	var db_read_histories []ReadHistory
	if err := DBMS.Where("user_id = ?", user.ID).Order("created_at desc").Limit(limit).Find(&db_read_histories).Error; err != nil {
		return []Data.ReadHistory{}, err
	}
	var api_histories []Data.ReadHistory
	for _, db_read_history := range db_read_histories {
		api_histories = append(api_histories, ConvertToApiReadHistory(db_read_history, user.UserUniqueID))
	}
	return api_histories, nil
}

func (repo DBRepoImpl) ModifySearchHistory(user_unique_Id string, text string, is_add_or_remove bool) ([]string, error) {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return []string{}, err
	}
	if is_add_or_remove {
		// 追加
		db_search_history := SearchHistory{
			SearchWord: text,
			searchAt:   time.Now(),
		}
		err := DBMS.Model(&user).Association("SearchHistories").Append(&db_search_history)
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
	var db_search_history []SearchHistory
	if err := DBMS.Model(&user).Association("SearchHistories").Find(&db_search_history); err != nil {
		return []string{}, err
	}
	// SearchAtでdescソートしてdbSearchHistに入れる
	sort.Slice(db_search_history, func(i, j int) bool {
		return db_search_history[i].searchAt.After(db_search_history[j].searchAt)
	})
	var api_search_histories []string
	for _, db_search_history := range db_search_history {
		api_search_histories = append(api_search_histories, db_search_history.SearchWord)
	}
	return api_search_histories, nil
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

func (repo DBRepoImpl) ModifyFavoriteArticle(user_unique_Id string, article_url string, is_add_or_remove bool) error {
	var user User
	if err := DBMS.Where("user_unique_id = ?", user_unique_Id).First(&user).Error; err != nil {
		return err
	}
	// Articleを取得する
	var article Article
	if err := DBMS.Where("url = ?", article_url).First(&article).Error; err != nil {
		return err
	}
	if is_add_or_remove {
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
	// APIリクエスト制限はアカウントタイプごとに設定されてテーブルに保存されている
	var api_config ApiLimitConfig
	if err := DBMS.Where(&ApiLimitConfig{AccountType: user.AccountType}).First(&api_config).Error; err != nil {
		return Data.ApiConfig{}, err
	}
	// DB型からAPI型に変換する
	return Data.ApiConfig{
		AccountType:                 api_config.AccountType,
		RefreshArticleInterval:      api_config.RefreshArticleInterval,
		FetchArticleRequestInterval: api_config.FetchArticleRequestInterval,
		FetchArticleRequestLimit:    api_config.FetchArticleRequestLimit,
		FetchTrendRequestInterval:   api_config.FetchTrendRequestInterval,
		FetchTrendRequestLimit:      api_config.FetchTrendRequestLimit,
	}, nil
}

func (repo DBRepoImpl) ModifyApiRequestLimit(modify_type string, api_config Data.ApiConfig) error {
	if err := DBMS.Transaction(func(tx *gorm.DB) error {
		switch modify_type {
		case "Add":
			// 追加
			db_api_config := ApiLimitConfig{
				AccountType:                 api_config.AccountType,
				RefreshArticleInterval:      api_config.RefreshArticleInterval,
				FetchArticleRequestInterval: api_config.FetchArticleRequestInterval,
				FetchArticleRequestLimit:    api_config.FetchArticleRequestLimit,
				FetchTrendRequestInterval:   api_config.FetchTrendRequestInterval,
				FetchTrendRequestLimit:      api_config.FetchTrendRequestLimit,
			}
			if err := tx.Create(&db_api_config).Error; err != nil {
				tx.Rollback()
				return err
			}
		case "Update":
			// 更新する
			db_api_config := ApiLimitConfig{
				AccountType:                 api_config.AccountType,
				RefreshArticleInterval:      api_config.RefreshArticleInterval,
				FetchArticleRequestInterval: api_config.FetchArticleRequestInterval,
				FetchArticleRequestLimit:    api_config.FetchArticleRequestLimit,
				FetchTrendRequestInterval:   api_config.FetchTrendRequestInterval,
				FetchTrendRequestLimit:      api_config.FetchTrendRequestLimit,
			}
			if err := tx.Model(&ApiLimitConfig{}).Where(&ApiLimitConfig{AccountType: api_config.AccountType}).Updates(&db_api_config).Error; err != nil {
				tx.Rollback()
				return err
			}
		case "Delete":
			// 削除する
			if err := tx.Where(&ApiLimitConfig{AccountType: api_config.AccountType}).Delete(&ApiLimitConfig{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		case "UnscopedDelete":
			// 物理的に削除する
			if err := tx.Where(&ApiLimitConfig{AccountType: api_config.AccountType}).Unscoped().Delete(&ApiLimitConfig{}).Error; err != nil {
				tx.Rollback()
				return err
			}
		default:
			// エラーを返す
			return errors.New("invalid modify_type")
		}
		return nil
	}); err != nil {
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

func (r DBRepoImpl) ModifyExploreCategory(modify_type string, category Data.ExploreCategory) error {
	if err := DBMS.Transaction(func(tx *gorm.DB) error {
		explore_category := ExploreCategory{
			CategoryName: category.CategoryName,
			Description:  category.CategoryDescription,
			Country:      category.CategoryCountry,
			image_url:    category.CategoryImage,
		}
		switch modify_type {
		case "Add":
			// カテゴリを追加する
			if err := tx.Create(&explore_category).Error; err != nil {
				tx.Rollback()
				return err
			}
		case "Update":
			// カテゴリを更新する
			if err := tx.Model(&explore_category).Where("category_name = ?", &explore_category.CategoryName).Updates(&explore_category).Error; err != nil {
				tx.Rollback()
				return err
			}
		case "Delete":
			// カテゴリを削除する
			if err := tx.Where("category_name = ?", &explore_category.CategoryName).Delete(&explore_category).Error; err != nil {
				tx.Rollback()
				return err
			}
		case "UnscopedDelete":
			// カテゴリを物理的に削除する
			if err := tx.Where("category_name = ?", &explore_category.CategoryName).Unscoped().Delete(&explore_category).Error; err != nil {
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
