package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/**
提供article 相关缓存支持
*/
const (
	articleKeyPrefix    = "article:"
	articleAcsKeyPrefix = "article_view:"
)

//用于存储在缓存系统中的article内容 hset用于
type ArticleC struct {
	Id      int64  `redis:"id"`
	Title   string `redis:"title"`
	Publish string `redis:"publish"`
	Summary string `redis:"summary"`
	Content string `redis:"content"`
	Tags    string `redis:"tags"`
	//LastAccess int64     `redis:"access"` //文章的最后一次访问 移除这个字段使用list或zset
}

func (ac *ArticleC) toArticle() *Article {
	a := new(Article)
	a.Title = ac.Title
	a.Publish, _ = time.Parse(time.UnixDate, ac.Publish)
	a.Content = ac.Content
	a.Id = ac.Id
	a.Summary = ac.Summary

	//反序列化tags
	tags := make([]*Tag, 0)
	_ = json.Unmarshal([]byte(ac.Tags), &tags)
	a.Tags = tags

	return a
}

func toArticleC(a *Article) *ArticleC {
	tagsJson, _ := json.Marshal(a.Tags)
	ac := &ArticleC{
		Id:      a.Id,
		Title:   a.Title,
		Publish: a.Publish.Format(time.UnixDate),
		Summary: a.Summary,
		Content: a.Content,
		Tags:    string(tagsJson),
	}
	return ac
}

func cacheArticle(a *Article) error {
	conn := client.Get()
	key := articleKey(a.Id)
	ac := toArticleC(a)
	r, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(ac)...)
	if err != nil {
		return err
	}
	s, _ := redis.String(r, err)
	beego.BeeLogger.Info("%s", s)
	err = conn.Flush()
	return err
}

func updateArticleInCache(a *Article, fields []string) error {
	conn := client.Get()
	key := articleKey(a.Id)
	_, err := conn.Do("MULTI")
	if err != nil {
		return err
	}
	ac := toArticleC(a)
	for _, f := range fields {
		_, _ = conn.Do("HSET", key,
			strings.ToLower(f),
			articleField(ac, f))
	}
	_, _ = conn.Do("EXEC")
	return nil
}

func queryArticleInCache(id int64) (*Article, error) {
	conn := client.Get()
	key := articleKey(id)
	values, err := redis.Values(conn.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	if values == nil || len(values) == 0 {
		return nil, nil
	}
	articleC := new(ArticleC)
	err = redis.ScanStruct(values, articleC)
	return articleC.toArticle(), err
}

func accessArticle(id, access int64) error {
	conn := client.Get()
	key := articleAccessKey(id)
	_, err := conn.Do("SADD", key, access)
	return err
}

//获取访问量 最前面size的文章
/**
此处需要设计一个权重函数
*/
func QueryTopArticlesInCache(size int) ([]int64, error) {

	return nil, nil
}

func clearArticle(id int64) {
	conn := client.Get()
	key := articleKey(id)
	_, _ = conn.Do("DEL", key)
}

func articleAccessKey(id int64) string {
	return fmt.Sprintf("%s%s",
		articleAcsKeyPrefix,
		strconv.FormatInt(id, 10))
}

func articleKey(id int64) string {
	return fmt.Sprintf("%s%s",
		articleKeyPrefix,
		strconv.FormatInt(id, 10))
}

func articleField(ac *ArticleC, field string) interface{} {
	v := reflect.ValueOf(*ac)
	f := v.FieldByName(field)
	if f.Kind() == reflect.String {
		return f.String()
	} else {
		return f.Int()
	}
}
