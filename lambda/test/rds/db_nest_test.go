package RDS

import (
	// "errors"
	// "log"
	"log"
	"testing"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func TestDbNest(t *testing.T) {
	// SQLLiteに接続する
	InMemoryStr := "file::memory:"
	// fileSqlStr := "test_1.db"
	DB, err := gorm.Open(sqlite.Open(InMemoryStr))
	if err != nil {
		panic("failed to connect database")
	}
	// テストを実行する
	t.Run("DbNestedStructTest", func(t *testing.T) {
		// テストを実行する
		result, err := DbNestedStructTest(DB)
		// エラーが発生した場合はテスト失敗
		if err != nil {
			t.Fatal(err)
		}
		// 結果がfalseの場合はテスト失敗
		if result == false {
			t.Fatal("Test Failed")
		}
	})
}

// User は複数の CreditCards を持ちます。UserID は外部キーとなります。
type User struct {
	gorm.Model
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func TestHasMany(t *testing.T) {
	// SQLLiteに接続する
	InMemoryStr := "file::memory:"
	// fileSqlStr := "test_1.db"
	DB, err := gorm.Open(sqlite.Open(InMemoryStr))
	if err != nil {
		panic("failed to connect database")
	}
	// テストを実行する
	t.Run("HasMany", func(t *testing.T) {
		err := DB.Debug().AutoMigrate(&User{}, &CreditCard{})
		if err != nil {
			log.Println("AutoMigrate Error:", err)
			t.Fatal(err)
		}
		// テストを実行する
		// ユーザーを作成する
		user := User{
			CreditCards: []CreditCard{
				{Number: "411111111111"},
				{Number: "411111111112"},
				{Number: "411111111113"},
			},
		}
		DB.Create(&user)
		// ユーザーを取得する
		var user2 User
		DB.Preload("CreditCards").First(&user2)
		// ユーザーのクレジットカードを取得する
		var creditCards []CreditCard
		DB.Model(&user2).Association("CreditCards").Find(&creditCards)
		// クレジットカードの数を確認する
		if len(creditCards) != 2 {
			t.Fatal("creditCards length is not 2")
		}
		// 取得出来ているか確認する
		if creditCards[0].Number != "411111111111" {
			t.Fatal("creditCards[0].Number is not 411111111111")
		}
		if creditCards[1].Number != "411111111112" {
			t.Fatal("creditCards[1].Number is not 411111111112")
		}
	})
}
