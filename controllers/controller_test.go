package controllers

import (
	"blogServer/models"
	"testing"
)

func TestArticleController_AddArticle(t *testing.T) {
	ok, str := validArticleForUpdate(&models.Article{
		Id: 1,
	}, []string{
		"Title",
	})

	t.Log(ok, str)
}
