package ApiFunction

import (
	"batch/Repo"
	"log"
	"os"
	"strconv"
)

func InitDataBase(is_test_mode bool) (Repo.DBRepository, int, error) {
	refresh_interval := 0
	db_repo := Repo.DBRepoImpl{}
	if is_test_mode {
		refresh_interval = 15
		// DB初期化処理
		db_repo.ConnectDB(true)
		db_repo.AutoMigrate()
	} else {
		get_env_refresh_interval, err := strconv.Atoi(os.Getenv("REFRESH_INTERVAL"))
		if err != nil {
			log.Println("BATCH RefreshArticles ERROR! :", err)
			return db_repo, refresh_interval, err
		}
		refresh_interval = get_env_refresh_interval
		// DB初期化処理
		db_repo.ConnectDB(false)
	}
	return db_repo, refresh_interval, nil
}

func Batch(dbRepo Repo.DBRepository, refreshInterval int) (bool, error) {
	// バッチ処理を実行する
	// サイトテーブルから読み込んで記事を更新する
	// それにより購読サイトのフィードの鮮度を維持する
	res, err := RefreshArticles(dbRepo, refreshInterval)
	if err != nil {
		return false, err
	}
	log.Println(res)
	return true, nil
}
