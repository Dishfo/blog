package test

import (
	"blogServer/models"
	_ "blogServer/routers"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
	models.RegisterDB()
	models.InitRedis()
}

// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}

func TestAddArticle(t *testing.T) {

}

func TestLogin(t *testing.T) {
	client := &http.Client{}
	r, _ := http.NewRequest("POST", "/v1/admin/login", nil)
	r.PostForm = url.Values{
		"user": {"dishfo"},
		"pass": {"159357ghj"},
	}

	client.Do(r)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())
}

func TestTagQuery(t *testing.T) {
	client := &http.Client{}
	r, _ := http.NewRequest("POST", "/v1/admin/login", nil)
	r.PostForm = url.Values{
		"user": {"dishfo"},
		"pass": {"159357ghj"},
	}

	client.Do(r)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

}
