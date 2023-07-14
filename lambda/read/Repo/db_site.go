package Repo

import (
	"read/Data"
	"time"

	"gorm.io/gorm"
)

type Site struct {
	// gorm.Modelをつけると、idとCreatedAtとUpdatedAtとDeletedAtが作られる
	// gormでは、gorm.Modelでdeleted_atを作成している場合、Deleteすると、自動的に論理削除になるという仕様
	gorm.Model
	SiteName          string
	SiteUrl           string
	RssUrl            string
	IconUrl           string
	Description       string
	SiteFeeds         []SiteArticle
	Tags              []Tag
	Category          string
	LastModified      time.Time
	SubscriptionCount int
}

type SiteArticle struct {
	gorm.Model
	SiteID       uint
	ArticleIndex int
	Title        string
	Url          string
	IconUrl      string
	Description  string
	ReadLater    bool
	PublishedAt  time.Time
}

type Tag struct {
	gorm.Model
	TagName string
	SiteID  uint
}

type ExploreCategory struct {
	gorm.Model
	CategoryName string
	Description  string
	Country      string
}

// API型からDB型に変換する
func convertApiSiteToDb(site Data.WebSite, articles []Data.Article) Site {
	var siteArticles []SiteArticle
	for _, siteArticle := range articles {
		var publishedAt time.Time
		time, err := time.Parse(time.RFC3339, siteArticle.LastModified)
		if err != nil {
			publishedAt = time.UTC()
		} else {
			publishedAt = time
		}
		siteArticles = append(siteArticles, SiteArticle{
			Title:       siteArticle.Title,
			Url:         siteArticle.Link,
			IconUrl:     siteArticle.Image.Link,
			Description: siteArticle.Description,
			PublishedAt: publishedAt,
		})
	}
	return Site{
		SiteName:    site.SiteName,
		SiteUrl:     site.SiteURL,
		RssUrl:      site.SiteRssURL,
		SiteFeeds:   siteArticles,
		IconUrl:     site.SiteImage,
		Description: site.SiteDescription,
		Category:    site.SiteCategory,
	}
}

// DB型からAPI型に変換する
func convertDbSiteToApi(site Site) (Data.WebSite, []Data.Article) {
	var siteArticles []Data.Article
	for _, siteArticle := range site.SiteFeeds {
		siteArticles = append(siteArticles, Data.Article{
			Index:        0,
			Title:        siteArticle.Title,
			Description:  siteArticle.Description,
			Link:         siteArticle.Url,
			Image:        Data.RssFeedImage{Link: siteArticle.IconUrl},
			Site:         site.SiteName,
			LastModified: siteArticle.PublishedAt.Format(time.RFC3339),
			IsReedLate:   siteArticle.ReadLater,
			Category:     site.Category,
			SiteUrl:      site.SiteUrl,
		})
	}
	return Data.WebSite{
		SiteName:        site.SiteName,
		SiteImage:       site.IconUrl,
		SiteDescription: site.Description,
		SiteURL:         site.SiteUrl,
		SiteRssURL:      site.RssUrl,
		LastModified:    site.LastModified.Format(time.RFC3339),
		SiteCategory:    site.Category,
	}, siteArticles
}
