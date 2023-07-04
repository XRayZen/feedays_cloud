package SecretManager

import (
	// "os"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func createClient(region string) (*secretsmanager.SecretsManager, error) {
    session, err := session.NewSession(&aws.Config{
        Region: aws.String(region),
    })
    if err != nil {
        return nil, err
    }
    return secretsmanager.New(session), nil
}

func getSecretValue(client *secretsmanager.SecretsManager, secretName string,secretStage string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
		VersionStage: aws.String(secretStage),
	}
	result, err := client.GetSecretValue(input)
	if err != nil {
		return "", err
	}
	return *result.SecretString, nil
}

func DB_Secret()(string,string,error){
	// AWS Secrets ManagerからDB接続情報を取得する
	// usernameとpasswordはSecrets Managerに保存している
	region := os.Getenv("region")
	client,err := createClient(region)
	if err != nil {
		return "","",err
	}
	secret_stage := os.Getenv("secret_stage")
	// シークレットネームはDBユーザー名と同じにしている
	secret_name := os.Getenv("db_username")
	db_username,err := getSecretValue(client,secret_name,secret_stage)
	if err != nil {
		return "","",err
	}
	db_password,err := getSecretValue(client,secret_name,secret_stage)
	if err != nil {
		return "","",err
	}
	return db_username,db_password,nil
}

func Secret_read_test() (bool,error) {
	db_username,db_password,err := DB_Secret()
	if err != nil {
		return false,err
	}
	// ちゃんと取得できているか確認してできたらtrueを返す
	// usernameがadminでpasswordが文字数10以上であることを確認する
	if db_username == "admin" && len(db_password) >= 10 {
		return true,nil
	} else {
		return false,nil
	}
}
