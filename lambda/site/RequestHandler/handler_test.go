package RequestHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"site/Data"
	"site/Repo"
	"sort"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

// "site/APIFunction"
// "read/Repo"

// Lambda・Site関数の統合テスト
// リクエストごとにテストを分ける
func TestParseHandlerBySearch(t *testing.T) {
	// DBを初期化
	dbRepo := Repo.DBRepoImpl{}
	setup(t, dbRepo)
	// テスト用のリクエストを作成
	searchRequestByURL := Data.ApiSearchRequest{
		Word:       "https://gigazine.net/",
		SearchType: "url",
	}
	// 記事へのキーワード検索
	searchRequestByKeyword := Data.ApiSearchRequest{
		Word:       "GIGAZINE",
		SearchType: "keyword",
	}
	searchRequestBySiteName := Data.ApiSearchRequest{
		Word:       "GIGAZINE",
		SearchType: "siteName",
	}
	// リクエストをJSONに変換
	searchRequestByKeywordJSON, err := json.Marshal(searchRequestByKeyword)
	if err != nil {
		t.Fatal(err)
	}
	searchRequestByUrlJSON, err := json.Marshal(searchRequestByURL)
	if err != nil {
		t.Fatal(err)
	}
	searchRequestBySiteNameJSON, err := json.Marshal(searchRequestBySiteName)
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
				request_json:  string(searchRequestByKeywordJSON),
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
				request_json:  string(searchRequestByUrlJSON),
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
				request_json:  string(searchRequestBySiteNameJSON),
				request_json2: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", dbRepo, tt.args.request_type, tt.args.user_id, tt.args.request_json, tt.args.request_json2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// JSONを構造体に変換
			var searchResult Data.SearchResult
			err = json.Unmarshal([]byte(result), &searchResult)
			if err != nil {
				t.Fatal(err)
			}
			// リザルトが成功してたらパス
			if searchResult.ApiResponse == "accept" {
				return
			}
			// リザルトが失敗してたらエラー
			t.Errorf("ParseRequestType() = %v, want %v", searchResult.ApiResponse, "accept")
		})
	}
}

// サイトを購読する機能の統合テスト
func TestParseHandlerBySubscribe(t *testing.T) {
	// DBを初期化
	dbRepo := Repo.DBRepoImpl{}
	setup(t, dbRepo)
	// 購読サイトのリクエストを作成
	subscribeWebSite := Data.WebSite{
		SiteName: "GIGAZINE",
		SiteURL:  "https://gigazine.net/",
	}
	newSubscribeWebSite := Data.WebSite{
		SiteName:   "理想ちゃんねる",
		SiteURL:    "http://ideal2ch.livedoor.biz/",
		SiteRssURL: "http://ideal2ch.livedoor.biz/index.rdf",
	}
	// リクエストをJSONに変換
	subscribeRequestJSON, err := json.Marshal(subscribeWebSite)
	if err != nil {
		t.Fatal(err)
	}
	newSubscribeRequestJSON, err := json.Marshal(newSubscribeWebSite)
	if err != nil {
		t.Fatal(err)
	}
	isSubscribeByTrue := true
	isSubscribeByTrueJason, err := json.Marshal(isSubscribeByTrue)
	if err != nil {
		t.Fatal(err)
	}
	isSubscribeByFalse := false
	isSubscribeByFalseJason, err := json.Marshal(isSubscribeByFalse)
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
				request_json:  string(subscribeRequestJSON),
				request_json2: string(isSubscribeByTrueJason),
			},
		},
		// 購読解除テスト
		{
			name: "購読解除テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "SubscribeSite",
				request_json:  string(subscribeRequestJSON),
				request_json2: string(isSubscribeByFalseJason),
			},
		},
		// 新規サイト登録テスト
		{
			name: "新規サイト登録テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_type:  "SubscribeSite",
				request_json:  string(newSubscribeRequestJSON),
				request_json2: string(isSubscribeByTrueJason),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", dbRepo, tt.args.request_type, tt.args.user_id, tt.args.request_json, tt.args.request_json2)
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
				if dbRepo.IsExistSite(newSubscribeWebSite.SiteURL) == false {
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
	getArticlesRequestByLatest := Data.FetchArticlesRequest{
		RequestType: "Latest",
		SiteUrl:     "https://gigazine.net/",
		ReadCount:   10,
	}
	// 今日の0時を取得
	now := time.Now()
	nowRfc3339Str := now.Format(time.RFC3339)
	getArticlesRequestByOldest := Data.FetchArticlesRequest{
		RequestType:    "Old",
		SiteUrl:        "https://gigazine.net/",
		ReadCount:      10,
		OldestModified: nowRfc3339Str,
	}
	// 昨日の0時を取得
	yesterday := now.AddDate(0, 0, -1)
	yesterdayRfc3339Str := yesterday.Format(time.RFC3339)
	getArticlesRequestByUpdate := Data.FetchArticlesRequest{
		RequestType:    "Update",
		SiteUrl:        "https://gigazine.net/",
		ReadCount:      10,
		OldestModified: yesterdayRfc3339Str,
	}
	// リクエストをJSONに変換
	getArticlesRequestByLatestJSON, err := json.Marshal(getArticlesRequestByLatest)
	if err != nil {
		t.Fatal(err)
	}
	getArticlesRequestByOldestJSON, err := json.Marshal(getArticlesRequestByOldest)
	if err != nil {
		t.Fatal(err)
	}
	getArticlesRequestByUpdateJSON, err := json.Marshal(getArticlesRequestByUpdate)
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
				request_json:  string(getArticlesRequestByLatestJSON),
				request_json2: "",
			},
		},
		// 最古記事取得テスト
		{
			name: "最古記事取得テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_json:  string(getArticlesRequestByOldestJSON),
				request_json2: "",
			},
		},
		// 更新記事取得テスト
		{
			name: "更新記事取得テスト",
			args: args{
				access_ip:     "",
				user_id:       "test",
				request_json:  string(getArticlesRequestByUpdateJSON),
				request_json2: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRequestType("", Repo.DBRepoImpl{}, "FetchArticles", tt.args.user_id, tt.args.request_json, tt.args.request_json2)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != "true" {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
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
	insertArticles := []Data.Article{}
	// 最後から10件までを取得
	for i := len(articles) - 1; i > len(articles)-11; i-- {
		insertArticles = append(insertArticles, articles[i])
	}
	dbInsertArticles := []Repo.Article{}
	for _, article := range insertArticles {
		publicationDate, _ := time.Parse(time.RFC3339, article.PublishedAt)
		dbArticle := Repo.Article{
			Title:       article.Title,
			Url:         article.Link,
			Description: article.Description,
			PublishedAt: publicationDate,
		}
		dbInsertArticles = append(dbInsertArticles, dbArticle)
	}

	dbSite := Repo.Site{
		SiteName:     "GIGAZINE",
		RssUrl:       "https://gigazine.net/news/rss_2.0/",
		SiteUrl:      "https://gigazine.net/",
		Description:  "ギガジンのRSSフィード",
		SiteArticles: dbInsertArticles,
	}
	if err := Repo.DBMS.Create(&dbSite).Error; err != nil {
		t.Fatal(err)
	}
	// ユーザーを作成
	dbUser := Repo.User{
		UserName:     "test",
		UserUniqueID: "test",
	}
	if err := Repo.DBMS.Create(&dbUser).Error; err != nil {
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

// 並列処理で記事のイメージURLを取得する
func getArticleImageURLs(articles []Data.Article) ([]Data.Article, error) {
	// 並列処理で記事のイメージURLを取得する
	// 1. og:imageを取得する
	// 3. それでもなければfavicon.icoを取得する
	ch := make(chan Data.Article)
	for _, article := range articles {
		go func(article Data.Article) {
			// 1. og:imageを取得する
			doc, err := getHtmlGoQueryDoc(article.Link)
			if err != nil {
				ch <- article
				return
			}
			imageUrl, err := getArticleImageURL(doc, article.Link)
			if err != nil {
				ch <- article
				return
			}
			article.Image = Data.RssFeedImage{
				Link: imageUrl,
			}
			ch <- article
		}(article)
	}
	for i := 0; i < len(articles); i++ {
		articles[i] = <-ch
	}
	// articleを日時でソートする
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].PublishedAt > articles[j].PublishedAt
	})
	return articles, nil
}

// 記事のイメージURLを取得する
func getArticleImageURL(doc *goquery.Document, articleUrl string) (string, error) {
	// 記事のイメージURLを取得する
	// 1. og:imageを取得する
	// 3. それでもなければfavicon.icoを取得する
	imageUrl := ""
	// 1. og:imageを取得する
	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		property, exists := s.Attr("property")
		if exists {
			if property == "og:image" {
				imageUrl = s.AttrOr("content", "")
				return
			}
		}
	})
	// 2. それでもなければfavicon.icoを取得する
	if imageUrl == "" {
		imageUrl = articleUrl + "/favicon.ico"
	}
	return imageUrl, nil
}

func getHtmlGoQueryDoc(url string) (*goquery.Document, error) {
	// /を消す
	if url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("goquery error: %v", err)
	}
	return doc, nil
}
