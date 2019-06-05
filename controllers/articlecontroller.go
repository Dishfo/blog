package controllers

import (
	"blogServer/models"
	"blogmesssage"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	visitedCookieName = "visited"
	recommendUrl      = "/v1/article/"
	cookieMaxDay      = 180
)

type ArticleController struct {
	beego.Controller
}

var (
	validFieldName = "(Title|Tags|Summary|Content)"
	regex          *regexp.Regexp
)

func init() {
	var err error
	regex, err = regexp.Compile(validFieldName)
	if err != nil {
		log.Fatal(err.Error())
	}
}

//read
//@Param pageNo  query int true   哪一页
//@Param pageSize  query int true  一页的大小
//@Success 200 {object} models.QueryResult
//@router /list [get]
func (c *ArticleController) GetArticleList() {
	var result models.QueryResult
	pageNo, err := c.GetInt("pageNo")
	pageSize, err := c.GetInt("pageSize")

	if err != nil {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
	} else {
		artilces, err := models.QueryArticleList(pageNo, pageSize)
		if err != nil {
			result.OperationResult =
				models.NewOperationResult(models.InternalErr)
		} else {
			result.OperationResult =
				models.NewOperationResult(models.SUCCEED)
			result.Value = artilces
		}
	}
	c.Data["json"] = result
	c.ServeJSON()
}

//@Param id path int64 true "Specific id of article"
//@Success 200 {object} models.QueryResult
//@router /getArticle/:id [get]
func (c *ArticleController) GetArticle() {
	var result models.QueryResult
	id, err := c.GetInt64(":id", 0)
	if err != nil || id <= 0 {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
	} else {
		a, err := models.QueryArticleById(id)
		if err != nil {
			result.OperationResult =
				models.NewOperationResult(models.InternalErr)
		} else {
			result.OperationResult =
				models.NewOperationResult(models.SUCCEED)
			result.Value = a
			appendVisitRecord(c, id)
		}
	}

	c.Data["json"] = result
	c.ServeJSON()
}

//@Param tags query []int64 true "Specific tags"
//@Success 200 {object} models.QueryResult
//@router /getArticleByTag [get]
func (c *ArticleController) GetArticleByTag() {
	var result models.QueryResult
	var err error
	var num int64
	tagIds := make([]int64, 0)
	tags := c.GetString("tags")
	ids := strings.Split(tags, ",")
	for _, id := range ids {
		num, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			break
		}
		tagIds = append(tagIds, num)
	}

	if err != nil {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
	} else {
		tags, err := models.QueryTagsById(tagIds)
		if err != nil {
			result.OperationResult =
				models.NewOperationResult(models.InternalErr)
			goto out
		}
		artilces, err := models.QueryArticleByTag(tags)
		if err != nil {
			result.OperationResult =
				models.NewOperationResult(models.InternalErr)
		} else {
			result.OperationResult =
				models.NewOperationResult(models.SUCCEED)
			result.Value = artilces
		}
	}
out:
	c.Data["json"] = result
	c.ServeJSON()
}

//@Param size query int true "Number of articles required"
//@Success 200 {object} models.QueryResult
//@router /getTopArticles [get]
func (c *ArticleController) GetTopArticles() {
	var result models.QueryResult
	size, err := c.GetInt("size")
	if err != nil || size <= 0 {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
	} else {
		articles, err := models.QueryTopArticle(size)
		if err != nil {
			result.OperationResult =
				models.NewOperationResult(models.InternalErr)
		} else {
			result.OperationResult =
				models.NewOperationResult(models.SUCCEED)
			result.Value = articles
		}
	}

	c.Data["json"] = result
	c.ServeJSON()
}

//write
//@Param article body models.Article true "create article"
//@Success 200 {object} models.OperationResult
//@router /add [post]
func (c *ArticleController) AddArticle() {
	var result models.OperationResult
	a := new(models.Article)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, a)
	ok, _ := validArticleForCreate(a)
	if err != nil || !ok {
		result = models.NewOperationResult(models.InvalidArg)
	} else {
		a.Publish = time.Now()
		err := models.CreateArticle(a)
		if err != nil {
			result = models.NewOperationResult(models.InternalErr)
		} else {
			result = models.NewOperationResult(models.SUCCEED)
			onAddArticle(a)
		}
	}

	c.Data["json"] = result
	c.ServeJSON()
}

//@Param article body models.ArticleEditWrapper true
//@Success 200 {object} models.OperationResult
//@router /edit [put]
func (c *ArticleController) EditArticle() {
	var result models.OperationResult
	aw := new(models.ArticleEditWrapper)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, aw)
	ok, _ := validArticleForUpdate(aw.Value, aw.Fields)
	if err != nil || !ok {
		result = models.NewOperationResult(models.InvalidArg)
	} else {
		a := aw.Value
		err := models.UpdateArticle(a, aw.Fields)
		if err != nil {
			result = models.NewOperationResult(models.InternalErr)
		} else {
			result = models.NewOperationResult(models.SUCCEED)
			onUpdateArticle(a, aw.Fields)
		}
	}
	c.Data["json"] = result
	c.ServeJSON()
}

//@Param id path int true "id of article"
//@Success 200 {object} models.OperationResult
//@router /delete/:id [delete]
func (c *ArticleController) RemoveArticle() {
	var result models.OperationResult
	id, err := c.GetInt64(":id", 0)
	if err != nil || id <= 0 {
		result = models.NewOperationResult(models.InvalidArg)
	} else {
		err = models.DeleteArticle(id)
		if err != nil {
			result = models.NewOperationResult(models.InternalErr)
		} else {
			result = models.NewOperationResult(models.SUCCEED)
		}
	}
	c.Data["json"] = result
	c.ServeJSON()
}

//todo 添加推荐文章的数目限制
//@Param size query int true "number of recommend article"
//@Success 200 {object} models.QueryResult
//@router /recommended   [get]
func (c *ArticleController) RecommendArticles() {
	var result models.QueryResult
	var visited = parseVisited(c.Ctx.Request)
	var articles []*models.Article
	var size int
	visitIds := make([]int64, 0)

	err := json.Unmarshal([]byte(visited.Value), &visitIds)
	if err != nil || !validIdArray(visitIds) {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
		goto out
	}

	size, err = c.GetInt("size")
	if err != nil || size <= 0 {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
		goto out
	}

	articles, err = models.GetRecommendArticles(visitIds, size)
	if err != nil {
		result.OperationResult =
			models.NewOperationResult(models.InvalidArg)
	} else {
		result.OperationResult =
			models.NewOperationResult(models.SUCCEED)
		result.Value = articles
	}
out:
	c.Data["json"] = result
	c.ServeJSON()
}

//todo 实现结构体参数的验证
func validArticleForCreate(a *models.Article) (bool, string) {
	valid := validation.Validation{}
	valid.Required(a.Title, "title")
	valid.Required(a.Summary, "summary")
	valid.Required(a.Content, "content")
	if len(a.Tags) > 0 {
		for _, tag := range a.Tags {
			valid.Required(tag.Id, "tag_id")
		}
	}
	validStringLength(&valid, a)
	if len(valid.Errors) > 0 {
		sb := new(strings.Builder)
		for _, e := range valid.Errors {
			sb.WriteString(e.Name + " ")
			sb.WriteString(e.String())
			sb.WriteString("\n")
		}
		return false, sb.String()
	}

	return true, ""
}

func validIdArray(ids []int64) bool {
	for _, id := range ids {
		if id <= 0 {
			return false
		}
	}

	return true
}

func validArticleForUpdate(a *models.Article, fields []string) (bool, string) {
	valid := validation.Validation{}
	valid.Required(a.Id, "id")
	validStringLength(&valid, a)

	//判断field 是否合法
	for _, field := range fields {
		valid.Match(field, regex, field)
	}

	if len(valid.Errors) > 0 {
		sb := new(strings.Builder)
		for _, e := range valid.Errors {
			sb.WriteString(e.Name + " ")
			sb.WriteString(e.String())
			sb.WriteString("\n")
		}
		return false, sb.String()
	}

	return true, ""
}

func validStringLength(v *validation.Validation, a *models.Article) {
	v.MaxSize(a.Content, 65535, "content_len")
	v.MaxSize(a.Title, 128, "title_len")
	v.MaxSize(a.Summary, 256, "summary_len")
}

func onAddArticle(a *models.Article) {
	msg := new(blogmesssage.ArticleMessage)
	msg.ArticleId = a.Id
	msg.ArticleContent = a.Content
	msg.Op = blogmesssage.CreateArticle
	models.SendArticleMessage(msg)
	//投递消息到消息队列中
}

func onUpdateArticle(a *models.Article, fields []string) {
	isEditContent := false
	for _, f := range fields {
		if f == "Content" {
			isEditContent = true
			break
		}
	}

	if !isEditContent {
		return
	}
	msg := new(blogmesssage.ArticleMessage)
	msg.ArticleId = a.Id
	msg.ArticleContent = a.Content
	msg.Op = blogmesssage.UpdateArticle
	models.SendArticleMessage(msg)
}

//在访问某一篇文章后使用
func appendVisitRecord(c *ArticleController, id int64) {
	visited := parseVisited(c.Ctx.Request)
	var visitedIds []int64
	err := json.Unmarshal([]byte(visited.Value), &visitedIds)
	if err != nil {
		visitedIds = make([]int64, 0)
	}

	visitedIds = append(visitedIds, id)
	b, _ := json.Marshal(visitedIds)
	visited.Value = string(b)
	http.SetCookie(c.Ctx.ResponseWriter, visited)
}

func parseVisited(r *http.Request) *http.Cookie {
	var visited *http.Cookie
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == visitedCookieName {
			visited = cookie
		}
	}

	if visited == nil {
		visited = new(http.Cookie)
		visited.Value = "[]"
		visited.Name = visitedCookieName
		visited.HttpOnly = true
		visited.Path = "/v1/article/"
		visited.Expires = time.Now().Add(time.Hour * 24 * cookieMaxDay)
	}

	return visited
}
