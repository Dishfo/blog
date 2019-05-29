package models

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	RegisterDB()
	InitRedis()

}

func TestTagQuery(t *testing.T) {

	tag := &Tag{
		Name: "golang",
	}

	o := orm.NewOrm()
	o.Begin()
	_, err := o.Insert(tag)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		a := &Article{
			Title:   "1111",
			Publish: time.Now(),
		}
		a.Tags = []*Tag{
			tag,
		}

		o.Insert(a)

		m2m := o.QueryM2M(a, "Tags")
		m2m.Add(tag)
	}
	/*a2 := new(Article)
	o.QueryTable("article").
		Filter("Tags__Tag__Id",tag.Id).One(a2)
	o.LoadRelated(a2,"Tags")*/
	o.Commit()
	n := time.Now()
	for i := 0; i < 2000; i++ {
		queryRelatedArticle(1)
	}

	e := time.Now()
	t.Log(e.Sub(n).Seconds() * 1000)
	/*if err!=nil {
		t.Fatal(err)
	}
	b,_ := json.MarshalIndent(articles," ","  ")
	t.Log(string(b))*/

}

func TestTagInsert(t *testing.T) {
	tag := &Tag{
		Name: "golang",
	}

	tag1 := &Tag{
		Name: "java",
	}

	tag2 := &Tag{
		Name: "c++",
	}

	insertTagInSql(tag)
	insertTagInSql(tag1)
	insertTagInSql(tag2)

	deleteTagInSql(2)

	tags, err := queryAllTagsInSql()
	log.Println(err)
	b, _ := json.MarshalIndent(tags, "", " ")
	log.Println(string(b))

	tag.Id = 0
	tag.Name = "c"
	CreateTag(tag)

	tags, _ = QueryAllTags()
	for _, tag := range tags {
		DeleteTag(tag.Id)
	}
}

func TestTagCache(t *testing.T) {
	tags, err := getTagsInCache()
	if err != nil {
		t.Fatal(err)
	}

	err = cacheTagsInCache(&Tag{
		Id:   54,
		Name: "awdawad",
	})
	t.Log(err)
	t.Log(tags)
}

func TestArticleCache(t *testing.T) {
	a := &Article{
		Title:   "test",
		Summary: "wada",
		Publish: time.Now(),
		Content: "dwada",
		Tags: []*Tag{
			&Tag{
				Name: "json",
				Id:   1,
			},
			&Tag{
				Id:   2,
				Name: "c++",
			},
		},
	}
	err := cacheArticle(a)
	if err != nil {
		t.Fatal(err)
	}

	//err = CreateArticle(a)
	//if err!=nil {
	//	log.Println(err)
	//}

	a.Title = "dishfo test"

	updateArticleInCache(a, []string{
		"Title",
	})
	a.Tags = []*Tag{
		&Tag{
			Id:   7,
			Name: "golang",
		},
	}
	updateArticleInCache(a, []string{
		"Tags",
	})

	a, err = queryArticleInCache(0)

	b, _ := json.MarshalIndent(a, "", "	")
	log.Println(string(b))
	//QueryArticleById(a.Id)
}