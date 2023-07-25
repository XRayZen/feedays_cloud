package ApiFunction

import (
	"log"
	"os"
	"read/Repo"
	"strconv"
)

func Batch() (bool, error) {
	// バッチ処理を実行する
	// サイトテーブルから読み込んで記事を更新する
	// それにより購読サイトのフィードの鮮度を維持する
	RefreshInterval, err := strconv.Atoi(os.Getenv("REFRESH_INTERVAL"))
	if err != nil {
		log.Println("BATCH RefreshArticles ERROR! :", err)
		return false, err
	}
	res, err := RefreshArticles(Repo.DBRepoImpl{}, RefreshInterval)
	if err != nil {
		return false, err
	}
	log.Println(res)
	return true, nil
}
