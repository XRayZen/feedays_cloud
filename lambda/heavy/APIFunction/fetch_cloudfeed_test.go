package APIFunction

import (
	"testing"
)

func TestFetchCloudFeed(t *testing.T) {

}

func Test_getArticleImageUrls(t *testing.T) {
	_, arg_articles, err := NewSite("https://www.4gamer.net/")
	if err != nil {
		t.Errorf("NewSite error: %v", err)
	}
	res_articles, err := getArticleImageURLs(arg_articles)
	if err != nil {
		t.Errorf("getArticleImageURLs error: %v", err)
	}
	if len(res_articles) == 0 {
		t.Errorf("getArticleImageURLs error: %v", err)
	}
	t.Run("getArticleImageURLs", func(t *testing.T) {
		// イメージURLが取得できているか確認する
		for _, article := range res_articles {
			if article.Image.Link == "" {
				// 記事URLとエラーを出力する
				t.Errorf("getArticleImageURLs error: %v ,siteURL: %v", err, article.Link)
			}
		}
	})
}
