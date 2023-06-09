// パッケージ名が大文字だとパブリック
// パッケージ名が小文字だとプライベート
package RequestHandler

import (
	"errors"
	"read/DBRepo"
	// "read/Data"
)

func ParseRequestType(diDBRepo DBRepo.DBRepository, requestType string, userID string) (res string, err error) {
	// DBからデータを取得するだけの処理をする
	switch requestType {
	case "ExploreCategories":
		// エントリポイントでDIしたのを入れる
		str := Explore{
			DBrepo: diDBRepo,
		}
		res,err= str.GetExploreCategories(userID)
		return res, err
	case "Ranking":
		str := Ranking{
			DBrepo: diDBRepo,
		}
		res,err= str.GetRanking(userID)
		return res, err
	default:
		return "", errors.New("invalid request type")
	}
}
