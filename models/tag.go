package models

import "log"

type Tag struct {
	Id       int64
	Name     string     `orm:"unique"`
	Articles []*Article `orm:"reverse(many)"`
}

func QueryAllTags() ([]*Tag, error) {
	tags, err := getTagsInCache()
	if err != nil {
		return nil, err
	}
	if tags != nil || len(tags) != 0 {
		log.Println("hit in cache")
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
	/*tags, err := queryTagsByIdInCache(ids)
	if err != nil {
		return nil, err
	}
	if tags != nil || len(tags) != 0 {
		log.Println("hit in cache")
		return tags, nil
	}*/
	return queryTagsByIdInSql(ids)
}
