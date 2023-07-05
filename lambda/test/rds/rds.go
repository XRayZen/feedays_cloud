package RDS

import (
	"fmt"
	"os"
	"test/secret_manager"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// テスト用テーブル
type TestTable struct {
	ID          int `gorm:"primary_key"`
	Name        string
	Address     string
	Description string
}

// DB接続情報のキャッシュ
var (
	dbName, dbUser, dbPass, dbAddress string
	DBMS                              *gorm.DB
	// debug                             bool
)

func RDS_Connect(DbType, DbName, dbUser, dbPass, dbEndPoint string) (*gorm.DB, error) {
	// DB接続
	DBMS, err := gorm.Open(DbType, dbUser+":"+dbPass+"@tcp("+dbEndPoint+")/"+DbName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	return DBMS, nil
}

func GormConnect() (*gorm.DB, error) {
	// 環境変数からDB接続情報を取得する
	// AWS Secrets ManagerからDB接続情報を取得する
	// usernameとpasswordはSecrets Managerに保存している
	// シークレットネームはDBユーザー名と同じにしている
	User, Pass, err := SecretManager.DB_Secret()
	if err != nil {
		fmt.Println("SecretManager.DB_Secret() Error:", err)
		return nil, err
	}
	dbUser = User
	dbPass = Pass
	dbAddress = os.Getenv("rds_endpoint")
	dbName = os.Getenv("db_name")
	DBMS, err := RDS_Connect("mysql", dbName, dbUser, dbPass, dbAddress)
	if err != nil {
		fmt.Println("DB接続失敗:", err)
		return nil, err
	}
	DBMS.LogMode(true)
	return DBMS, nil
}

func GormMigrateTable(DBMS *gorm.DB) *gorm.DB {
	DBMS.AutoMigrate(&TestTable{})
	return DBMS
}

func RdsWriteReadTest() (bool, error) {
	DB, err := GormConnect()
	if err != nil {
		return false, err
	}
	// テーブル作成
	DB = GormMigrateTable(DB)
	// トランザクション開始
	tx := DB.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}
	// テーブルにデータを追加
	tx.Create(&TestTable{Name: "test", Address: "test", Description: "test"})
	// トランザクション終了
	tx.Commit()
	// テーブルからデータを取得“
	var test TestTable
	DB.First(&test, "name = ?", "test")
	fmt.Println(test)
	return true, nil
}
