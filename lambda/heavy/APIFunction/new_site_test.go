package APIFunction

import (
	// "heavy/Data"
	// "heavy/Data"
	"testing"
)

// RSSのURLを取得する
func TestNewSite(t *testing.T) {
	// https://jp.ign.com/
	// https://automaton-media.com/
	// https://www.4gamer.net/
	// https://gigazine.net/
	// https://randomwalker.blog.fc2.com/
	// https://takezo50.com/
	// https://kabumatome.doorblog.jp/
	// https://techlife.cookpad.com/
	// https://codezine.jp/
	// https://techblog.yahoo.co.jp/
	// https://developer.hatenastaff.com/
	// TestCaseを書く
	type args struct {
		siteUrl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "正常系",
			args: args{
				siteUrl: "https://jp.ign.com/",
			},
			want: "https://jp.ign.com/feed.xml",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://automaton-media.com/",
			},
			want: "https://automaton-media.com/feed",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://www.4gamer.net/",
			},
			want: "https://www.4gamer.net/rss/index.xml",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://gigazine.net/",
			},
			want: "https://gigazine.net/news/rss_2.0",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://randomwalker.blog.fc2.com/",
			},
			want: "https://randomwalker.blog.fc2.com/?xml",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://takezo50.com/",
			},
			want: "https://takezo50.com/feed",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://kabumatome.doorblog.jp/",
			},
			want: "https://kabumatome.doorblog.jp/index.rdf",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://techlife.cookpad.com/",
			},
			want: "https://techlife.cookpad.com/rss",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://codezine.jp/",
			},
			want: "https://codezine.jp/rss/new/index.xml",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://techblog.yahoo.co.jp/",
			},
			want: "https://techblog.yahoo.co.jp/index.xml",
		},
		{
			name: "正常系",
			args: args{
				siteUrl: "https://developer.hatenastaff.com/",
			},
			want: "https://developer.hatenastaff.com/rss",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Url := tt.args.siteUrl
			// Urlに/がある場合は、/を削除する
			if Url[len(Url)-1:] == "/" {
				Url = Url[:len(Url)-1]
			}
			webSite, articles, err := NewSite(Url)
			if err != nil {
				t.Error(err)
			}
			if webSite.SiteURL != Url {
				t.Errorf("NewSite() = %v, want %v", webSite.SiteURL, Url)
			}
			if len(articles) == 0 {
				t.Errorf("NewSite() = %v, want %v", len(articles), "0")
			}
		})
	}
}

// RSSFeedをパースする
func TestParseRSSFeed(t *testing.T) {
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
			articles, err := parseRssFeed(tt.args.rssUrl)
			if err != nil {
				t.Error(err)
			}
			if len(articles) < tt.want {
				t.Errorf("parseRssFeed() = %v, want %v", len(articles), tt.want)
			}
		})
	}

}
