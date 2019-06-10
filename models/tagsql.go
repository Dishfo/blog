package models

/**
提供从mysql数据库读写tag
*/
const (
	articleTable = "article"
	tagTable     = "tag"
)

func queryAllTagsInSql() ([]*Tag, error) {
	tags := make([]*Tag, 0)
	err := dbInstance.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func deleteTagInSql(id int64) error {
	err := dbInstance.Delete(Tag{}, id).Error
	return err
}

func insertTagInSql(tag *Tag) (int64, error) {
	db := dbInstance.Create(tag)
	return db.RowsAffected, db.Error
}

func queryTagsByIdInSql(ids []int64) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	dbInstance.Where("id in (?) ", ids).Find(&tags)
	return tags, dbInstance.Error
}

func queryRelatedArticle(id int64) ([]*Article, error) {
	articles := make([]*Article, 0)
	db := dbInstance.
		Where("id = ?", id).
		First(&Tag{}).
		Related(&articles, "Articles").Find(&articles)
	if db.Error != nil {
		return nil, db.Error
	}
	return articles, nil
}
