package DBRepo

import (
	// "read/Data"
	"read/Data"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// テスト用のモックデータを生成する
func InitDataBase() DBRepo {
	dbRepo := DBRepoImpl{}
	// MockModeでRDSではなくインメモリーsqliteに接続する
	if err := dbRepo.ConnectDB(true); err != nil {
		panic("failed to connect database")
	}
	if err := dbRepo.AutoMigrate(); err != nil {
		panic("failed to migrate database")
	}
	// カテゴリを生成する
	var categories = []ExploreCategory{
		{
			CategoryName: "CategoryName",
			Country:      "JP",
		},
	}
	// カテゴリを保存する
	DBMS.Create(&categories)
	return dbRepo
}

// GIGAZINEのサイトを取得して生成する

func TestDbRepoTest(t *testing.T) {
	// まずはUserを生成する
	dbRepo := InitDataBase()
	// 登録するUserを生成する
	user := Data.UserConfig{
		UserName:      "MockUser",
		UserUniqueID:  "0000",
		AccountType:   "Free",
		Country:       "JP",
		SearchHistory: []Data.SearchHistory{},
		ClientConfig:  Data.ClientConfig{},
		ReadHistory:   []Data.ReadHistory{},
	}

	// Userを取得する
	t.Run("GetUser", func(t *testing.T) {
		user.ClientConfig.ApiConfig.FetchArticleRequestInterval = 1000
		user.ClientConfig.ApiConfig.FetchTrendRequestInterval = 2000
		// Userを登録する
		if err := dbRepo.RegisterUser(user); err != nil {
			t.Errorf("failed to create user: %v", err)
		}
		// Userを取得する
		user, err := dbRepo.SearchUserConfig(user.UserUniqueID)
		// 取得出来たのがMockUserか確認する
		if err != nil || user.UserName != "MockUser" || user.ClientConfig.ApiConfig.FetchArticleRequestInterval != 1000 || user.ClientConfig.ApiConfig.FetchTrendRequestInterval != 2000 {
			t.Errorf("failed to get user: %v", err)
		}
		user.ClientConfig.ApiConfig.RefreshArticleInterval = 4000
		// Userを更新する
		if err := dbRepo.UpdateUser(user.UserUniqueID, user); err != nil {
			t.Errorf("failed to update user: %v", err)
		}
		// Userが更新されているか確認する
		user, err = dbRepo.SearchUserConfig(user.UserUniqueID)
		// 数字が反映されていない場合はエラー
		if err != nil || user.ClientConfig.ApiConfig.RefreshArticleInterval != 4000 {
			// 検証に失敗した場合はエラー
			t.Errorf("failed to update validation: %v", err)
		}
		// 閲覧履歴を追加する
		readActivity := Data.ReadHistory{
			Link:           "link",
			AccessAt:       time.Now().Format(time.RFC3339),
			AccessPlatform: "PC",
			AccessIP:       "10.9.9.9",
		}
		if err := dbRepo.AddReadHistory(user.UserUniqueID, readActivity); err != nil {
			t.Errorf("failed to add read history: %v", err)
		}
		//閲覧履歴を取得して検証する
		user, err  = dbRepo.SearchUserConfig(user.UserUniqueID)
		if err != nil || len(user.ReadHistory) == 0 {
			t.Errorf("failed to get read history: %v", err)
		}
	})
}

// DBの関連付けのデータをCRUDテスト

type UserAs struct {
	ID       int
	Name     string
	Articles []Article
}

type Article struct {
	ID       int
	UserAsID uint
	Name     string
}

func TestDbRepoTest2(t *testing.T) {
	// もしモックモードが有効ならSqliteに接続する
	InMemoryStr := "file::memory:"
	DB, err := gorm.Open(sqlite.Open(InMemoryStr))
	if err != nil {
		panic("failed to connect database")
	}
	isDbConnected = true
	DBMS = DB
	// テーブルを作成する
	DBMS.AutoMigrate(&UserAs{}, &Article{})
	// ユーザーを作成する
	user := UserAs{
		Name: "User",
		// Articles: []Article{},
	}
	DBMS.Create(&user)
	t.Run("GORMAssociationTest", func(t *testing.T) {
		if err := DBMS.Model(&user).Association("Articles").Append(&Article{Name: "Article"}); err != nil {
			t.Errorf("failed to append article: %v", err)
		}
		var articles []Article
		if err := DBMS.Model(&user).Association("Articles").Find(&articles); err != nil {
			t.Errorf("failed to find article: %v", err)
		}
		if articles[0].Name != "Article" {
			t.Errorf("failed to find article: %v", err)
		}
	})
}
