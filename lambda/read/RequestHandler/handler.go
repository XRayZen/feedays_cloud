// パッケージ名が大文字だとパブリック
// パッケージ名が小文字だとプライベート
package RequestHandler

import (
	"errors"
	"read/Repo"
)

func ParseRequestType(diDBRepo Repo.DBRepository, requestType string, userID string) (res string, err error) {
	// DB接続
	diDBRepo.ConnectDB(false)
	diDBRepo.AutoMigrate()
	// DBからデータを取得するだけの処理をする
	switch requestType {
	case "ExploreCategories":
		res, err = GetExploreCategories(diDBRepo, userID)
		return res, err
	default:
		return "", errors.New("invalid request type")
	}
}
