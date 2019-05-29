package controllers

import (
	"blogServer/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"strconv"
)

type TagController struct {
	beego.Controller
}

//@Success 200 {object} models.QueryResult
//@router /list [get]
func (c *TagController) GetAllTags() {
	result := models.QueryResult{}
	tags, err := models.QueryAllTags()

	if err != nil {
		result.OperationResult =
			models.NewOperationResult(models.InternalErr)
	} else {
		result.OperationResult =
			models.NewOperationResult(models.SUCCEED)
		result.Value = tags
	}

	c.Data["json"] = result
	c.ServeJSON()
}

//@Param id body models.Tag true
//@Success 200 {object} models.OperationResult
//@Failure 401 non auth
//@router /add [post]
func (c *TagController) AddTag() {
	var result models.OperationResult
	tag := new(models.Tag)
	err := json.Unmarshal(c.Ctx.Input.RequestBody, tag)
	if err != nil || tag.Name == "" {
		result = models.NewOperationResult(models.InvalidArg)
	} else {
		err := models.CreateTag(tag)
		if err != nil {
			result = models.NewOperationResult(models.InternalErr)
		} else {
			result = models.NewOperationResult(models.SUCCEED)
		}
	}
	c.Data["json"] = result
	c.ServeJSON()
}

//@Param id path int true
//@Success 200 {object} models.OperationResult
//@Failure 401 non auth
//@router /delete/:id [delete]
func (c *TagController) RemoveTag() {
	var result models.OperationResult
	idstr := c.GetString(":id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		result = models.NewOperationResult(models.InvalidArg)
	} else {
		err = models.DeleteTag(id)
		if err != nil {
			result = models.NewOperationResult(models.InternalErr)
		} else {
			result = models.NewOperationResult(models.SUCCEED)
		}
	}
	c.Data["json"] = result
	c.ServeJSON()
}
