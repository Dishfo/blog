package models

import (
	"github.com/astaxie/beego"
)

type Tag struct {
	Id       int64      `gorm:"primary_key type:int auto_increment;"`
	Name     string     `gorm:"unique_index;type:varchar(20);"`
	Articles []*Article `gorm:"mamy2many:article_tags PRELOAD:false;"`
}

func (Tag) TableName() string {
	return "tag"
}

func QueryAllTags() ([]*Tag, error) {
	tags, err := getTagsInCache()
	if err != nil {
		return nil, err
	}
	if tags != nil || len(tags) != 0 {
		beego.BeeLogger.Info("hit in cache")
		return tags, nil
	}
	return queryAllTagsInSql()
}

func CreateTag(t *Tag) error {
	id, err := insertTagInSql(t)
	if err != nil {
		return err
	}
	t.Id = id
	_ = cacheTagsInCache(t)
	return nil
}

func DeleteTag(id int64) error {
	clearTagInCache(id)
	return deleteTagInSql(id)
}

func QueryTagsById(ids []int64) ([]*Tag, error) {
	return queryTagsByIdInSql(ids)
}
