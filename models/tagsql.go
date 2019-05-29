package models

import "github.com/astaxie/beego/orm"

/**
提供从mysql数据库读写tag
*/
const (
	articleTable = "article"
	tagTable     = "tag"
)

func queryAllTagsInSql() ([]*Tag, error) {
	tags := make([]*Tag, 0)
	o := orm.NewOrm()
	_, err := o.QueryTable(tagTable).All(&tags)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func deleteTagInSql(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Tag{
		Id: id,
	}, "Id")
	return err
}

func insertTagInSql(tag *Tag) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(tag)
}

func queryTagsByIdInSql(ids []int64) ([]*Tag, error) {
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return nil, err
	}
	defer o.Commit()

	tags := make([]*Tag, 0)
	_, err = o.QueryTable(tagTable).
		Filter("Id__in", ids).
		All(&tags)
	return tags, err
}

func queryRelatedArticle(id int64) ([]*Article, error) {
	articles := make([]*Article, 0)
	o := orm.NewOrm()
	err := o.Begin()
	defer o.Commit()
	if err != nil {
		return articles, err
	}
	_, err = o.QueryTable(articleTable).Filter("Tags__Tag__Id", id).All(&articles)
	if err != nil {
		return articles, err
	}

	for _, a := range articles {
		_, err = o.LoadRelated(a, "Tags")
		if err != nil {
			return nil, err
		}
	}

	return articles, nil
}
