package RequestHandler

import (
	"time"
	"write/DBRepo"
	"write/Data"
)

func ReportAPIActivity(ip string, dbRepo DBRepo.DBRepo, userId string, identInfo Data.UserAccessIdentInfo,
	activityType string) error {
	res := Data.Activity{
		UserID:       userId,
		ActivityType: activityType,
		ActivityTime: time.Now().UTC(),
		// Webかモバイルかデスクトップか
		AccessEnvironment: identInfo.PlatformType,
		AccessIP:          ip,
		ActivityID:        "",
	}
	return dbRepo.AddApiActivity(userId, res)
}
