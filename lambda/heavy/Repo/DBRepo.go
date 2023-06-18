package Repo

import (
	// "read/Data"
	// "read/Data"
	// "time"
)

// type DBRepo interface {
// 	// サイトURLをキーにDBに該当するサイトがあるか確認する
// 	IsExistSite(site_url string) (bool, error)
// 	// サイトURLをキーにDBに該当するサイトを返す
// 	GetSite(site_url string) (Data.WebSite, error)
// 	// サイトURLをキーに記事更新日時を取得する
// 	GetSiteLastModified(site_url string) (time.Time, error)
// 	// 新規サイトをDBに登録する
// 	RegisterSite(site Data.WebSite,articles []Data.Article) error
// 	// キーワード検索でDBに該当する記事を返す
// 	SearchArticlesByKeyword(keyword string) ([]Data.Article, error)
// }

// type DBRepoImpl struct {
// }

// func (r DBRepoImpl) IsExistSite(site_url string) (bool, error) {
// 	return false, nil
// }

// func (r DBRepoImpl) GetSite(site_url string) (Data.WebSite, error) {
// 	return Data.WebSite{}, nil
// }

// func (r DBRepoImpl) GetSiteLastModified(site_url string) (time.Time, error) {
// 	return time.Now(), nil
// }

// func (r DBRepoImpl) RegisterSite(site Data.WebSite,articles []Data.Article) error {
// 	return nil
// }

// func (r DBRepoImpl) SearchArticlesByKeyword(keyword string) ([]Data.Article, error) {
// 	return nil, nil
// }