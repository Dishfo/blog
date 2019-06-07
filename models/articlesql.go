package models

import (
	"database/sql"
	"strings"
)

//todo 将beego orm 替换为gorm

func queryArticleListInSql(pageno, size int) ([]*Article, error) {
	offset := pageno * size
	var articles []*Article
	dbInstance.
		Limit(size).
		Offset(offset).
		Preload("Tags").
		Find(&articles)

	return articles, nil
}

func queryArticleByIdInSql(id int64) (*Article, error) {
	a := new(Article)
	err := dbInstance.First(a, id).Error
	if err != nil {
		return nil, err
	}
	if a.Id == 0 {
		return nil, nil
	}
	return a, nil
}

//返回的article 中只需要包含指定的部分tag
func queryArticleByTagInSql(tags []*Tag) ([]*Article, error) {
	tagIds := make([]int64, 0)
	for _, tag := range tags {
		tagIds = append(tagIds, tag.Id)
	}
	articles := make([]*Article, 0)

	return articles, nil
}

func insertArticleInSql(a *Article) error {
	tags := a.Tags
	a.Tags = nil
	tx := dbInstance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(a).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(a).Association("Tags").Append(tags).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func updateArticleInSql(a *Article, fields []string) error {
	/*o := orm.NewOrm()
	var err error
	err = o.Begin()

	//修改相关的tag关系
	if setTag {
		m2m := o.QueryM2M(a, "Tags")
		_, err = m2m.Clear()
		if err == nil {
			_, err = m2m.Add(a.Tags)
		}
	}

	if err != nil {
		_ = o.Rollback()
	} else {
		_ = o.Commit()
	}
	*/
	var err error
	tx := dbInstance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	setTag := false
	filterFields := make([]string, 0)
	for _, f := range fields {
		if f != "Tags" {
			filterFields = append(filterFields, f)
		} else {
			setTag = true
		}
	}

	if len(filterFields) != 0 {
		err = tx.
			Model(a).
			Select(strings.Join(filterFields, ",")).
			Updates(a).Error
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	if setTag {
		err = tx.Model(a).Association("Tags").Replace(a.Tags).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error

}

func deleteArticleInSql(id int64) error {
	return dbInstance.Delete(&Article{}, id).Error
}

func unpackSqlResult(rows *sql.Rows) ([]*Article, error) {
	articles := make([]*Article, 0)
	var err error
	var lastId int64 = 0
	var lastA *Article
	var tagName string
	var tagId int64
	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Id,
			&a.Title,
			&a.Publish,
			&a.Summary,
			&tagId,
			&tagName)
		if a.Id != lastId {
			articles = append(articles, a)
			lastA = a
			a.Tags = make([]*Tag, 0)
		} else {
			a = lastA
		}
		if err != nil {
			return nil, err
		}
		if tagId != 0 {
			a.Tags = append(a.Tags, &Tag{
				Id:   tagId,
				Name: tagName,
			})
		}
	}
	return articles, nil
}
