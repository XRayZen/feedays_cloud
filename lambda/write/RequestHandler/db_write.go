package RequestHandler

import (
	"encoding/json"
	"write/DBRepo"
)

func GetUserInfo(repo DBRepo.DBRepo, userId string) (string, error) {
	res, err := repo.GetUserInfo(userId)
	if err != nil {
		return "", err
	}
	str, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func RegisterUser(repo DBRepo.DBRepo, userId string, value string, argument_value2 string) (string, error) {
	return "", nil
}

func ReportActivity(repo DBRepo.DBRepo, userId string, value string) (string, error) {
	return "", nil
}

func SyncConfig(repo DBRepo.DBRepo, userId string, value string) (string, error) {
	return "", nil
}

func CodeSync(repo DBRepo.DBRepo, userId string, value string) (string, error) {
	return "", nil
}
