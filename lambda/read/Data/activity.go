package Data

import "time"

type Activity struct {
	ActivityID        string    `json:"activityId"`
	ActivityType      string    `json:"activityType"`
	ActivityTime      time.Time `json:"activityTime"`
	UserID            string    `json:"userId"`
	AccessEnvironment string    `json:"accessEnvironment"`
	// AccessDevice string
	AccessIP string `json:"accessIP"`
}
