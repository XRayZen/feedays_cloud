package RDS

import (
	"os"
	"test/secret_manager"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	dbUser, dbPass, dbAddress, dbSalt string
	DBMS                              *gorm.DB
	debug                             bool
)

func handler()  {
	// 環境変数からDB接続情報を取得する
	// AWS Secrets ManagerからDB接続情報を取得する
	// usernameとpasswordはSecrets Managerに保存している
	// シークレットネームはDBユーザー名と同じにしている
	User, Pass, err := SecretManager.DB_Secret()
	if err != nil {
		panic(err)
	}
	dbUser = User
	dbPass = Pass
	dbAddress = os.Getenv("rds_endpoint")
	


}



