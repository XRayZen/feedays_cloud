package APIFunction

import (
	// "heavy/Data"
	"testing"

)

func TestNewSite(t *testing.T) {

}


// RSSのURLを取得する
func TestGetRSSUrl(t *testing.T) {
	// https://jp.ign.com/
    // https://automaton-media.com/
    // https://www.4gamer.net/
	// https://gigazine.net/
	//TODO: ignなどRSSURLをサイトURLを入れないところもあるので対処する
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
			want: "https://automaton-media.com/feed/",
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
			want: "https://gigazine.net/news/rss_2.0/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRSSUrl(tt.args.siteUrl)
			if err != nil {
				t.Error(err)
			}
			if got[0] != tt.want {
				t.Errorf("getRSSUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
