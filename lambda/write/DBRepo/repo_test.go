package DBRepo

import (
	// "read/Data"
	"read/Data"
	"testing"
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

func TestDbRepoTest(t *testing.T) {
	// まずはUserを生成する
	dbRepo := InitDataBase()
	// 登録するUserを生成する
	user := Data.UserConfig{
		UserName:      "MockUser",
		UserUniqueID:  "0000",
		AccountType:   "JP",
		SearchHistory: []Data.SearchHistory{},
		ClientConfig:  Data.ClientConfig{},
	}

	// Userを取得する
	t.Run("GetUser", func(t *testing.T) {
		// Userを登録する
		if err := dbRepo.RegisterUser(user); err != nil {
			t.Errorf("failed to create user: %v", err)
		}
		// Userを取得する
		user, err := dbRepo.GetUserConfig("0000")
		// 取得出来たのがMockUserか確認する
		if err != nil || user.UserName != "MockUser" {
			t.Errorf("failed to get user: %v", err)
		}
		user.ClientConfig = Data.ClientConfig{
			ApiConfig: Data.ApiConfig{
				FetchArticleRequestInterval: 4000,
			},
		}
		// Userを更新する
		if err := dbRepo.UpdateClientConfig(user.UserUniqueID, user); err != nil {
			t.Errorf("failed to update user: %v", err)
		}
		// Userが更新されているか確認する
		user, err = dbRepo.GetUserConfig(user.UserUniqueID)
		// 数字が反映されていない場合はエラー
		if err != nil || user.ClientConfig.ApiConfig.FetchArticleRequestInterval != 4000 {
			// 検証に失敗した場合はエラー
			t.Errorf("failed to validation user: %v", err)
		}
		// Userを削除する
		if err := dbRepo.DeleteUser(user.UserUniqueID); err != nil {
			t.Errorf("failed to delete user: %v", err)
		}
	})
}
