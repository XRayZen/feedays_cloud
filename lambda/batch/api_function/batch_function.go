package ApiFunction

import (
	"batch/Repo"
	"log"
	"os"
	"strconv"
)

func InitDataBase(isTestMode bool) (Repo.DBRepository, int, error) {
	RefreshInterval := 0
	DbRepo := Repo.DBRepoImpl{}
	if isTestMode {
		RefreshInterval = 15
		// DB初期化処理
		DbRepo.ConnectDB(true)
		DbRepo.AutoMigrate()
	} else {
		GetEnvRefreshInterval, err := strconv.Atoi(os.Getenv("REFRESH_INTERVAL"))
		if err != nil {
			log.Println("BATCH RefreshArticles ERROR! :", err)
			return DbRepo, RefreshInterval, err
		}
		RefreshInterval = GetEnvRefreshInterval
		// DB初期化処理
		DbRepo.ConnectDB(false)
		DbRepo.AutoMigrate()
	}
	return DbRepo, RefreshInterval, nil
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
