package APIFunction

import (
	"testing"
)

func Test_fetchRssArticles(t *testing.T) {
	// https://gigazine.net/news/rss_2.0
	type args struct {
		rssUrl string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "正常系",
			args: args{
				rssUrl: "https://gigazine.net/news/rss_2.0",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			articles, err := fetchRSSArticles(tt.args.rssUrl)
			if err != nil {
				t.Error(err)
			}
			if len(articles) < tt.want {
				t.Errorf("parseRssFeed() = %v, want %v", len(articles), tt.want)
			}
		})
	}
}

func Test_getArticleImageUrls(t *testing.T) {
	_, arg_articles, err := newSite("https://www.4gamer.net/")
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
