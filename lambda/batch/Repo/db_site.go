package Repo

import (
	"batch/Data"
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
	Country      string
}

// Data.WebSiteとArticlesからDB型に変換する
func convertApiSiteToDb(site Data.WebSite, articles []Data.Article) Site {
	// サイトの最終更新日時をtimeに変換
	var lastModified time.Time
	res_time, err := time.Parse(time.RFC3339, site.LastModified)
	if err != nil {
		lastModified = time.Now().UTC()
	} else {
		lastModified = res_time
	}
	var siteArticles []Article
	for _, siteArticle := range articles {
		var publishedAt time.Time
		time, err := time.Parse(time.RFC3339, siteArticle.LastModified)
		if err != nil {
			publishedAt = time.UTC()
		} else {
			publishedAt = time
		}
		siteArticles = append(siteArticles, Article{
			Title:       siteArticle.Title,
			Url:         siteArticle.Link,
			IconUrl:     siteArticle.Image.Link,
			Description: siteArticle.Description,
			PublishedAt: publishedAt,
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
		SiteArticles: siteArticles,
		IconUrl:      site.SiteImage,
		Description:  site.SiteDescription,
		Category:     site.SiteCategory,
		LastModified: lastModified,
		Tags:         tags,
	}
}

// DB型からAPI型に変換する
func convertDbSiteToApi(site Site) (Data.WebSite, []Data.Article) {
	var siteArticles []Data.Article
	for _, siteArticle := range site.SiteArticles {
		siteArticles = append(siteArticles, Data.Article{
			Title:        siteArticle.Title,
			Description:  siteArticle.Description,
			Link:         siteArticle.Url,
			Image:        Data.RssFeedImage{Link: siteArticle.IconUrl},
			Site:         site.SiteName,
			LastModified: siteArticle.PublishedAt.Format(time.RFC3339),
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

// API型からDB型に記事を変換する
func convertApiArticleToDb(articles []Data.Article) []Article {
	var siteArticles []Article
	for _, siteArticle := range articles {
		PublishedAt, err := time.Parse(time.RFC3339, siteArticle.LastModified)
		if err != nil {
			PublishedAt = time.Now().UTC()
		}
		siteArticles = append(siteArticles, Article{
			Title:       siteArticle.Title,
			Url:         siteArticle.Link,
			IconUrl:     siteArticle.Image.Link,
			Description: siteArticle.Description,
			PublishedAt: PublishedAt,
		})
	}
	return siteArticles
}
