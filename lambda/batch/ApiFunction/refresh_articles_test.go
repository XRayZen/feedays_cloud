package ApiFunction

import (
	"read/Repo"
	"testing"
)

func TestRefreshArticles(t *testing.T) {
	// 考えられるテストケースを網羅する
	// 更新する記事がない場合
	// 更新する記事がある場合
	type args struct {
		repo         Repo.DBRepository
		MockSiteLastModified int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "更新する記事がない場合",
			args: args{
				repo:         Repo.MockDBRepo{},
				MockSiteLastModified: 10,
			},
			want:    "No Update",
			wantErr: false,
		},
		{
			name: "更新する記事がある場合",
			args: args{
				repo:         Repo.MockDBRepo{},
				MockSiteLastModified: 30,
			},
			want:    "BATCH RefreshArticles SUCCESS!",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Repo.MockSiteLastModified = tt.args.MockSiteLastModified
			got, err := RefreshArticles(tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RefreshArticles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
