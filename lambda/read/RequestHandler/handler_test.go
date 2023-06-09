package RequestHandler

import (
	"encoding/json"
	"read/DBRepo"
	"read/Data"
	"testing"
)

func TestParseRequestType(t *testing.T) {
	// 正解のデータを用意する
	want_Ex, _ := json.Marshal(Data.ExploreCategories{
		CategoryName: "CategoryName",
	})
	want_ExStr := string(want_Ex)
	want_Rs, _ := json.Marshal(Data.Ranking{
		RankingName: "RankingName",
	})
	want_RsStr := string(want_Rs)
	type args struct {
		diDBRepo    DBRepo.DBRepository
		requestType string
		userID      string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常系 ExploreCategories",
			args: args{
				diDBRepo:    DBRepo.MockDBRepo{},
				requestType: "ExploreCategories",
				userID:      "userID",
			},
			want: want_ExStr,
		},
		{
			name: "正常系 Ranking",
			args: args{
				diDBRepo:    DBRepo.MockDBRepo{},
				requestType: "Ranking",
				userID:      "userID",
			},
			want: want_RsStr,
		},
		{
			name: "異常系 invalid request type",
			args: args{
				diDBRepo:    DBRepo.MockDBRepo{},
				requestType: "invalid request type",
				userID:      "userID",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRequestType(tt.args.diDBRepo, tt.args.requestType, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRequestType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRequestType() = %v, want %v", got, tt.want)
			}
		})
	}

}