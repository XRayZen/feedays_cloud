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
	SiteName     string
	SiteUrl      string
	RssUrl       string
	IconUrl      string
	Description  string
	SiteArticles []Article
	Tags         []Tag
	Category     string
	// 記事を更新したら、LastModifiedを更新する
	LastModified      time.Time
	SubscriptionCount int
}

type Article struct {
	gorm.Model
	SiteID      uint
	Title       string
	Url         string
	IconUrl     string
	Description string
	PublishedAt time.Time
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
	image_url    string
	Country      string
}

// Data.WebSiteとArticlesからDB型に変換する
func convertApiSiteToDb(site Data.WebSite, articles []Data.Article) Site {
	// サイトの最終更新日時をtimeに変換
	var last_modified time.Time
	res_time, err := time.Parse(time.RFC3339, site.LastModified)
	if err != nil {
		last_modified = time.Now().UTC()
	} else {
		last_modified = res_time
	}
	var site_articles []Article
	for _, site_article := range articles {
		var published_at time.Time
		time, err := time.Parse(time.RFC3339, site_article.PublishedAt)
		if err != nil {
			published_at = time.UTC()
		} else {
			published_at = time
		}
		site_articles = append(site_articles, Article{
			Title:       site_article.Title,
			Url:         site_article.Link,
			IconUrl:     site_article.Image.Link,
			Description: site_article.Description,
			PublishedAt: published_at,
		})
	}
	// タグを変換
	var tags []Tag
	for _, tag := range site.SiteTags {
		tags = append(tags, Tag{TagName: tag})
	}

	return Site{
		SiteName:     site.SiteName,
		SiteUrl:      site.SiteURL,
		RssUrl:       site.SiteRssURL,
		SiteArticles: site_articles,
		IconUrl:      site.SiteImage,
		Description:  site.SiteDescription,
		Category:     site.SiteCategory,
		LastModified: last_modified,
		Tags:         tags,
	}
}

// DB型からAPI型に変換する
func convertDbSiteToApi(site Site) (Data.WebSite, []Data.Article) {
	var site_articles []Data.Article
	for _, site_article := range site.SiteArticles {
		site_articles = append(site_articles, Data.Article{
			Title:       site_article.Title,
			Description: site_article.Description,
			Link:        site_article.Url,
			Image:       Data.RssFeedImage{Link: site_article.IconUrl},
			Site:        site.SiteName,
			PublishedAt: site_article.PublishedAt.Format(time.RFC3339),
			Category:    site.Category,
			SiteUrl:     site.SiteUrl,
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
	}, site_articles
}

// API型からDB型に記事を変換する
func convertApiArticleToDb(articles []Data.Article) []Article {
	var site_articles []Article
	for _, site_article := range articles {
		published_at, err := time.Parse(time.RFC3339, site_article.PublishedAt)
		if err != nil {
			published_at = time.Now().UTC()
		}
		site_articles = append(site_articles, Article{
			Title:       site_article.Title,
			Url:         site_article.Link,
			IconUrl:     site_article.Image.Link,
			Description: site_article.Description,
			PublishedAt: published_at,
		})
	}
	return site_articles
}
