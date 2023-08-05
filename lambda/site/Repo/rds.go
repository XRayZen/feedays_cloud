package Repo

import (
	"os"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/secretsmanager"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB接続情報のキャッシュ
var (
	dbType, dbName, dbUser, dbPass, dbEndPoint, dbPort string
	DBMS                                               *gorm.DB
	isDbConnected                                      bool
	Debug                                              bool
)

// DBに操作する時に呼び出してDB接続処理がしていないのならDB接続処理をする関数
// していたらそのままDB接続情報を返す
func DataBaseConnect() error {
	if isDbConnected {
		return nil
	}
	// DB接続情報を取得する
	dbName = os.Getenv("db_name")
	dbUser = os.Getenv("db_username")
	secret, err := FetchDbSecret()
	if err != nil {
		return err
	}
	dbPass = secret.Password
	dbType = "mysql"
	dbPort = os.Getenv("db_port")
	dbEndPoint = os.Getenv("rds_endpoint")
	db, err := RDS_Connect(dbType, dbName, dbUser, dbPass, dbEndPoint, dbPort)
	if err != nil {
		return err
	}
	DBMS = db
	isDbConnected = true
	return nil
}

// DBが切断されたらDB接続情報を削除する
func Disconnect() {
	DBMS = nil
	isDbConnected = false
}

func RDS_Connect(DbType, DbName, dbUser, dbPass, dbEndPoint, dbPort string) (*gorm.DB, error) {
	// DB接続
	CONNECT := dbUser + ":" + dbPass + "@tcp(" + dbEndPoint + ":" + dbPort + ")/" + DbName
	OPTION_STR := "?charset=utf8&&parseTime=True&loc=Local"
	cfg := mysql.Config{
		DSN: CONNECT + OPTION_STR,
	}

	// DBに管理者権限があるユーザーが作成されていないのでエラーが出る
	DBMS, err := gorm.Open(mysql.New(cfg), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DBMS, nil
}

// DBに使う全ての型をマイグレートする
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&User{})
}
