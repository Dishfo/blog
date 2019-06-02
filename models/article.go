package models

import (
	"github.com/astaxie/beego"
	"time"
)

//在修改某一篇文章时使用
type ArticleEditWrapper struct {
	Fields []string
	Value  *Article
}

type Article struct {
	Id      int64
	Title   string `orm:"size(128)"`
	Publish time.Time
	Summary string `orm:"size(256)"`
	Content string `orm:"size(65535);type(text)"`
	Tags    []*Tag `orm:"rel(m2m)"`
}

//提供对外的函数进行相关的读写操作而隐藏相关的数据库或缓存操作.
func QueryArticleList(pageNo, pageSize int) ([]*Article, error) {
	return queryArticleListInSql(pageNo, pageSize)
}

func QueryArticleById(id int64) (*Article, error) {
	a, err := queryArticleInCache(id)
	if a != nil {
		beego.BeeLogger.Info("hit a article cache")
		_ = accessArticle(a.Id, time.Now().Unix())
		return a, nil
	}

	a, err = queryArticleByIdInSql(id)
	if err != nil {
		return nil, err
	}
	if a != nil {
		_ = cacheArticle(a)
		_ = accessArticle(a.Id, time.Now().Unix())
	}

	return a, nil
}

func QueryArticleByTag(tags []*Tag) ([]*Article, error) {
	return queryArticleByTagInSql(tags)
}

//todo 需要设计权重函数
//QueryTopArticle 返回一个包含了不包含文章内容的列表
func QueryTopArticle(size int) ([]*Article, error) {
	ids, err := queryTopArticlesInCache(size)
	if err != nil {
		return nil, err
	}
	articles := make([]*Article, 0)
	for _, id := range ids {
		a, err := QueryArticleById(id)
		if err != nil || a == nil {
			continue
		}
		a.Content = ""
		articles = append(articles, a)
	}
	return articles, nil
}

func UpdateArticle(a *Article, fields []string) error {
	err := updateArticleInSql(a, fields)
	if err != nil {
		return err
	}

	//todo 如果已缓存就更新
	_ = updateArticleInCache(a, fields)
	return nil
}

func DeleteArticle(id int64) error {
	err := deleteArticleInSql(id)
	clearArticle(id)
	return err
}

func CreateArticle(a *Article) error {
	err := insertArticleInSql(a)
	if err != nil {
		return err
	}
	_ = cacheArticle(a)
	return nil
}

//todo 添加额外的线程每隔一段时间更新缓存中的内容
//GetRecommendArticles 根据参数提供的文章id 返回推荐的文章
func GetRecommendArticles(ids []int64) ([]*Article, error) {
	articleIds, err := getRecommendArticleIds(ids)
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, 0)
	for _, articleId := range articleIds {
		a, err := QueryArticleById(articleId)
		if err != nil {
			return nil, err
		}
		if a == nil || a.Id == 0 {
			continue
		}
		a.Content = ""
		articles = append(articles, a)
	}

	return articles, nil
}
