package Repo

import (
	// "read/Data"
	"time"
)


type DBRepo interface{
	// サイトURLをキーにDBに該当するサイトがあるか確認する
	IsExistSite(site_url string) (bool, error)
	// サイトURLをキーに記事更新日時を取得する
	GetSiteLastModified(site_url string) (time.Time, error)

}


type DBRepoImpl struct{
}

func (r DBRepoImpl) IsExistSite(site_url string) (bool, error) {
	return false, nil
}

func (r DBRepoImpl) GetSiteLastModified(site_url string) (time.Time, error) {
	return time.Now(), nil
}
