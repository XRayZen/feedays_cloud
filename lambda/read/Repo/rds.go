package Repo

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)



func handler()  {
	// 環境変数からDB接続情報を取得する
	db_name:= os.Getenv("db_name")

	// AWS Secrets ManagerからDB接続情報を取得する
	// usernameとpasswordはSecrets Managerに保存している
	
}

