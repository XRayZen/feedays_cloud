// パッケージ名が大文字だとパブリック
// パッケージ名が小文字だとプライベート
package RequestHandler

import (
	"errors"
	"read/Repo"
	// "read/Data"
)

func ParseRequestType(diDBRepo Repo.DBRepository, requestType string, userID string) (res string, err error) {
	// DBからデータを取得するだけの処理をする
	switch requestType {
	case "FetchRanking":
		res, err := FetchRanking(userID)
		return res, err
	case "ExploreCategories":
		// エントリポイントでDIしたのを入れる
		str := Explore{
			DBrepo: diDBRepo,
		}
		res, err = str.GetExploreCategories(userID)
		return res, err
	default:
		return "", errors.New("invalid request type")
	}
}
