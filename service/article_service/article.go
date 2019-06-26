package article_service

import (
	"encoding/json"
	"gin-blog/models"
	"gin-blog/pkg/gredis"
	"gin-blog/pkg/logging"
	"gin-blog/service/cache_service"
)

type Article struct {
	ID            int    `gorm:"primary_key" json:"id"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	ModifiedBy    string `json:"modifled_by" valid:"Required;MaxSize(100)"`

	PageNum  int
	PageSize int
}

func (a *Article) Add(form models.AddArticleForm) error {

	article := models.Article{
		TagID:         form.TagID,
		Desc:          form.Desc,
		CoverImageUrl: form.CoverImageUrl,
		CreatedBy:     form.CreatedBy,
		State:         form.State,
		Title:         form.Title,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit(form models.EditArticleForm) error {

	a.TagID = form.TagID
	a.Desc = form.Desc
	a.CoverImageUrl = form.CoverImageUrl
	a.ModifiedBy = form.ModifiedBy
	a.State = form.State
	a.Title = form.Title

	return models.EditArticle(a.ID, a.getMaps())
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)

	cache := cache_service.Article{
		TagID:     a.TagID,
		State:     a.State,
		CreatedBy: a.CreatedBy,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticles)
		}
	}

	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *Article) Delete() error {
	return models.DeleteArticle(a.ID)
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleByID(a.ID)
}

func (a *Article) Count() (int, error) {
	return models.GetArticleTotal(a.getMaps())
}

func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	if a.CreatedBy != "" {
		maps["created_by"] = a.CreatedBy
	}
	if a.Title != "" {
		maps["title"] = a.Title
	}
	if a.Desc != "" {
		maps["desc"] = a.Desc
	}
	if a.Content != "" {
		maps["content"] = a.Content
	}
	if a.ModifiedBy != "" {
		maps["modified_by"] = a.ModifiedBy
	}
	if a.CoverImageUrl != "" {
		maps["cover_image_url"] = a.CoverImageUrl
	}

	return maps
}
