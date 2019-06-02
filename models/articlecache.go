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

	articleViewTime  = "article_access"
	articleViewCount = "article_view"
	articleViewSort  = "article_seq"
)

var (
	maxTopArticle = 80
	maxCacheNum   = 50
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
	_, _ = redis.String(r, err)
	//beego.BeeLogger.Info("%s", s)
	//err = conn.Flush()

	return err
}

//到底要不要receive呢
func updateArticleInCache(a *Article, fields []string) error {
	conn := client.Get()
	key := articleKey(a.Id)
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	ac := toArticleC(a)
	for _, f := range fields {
		_ = conn.Send("HSET", key,
			strings.ToLower(f),
			articleField(ac, f))
	}
	_ = conn.Send("EXEC")
	_ = conn.Flush()
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

//暂时不记录访问的时间
func accessArticle(id, access int64) error {
	conn := client.Get()
	key := articleAccessKey(id)
	err := conn.Send("MULTI")
	_ = conn.Send("ZADD", articleViewTime, access, id)
	err = conn.Send("ZINCRBY", articleViewSort, -2, key)
	err = conn.Send("ZINCRBY", articleViewCount, 1, key)
	err = conn.Send("EXEC")
	err = conn.Flush()
	return err
}

func clearArticles() {
	conn := client.Get()
	ids, err := redis.Int64s(conn.Do("ZRANGE", articleViewTime, 0, -(maxCacheNum + 1)))
	if err != nil {
		return
	}
	args := []interface{}{
		articleViewTime,
	}

	_ = conn.Send("MUTLI")
	for _, id := range ids {
		key := articleKey(id)
		args = append(args, id)
		_ = conn.Send("DEL", key)
	}

	_ = conn.Send("ZREM", args...)
	_ = conn.Send("EXEC")
	_ = conn.Flush()
}

//
func clearAccessSeq() {
	conn := client.Get()
	err := conn.Send("MULTI")
	err = conn.Send("ZREMRANGEBYRANK", articleViewSort, maxTopArticle+1, -1)
	err = conn.Send("ZREMRANGEBYRANK", articleViewCount, 0, -100-1)
	err = conn.Send("ZINTERSTORE", articleViewSort, 2, articleViewSort,
		articleViewCount)
	err = conn.Send("EXEC")
	err = conn.Flush()
	if err != nil {
		beego.BeeLogger.Error("%s", err.Error())
	}
}

//获取访问量 最前面size的文章
/**
此处需要设计一个权重函数
*/
func queryTopArticlesInCache(size int) ([]int64, error) {
	conn := client.Get()
	mems, err := redis.Strings(conn.Do("ZRANGE", articleViewSort, 0, size-1))
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0)

	for _, mem := range mems {
		s := strings.TrimLeft(mem, articleAcsKeyPrefix)
		id, _ := strconv.ParseInt(s, 10, 64)
		ids = append(ids, id)
	}

	return ids, nil
}

func clearArticle(id int64) {
	conn := client.Get()
	key := articleKey(id)
	key2 := articleAccessKey(id)
	_ = conn.Send("MUTLI")
	_ = conn.Send("DEL", key)
	_ = conn.Send("ZREM", articleViewSort, key2)
	_ = conn.Send("ZREM", articleViewCount, key2)
	_ = conn.Send("ZREM", articleViewTime, id)
	_ = conn.Send("EXEC")
	_ = conn.Flush()
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
