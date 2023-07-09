package Repo

import "gorm.io/gorm"

type User struct {
	gorm.Model
	user_name         string
	user_id           int `gorm:"uniqueIndex"`
	account_type      string
	api_activity      ApiActivity        `gorm:"foreignKey:user_id"`
	favorite_site     []FavoriteSite     `gorm:"foreignKey:user_id"`
	subscription_site []SubscriptionSite `gorm:"foreignKey:user_id"`
	read_history      []ReadHistory      `gorm:"foreignKey:user_id"`
	search_history    []SearchHistory    `gorm:"foreignKey:user_id"`
	client_config     ClientConfig       `gorm:"foreignKey:user_id"`
}

type ClientConfig struct {
	gorm.Model
	user_id               uint `gorm:"uniqueIndex"`
	feed_refresh_interval int
	ui_config             UiConfig  `gorm:"foreignKey:client_config_id"`
	api_config            ApiConfig `gorm:"foreignKey:client_config_id"`
}

type ApiActivity struct {
	gorm.Model
	user_id         uint
	activity_type   string
	activity_time   string
	access_platform string
	access_ip       string
}

type FavoriteSite struct {
	gorm.Model
	user_id       uint
	favorite_type string
	site_id       uint
}

type SubscriptionSite struct {
	gorm.Model
	user_id  uint
	site_url string
}

type SearchHistory struct {
	gorm.Model
	user_id     uint
	search_word string
}

type ReadHistory struct {
	gorm.Model
	user_id  uint
	site_url string
	feed_url string
}

type ApiConfig struct {
	gorm.Model
	client_config_id             uint
	account_type                 string
	fetch_feed_request_interval  int
	fetch_feed_request_limit     int
	fetch_trend_request_interval int
	fetch_trend_request_limit    int
}

type UiConfig struct {
	gorm.Model
	client_config_id    uint
	mobile_text_size    int
	tablet_text_size    int
	theme_color         string
	theme_mode          string
	drawer_menu_opacity float32
}
