package RDS

import (
	// "crypto/tls"
	"fmt"
	"log"
	"os"
	FetchSecret "test/fetch_secret"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	dbName, dbUser, dbPass, dbAddress, dbPort string
	DBMS                                      *gorm.DB
	Debug                             bool
)

func RDS_Connect(DbType, DbName, dbUser, dbPass, dbEndPoint, dbPort string) (*gorm.DB, error) {
	// DB接続
	CONNECT := dbUser + ":" + dbPass + "@tcp(" + dbEndPoint + ":" + dbPort + ")/" + DbName + "?charset=utf8&&parseTime=True&loc=Local"
	cfg:=mysql.Config{
		DSN: CONNECT,
	}

	log.Println("RDS CONNECT STR:", cfg.DSN)
	// DBに管理者権限があるユーザーが作成されていないのでエラーが出る
	DBMS, err := gorm.Open(mysql.New(cfg), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DBMS, nil
}

func GormConnect() (*gorm.DB, error) {
	// AWS Secrets ManagerからDB接続情報を取得する
	// usernameとpasswordはSecrets Managerに保存している
	// シークレットネームはDBユーザー名と同じにしている
	User, Pass, err := FetchSecret.DB_Secret()
	if err != nil {
		fmt.Println("SecretManager.DB_Secret() Error:", err)
		return nil, err
	}
	dbUser = User
	dbPass = Pass
	dbAddress = os.Getenv("rds_endpoint")
	dbName = os.Getenv("db_name")
	dbPort = os.Getenv("db_port")
	DBMS, err := RDS_Connect("mysql", dbName, dbUser, dbPass, dbAddress, dbPort)
	if err != nil {
		fmt.Println("DB接続失敗:", err)
		return nil, err
	}
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
	DB.Transaction(func(tx *gorm.DB) error {
		// トランザクション内でのデータベース処理を行う(ここでは `db` ではなく `tx` を利用する)
		if err := tx.Create(&TestTable{Name: "test", Address: "test", Description: "test"}).Error; err != nil {
			// エラーが発生した場合、ロールバックする
			return err
		}
		// エラーがなければコミットする
		if err := tx.Create(&TestTable{Name: "test", Address: "test", Description: "test"}).Error; err != nil {
		return err
		}
		// nilが返却されるとトランザクション内の全処理がコミットされる
		return nil
	})
	// テーブルからデータを取得“
	var test TestTable
	DB.First(&test, "name = ?", "test")
	fmt.Println(test)
	return true, nil
}
