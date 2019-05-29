package models

import "github.com/astaxie/beego/orm"

func queryArticleListInSql(pageno, size int) ([]*Article, error) {
	offset := pageno * size
	articles := make([]*Article, 0)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id", "title", "publish", "summary").
		From(articleTable).
		Limit(size).
		Offset(offset)
	sql := qb.String()
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return nil, err
	}
	defer o.Commit()
	_, err = o.Raw(sql).QueryRows(&articles)
	for _, a := range articles {
		_, err = o.LoadRelated(a, "Tags")
	}
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func queryArticleByIdInSql(id int64) (*Article, error) {
	a := new(Article)
	err :=
		orm.NewOrm().
			QueryTable(articleTable).
			Filter("Id", id).
			One(a)

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
	o := orm.NewOrm()
	_, err := o.QueryTable(articleTable).
		Filter("Tags__Tag__id__in", tagIds).
		Distinct().
		All(&articles)
	return articles, err
}

func insertArticleInSql(a *Article) error {
	id, err := orm.NewOrm().Insert(a)
	if err != nil {
		return err
	}
	a.Id = id
	return nil
}

func updateArticleInSql(a *Article, fields []string) error {
	o := orm.NewOrm()
	var err error
	err = o.Begin()
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
		_, err = o.Update(a, filterFields...)
	}

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

	return err
}

func deleteArticleInSql(id int64) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Article{
		Id: id,
	})
	return err
}
