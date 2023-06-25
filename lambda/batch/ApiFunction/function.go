package ApiFunction


import (
	"log"
	"read/Repo"
)

func Batch() (bool, error) {
	// バッチ処理を実行する
	// サイトテーブルから読み込んで記事を更新する
	// それにより購読サイトのフィードの鮮度を維持する
	res,err:= RefreshArticles(Repo.DBRepoImpl{})
	if err != nil {
		return false, err
	}
	log.Println(res)
	return true, nil
}

