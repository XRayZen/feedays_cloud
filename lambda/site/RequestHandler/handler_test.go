package RequestHandler

import (
	"encoding/json"

	"site/Data"
	"site/Repo"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
)

// Lambda・Site関数の統合テスト
// リクエストごとにテストを分ける
func TestParseHandlerBySearch(t *testing.T) {
	// DBを初期化
	db_repo := Repo.DBRepoImpl{}
	setup(t, db_repo)
	// テスト用のリクエストを作成
	search_request_by_url := Data.ApiSearchRequest{
		Word:       "https://gigazine.net/",
		SearchType: "URL",
	}
	// 記事へのキーワード検索
	search_request_by_keyword := Data.ApiSearchRequest{
		Word:       "GIGAZINE",
		SearchType: "Keyword",
	}
	search_request_by_site_name := Data.ApiSearchRequest{
		Word:       "GIGAZINE",
		SearchType: "SiteName",
	}
	// リクエストをJSONに変換
	search_request_by_keyword_json, err := json.Marshal(search_request_by_keyword)
	if err != nil {
		t.Fatal(err)
	}
	search_request_by_url_json, err := json.Marshal(search_request_by_url)
	if err != nil {
		t.Fatal(err)
	}
	search_request_by_site_name_json, err := json.Marshal(search_request_by_site_name)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		access_ip     string
		user_id       string
		request_type  string
		request_json  string
		request_json2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// キーワード検索でサイトが見つかる場合
		{
			name: "キーワード検索でサイトが見つかる場合",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "Search",
				request_json:  string(search_request_by_keyword_json),
				request_json2: "",
			},
		},
		// URL検索でサイトが見つかる場合
		{
			name: "URL検索でサイトが見つかる場合",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "Search",
				request_json:  string(search_request_by_url_json),
				request_json2: "",
			},
		},
		// サイト名検索でサイトが見つかる場合
		{
			name: "サイト名検索でサイトが見つかる場合",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "Search",
				request_json:  string(search_request_by_site_name_json),
				request_json2: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", db_repo, tt.args.request_type, tt.args.user_id, tt.args.request_json, tt.args.request_json2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// JSONを構造体に変換
			var search_result Data.SearchResult
			err = json.Unmarshal([]byte(result), &search_result)
			if err != nil {
				t.Fatal(err)
			}
			// リザルトが成功してたらパス
			if search_result.ApiResponse == "accept" {
				return
			}
			// リザルトが失敗してたらエラー
			t.Errorf("ParseRequestType() = %v, want %v", search_result.ApiResponse, "accept")
		})
	}
}

// サイトを購読する機能の統合テスト
func TestParseHandlerBySubscribe(t *testing.T) {
	// DBを初期化
	db_repo := Repo.DBRepoImpl{}
	setup(t, db_repo)
	// 購読サイトのリクエストを作成
	subscribe_web_site := Data.WebSite{
		SiteName: "GIGAZINE",
		SiteURL:  "https://gigazine.net/",
	}
	new_subscribe_web_site := Data.WebSite{
		SiteName:   "理想ちゃんねる",
		SiteURL:    "http://ideal2ch.livedoor.biz/",
		SiteRssURL: "http://ideal2ch.livedoor.biz/index.rdf",
	}
	// リクエストをJSONに変換
	subscribe_request_json, err := json.Marshal(subscribe_web_site)
	if err != nil {
		t.Fatal(err)
	}
	new_subscribe_request_json, err := json.Marshal(new_subscribe_web_site)
	if err != nil {
		t.Fatal(err)
	}
	is_subscribe_by_true_json, err := json.Marshal(true)
	if err != nil {
		t.Fatal(err)
	}
	is_subscribe_by_false_json, err := json.Marshal(false)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		access_ip     string
		user_id       string
		request_type  string
		request_json  string
		request_json2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// 購読テスト
		{
			name: "購読テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "SubscribeSite",
				request_json:  string(subscribe_request_json),
				request_json2: string(is_subscribe_by_true_json),
			},
		},
		// 購読解除テスト
		{
			name: "購読解除テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "SubscribeSite",
				request_json:  string(subscribe_request_json),
				request_json2: string(is_subscribe_by_false_json),
			},
		},
		// 新規サイト登録テスト
		{
			name: "新規サイト登録テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "SubscribeSite",
				request_json:  string(new_subscribe_request_json),
				request_json2: string(is_subscribe_by_true_json),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", db_repo, tt.args.request_type, tt.args.user_id, tt.args.request_json, tt.args.request_json2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// リザルトがSuccess Subscribe SiteかSuccess Unsubscribe SiteかSuccess Register Siteでなければエラー
			if result != "Success Subscribe Site" && result != "Success Unsubscribe Site" && result != "Success Register Site" {
				t.Errorf("ParseRequestType() = %v, want %v", result, "Success Subscribe Site or Success Unsubscribe Site")
				return
			}
			// もし新規サイト登録だったらDBに登録されているか確認
			if result == "Success Register Site" {
				if db_repo.IsExistSite(new_subscribe_web_site.SiteURL) == false {
					t.Errorf("ParseRequestType() = %v, want %v", result, "Success Register Site")
					return
				}
			}
		})
	}
}

// 記事を取得する機能の統合テスト
func TestParseHandlerByFetchArticles(t *testing.T) {
	// DBを初期化
	setup(t, Repo.DBRepoImpl{})
	// 記事取得のリクエストを作成
	get_articles_request_by_latest := Data.FetchArticlesRequest{
		RequestType: "Latest",
		SiteUrl:     "https://gigazine.net/",
		ReadCount:   10,
	}
	// 今日の0時を取得
	now := time.Now()
	now_rfc3339_str := now.Format(time.RFC3339)
	get_articles_request_by_oldest := Data.FetchArticlesRequest{
		RequestType:    "Old",
		SiteUrl:        "https://gigazine.net/",
		ReadCount:      10,
		OldestModified: now_rfc3339_str,
	}
	// 昨日の0時を取得
	yesterday := now.AddDate(0, 0, -1)
	yesterday_rfc3339_str := yesterday.Format(time.RFC3339)
	get_articles_request_by_update := Data.FetchArticlesRequest{
		RequestType:    "Update",
		SiteUrl:        "https://gigazine.net/",
		ReadCount:      10,
		OldestModified: yesterday_rfc3339_str,
	}
	// リクエストをJSONに変換
	get_articles_request_by_latest_json, err := json.Marshal(get_articles_request_by_latest)
	if err != nil {
		t.Fatal(err)
	}
	get_articles_request_by_oldest_json, err := json.Marshal(get_articles_request_by_oldest)
	if err != nil {
		t.Fatal(err)
	}
	get_articles_request_by_update_json, err := json.Marshal(get_articles_request_by_update)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		access_ip     string
		user_id       string
		request_type  string
		request_json  string
		request_json2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// 最新記事取得テスト
		{
			name: "最新記事取得テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "FetchArticle",
				request_json:  string(get_articles_request_by_latest_json),
				request_json2: "",
			},
		},
		// 最古記事取得テスト
		{
			name: "最古記事取得テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "FetchArticle",
				request_json:  string(get_articles_request_by_oldest_json),
				request_json2: "",
			},
		},
		// 更新記事取得テスト
		{
			name: "更新記事取得テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "FetchArticle",
				request_json:  string(get_articles_request_by_update_json),
				request_json2: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", Repo.DBRepoImpl{}, tt.args.request_type, tt.args.user_id, tt.args.request_json, tt.args.request_json2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// まずは最初のテストケースが上手く行くようにする
			// レスポンスを構造体に変換
			var fetch_article_response Data.FetchArticleResponse
			err = json.Unmarshal([]byte(result), &fetch_article_response)
			if err != nil {
				t.Fatal(err)
			}
			if fetch_article_response.ResponseType != "success" {
				t.Errorf("ParseRequestType() = %v, want %v", fetch_article_response.ResponseType, "success")
				return
			}
		})
	}
}

// ModifyExploreCategory ChangeSiteCategory DeleteSiteの統合テスト
func TestParseHandlerByModifyExploreCategory(t *testing.T) {
	repo := Repo.DBRepoImpl{}
	setup(t, repo)
	// リクエストデータを作成
	category := Data.ExploreCategory{
		CategoryName:    "test Category",
		CategoryCountry: "Japan",
	}
	categoryJSON, _ := json.Marshal(category)
	site_url := "https://gigazine.net/"
	site_url_json, _ := json.Marshal(site_url)
	category_name := "test Category"
	category_name_json, _ := json.Marshal(category_name)
	is_unscoped_json, _ := json.Marshal(true)
	type args struct {
		access_ip     string
		user_id       string
		request_type  string
		request_json  string
		request_json2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// ExploreCategory追加テスト
		{
			name: "ExploreCategory追加テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "ModifyExploreCategory",
				request_json:  string(categoryJSON),
				request_json2: "Add",
			},
			want: "Success ModifyExploreCategory",
		},
		// ExploreCategory削除テスト
		{
			name: "ExploreCategory削除テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "ModifyExploreCategory",
				request_json:  string(categoryJSON),
				request_json2: "UnscopedDelete",
			},
			want: "Success ModifyExploreCategory",
		},
		// ChangeSiteCategory
		{
			name: "ChangeSiteCategory",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "ChangeSiteCategory",
				request_json:  string(site_url_json),
				request_json2: string(category_name_json),
			},
			want: "Success ChangeSiteCategory",
		},
		// DeleteSite
		{
			name: "DeleteSite",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "DeleteSite",
				request_json:  string(site_url_json),
				request_json2: string(is_unscoped_json),
			},
			want: "Success DeleteSite",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", repo, tt.args.request_type, tt.args.user_id, tt.args.request_json, tt.args.request_json2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("ParseRequestType() = %v, want %v", result, tt.want)
			}
		})
	}
}

func setup(t *testing.T, dbRepo Repo.DBRepository) {
	dbRepo.ConnectDB(true)
	dbRepo.AutoMigrate()

	articles, err := fetchRSSArticles("https://gigazine.net/news/rss_2.0/")
	if err != nil {
		t.Fatal(err)
	}
	insert_articles := []Data.Article{}
	// 最後から10件までを取得
	for i := len(articles) - 1; i > len(articles)-11; i-- {
		insert_articles = append(insert_articles, articles[i])
	}
	db_insert_articles := []Repo.Article{}
	for _, article := range insert_articles {
		publication_date, _ := time.Parse(time.RFC3339, article.PublishedAt)
		db_article := Repo.Article{
			Title:       article.Title,
			Url:         article.Link,
			Description: article.Description,
			PublishedAt: publication_date,
		}
		db_insert_articles = append(db_insert_articles, db_article)
	}

	db_site := Repo.Site{
		SiteName:     "GIGAZINE",
		RssUrl:       "https://gigazine.net/news/rss_2.0/",
		SiteUrl:      "https://gigazine.net/",
		Description:  "ギガジンのRSSフィード",
		SiteArticles: db_insert_articles,
	}
	if err := Repo.DBMS.Create(&db_site).Error; err != nil {
		t.Fatal(err)
	}
	// ユーザーを作成
	db_user := Repo.User{
		UserName:     "test",
		UserUniqueID: "test",
	}
	if err := Repo.DBMS.Create(&db_user).Error; err != nil {
		t.Fatal(err)
	}
}

// 指定されたサイトのRSS_URLからRSSフィードを取得して記事リストとして返す
func fetchRSSArticles(rssUrl string) ([]Data.Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssUrl)
	if err != nil {
		return nil, err
	}
	articles := []Data.Article{}
	for _, v := range feed.Items {
		// Feedのカテゴリはタグにしておく
		category := ""
		if len(v.Categories) > 0 {
			category = v.Categories[0]
		}
		article := Data.Article{
			Title:       v.Title,
			Link:        v.Link,
			Description: v.Description,
			Category:    category,
			Site:        feed.Title,
			PublishedAt: v.PublishedParsed.Format(time.RFC3339),
		}
		articles = append(articles, article)
	}
	return articles, nil
}
