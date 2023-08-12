// パッケージ名が大文字だとパブリック
// パッケージ名が小文字だとプライベート
package RequestHandler

import (
	"errors"
	"read/Repo"
)

func ParseRequestType(di_db_repo Repo.DBRepository, request_type string, user_id string) (res string, err error) {
	// DB接続
	di_db_repo.ConnectDB(false)
	di_db_repo.AutoMigrate()
	// DBからデータを取得するだけの処理をする
	switch request_type {
	case "ExploreCategories":
		res, err = GetExploreCategories(di_db_repo, user_id)
		return res, err
	default:
		return "", errors.New("invalid request type")
	}
}
