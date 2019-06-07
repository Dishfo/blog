package models

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	InitModels()
	runtime.GOMAXPROCS(2)
}

func TestTagQuery(t *testing.T) {

	//tag := &Tag{
	//	Name: "golang",
	//}
	//
	//o := orm.NewOrm()
	//o.Begin()
	//_, err := o.Insert(tag)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//for i := 0; i < 3; i++ {
	//	a := &Article{
	//		Title:   "1111",
	//		Publish: time.Now(),
	//	}
	//	a.Tags = []*Tag{
	//		tag,
	//	}
	//
	//	o.Insert(a)
	//
	//	m2m := o.QueryM2M(a, "Tags")
	//	m2m.Add(tag)
	//}
	///*a2 := new(Article)
	//o.QueryTable("article").
	//	Filter("Tags__Tag__Id",tag.Id).One(a2)
	//o.LoadRelated(a2,"Tags")*/
	//o.Commit()
	//n := time.Now()

	tags, err := QueryAllTags()
	if err != nil {
		t.Fatal(err.Error())
	}

	log.Println(tags)
	for i := 0; i < 2000; i++ {
		queryRelatedArticle(1)
	}

	//e := time.Now()
	//t.Log(e.Sub(n).Seconds() * 1000)
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

	/*tags, _ = QueryAllTags()
	for _, tag := range tags {
		//DeleteTag(tag.Id)
	}*/
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

var (
	contentFile = "/home/dishfo/文档/content.txt"
)

func TestArticleCacheMod(t *testing.T) {
	//data, _ := ioutil.ReadFile(contentFile)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			a := &Article{
				Title:   "test",
				Summary: "wada",
				Publish: time.Now(),
				Content: string("123"),
				Tags: []*Tag{
					&Tag{
						Id: 23,
					},
					&Tag{
						Id: 20,
					},
				},
			}
			for i := 0; i < 5; i++ {
				a.Id = 0
				err := CreateArticle(a)
				if err != nil {
					t.Log(err)
				}
				//QueryArticleById(a.Id)

			}
			wg.Done()
		}()
	}

	wg.Wait()

	//1clearArticles()

}

func TestQueryAricle(t *testing.T) {
	articles, err := QueryArticleList(0, 100)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(len(articles))
}

func TestUpdateArticle(t *testing.T) {
	articles, _ := QueryArticleList(0, 1)
	articles[0].Summary = "11111"
	articles[0].Tags = []*Tag{
		&Tag{
			Id: 23,
		},
		&Tag{
			Id: 20,
		},
	}
	var err error
	var wg sync.WaitGroup
	for i := 0; i <= 30; i++ {
		wg.Add(1)
		go func() {
			err = UpdateArticle(articles[0], []string{
				"Content", "Tags",
			})
			wg.Done()
		}()
	}

	wg.Wait()
	if err != nil {
		t.Log(err)
	}
}

func TestRpcCall(t *testing.T) {
	/*content, err := json.Marshal([]int{1, 2, 3, 4, 5})
	if err != nil {
		return
	}
	q, err := maCh.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	err = maCh.Publish("rpc",
		"recommended",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          content,
			CorrelationId: "one",
			ReplyTo:       q.Name,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
	msgs, err := maCh.Consume(q.Name, "recs", true,
		false,
		false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for d := range msgs {
		log.Println(string(d.Body))
	}*/
	ids := []int{
		1, 2, 3, 4,
	}

	toJson := make(map[string]interface{})
	toJson["ids"] = ids
	toJson["size"] = 0
	b, _ := json.Marshal(toJson)
	t.Log(string(b))

}

func TestEmpty(t *testing.T) {

}

func TestRedisPipeLine(t *testing.T) {
	conn := client.Get()
	conn.Send("SET", "A", 123456)
	conn.Send("SET", "A", 123456)
	conn.Flush()

	str, err := redis.String(conn.Do("GET", "A"))
	t.Log(str, err)
}
