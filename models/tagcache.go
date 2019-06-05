package models

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

/**
提供tag缓存支持
*/
const (
	tagPrefix = "tag:"
	tagSet    = "tagids"
)

//将所有的文章标签缓存至redis中
func loadTags() error {
	tags, err := queryAllTagsInSql()
	if err != nil {
		return err
	}
	conn := client.Get()
	for _, tag := range tags {
		b, _ := json.Marshal(tag)
		tagval := string(b)
		tagkey := tagIdKey(tag.Id)
		err = conn.Send("SET", tagkey, tagval)
		if err == nil {
			_ = conn.Send("SADD", tagSet, tag.Id)
		}
	}
	_ = conn.Flush()
	return nil
}

func getTagsInCache() ([]*Tag, error) {
	tags := make([]*Tag, 0)
	conn := client.Get()
	ids, err := redis.Int64s(conn.Do("SMEMBERS", tagSet))
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		key := tagIdKey(id)
		r, err := conn.Do("GET", key)
		if err != nil {
			return nil, err
		}
		if r == nil {
			continue
		}
		tag := new(Tag)
		_ = json.Unmarshal(r.([]byte), tag)
		tags = append(tags, tag)
	}

	return tags, nil
}

func queryTagByIdInCache(id int64) (*Tag, error) {
	conn := client.Get()
	tag := new(Tag)
	key := tagIdKey(id)

	r, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	_ = json.Unmarshal(r.([]byte), tag)

	return tag, nil
}

func queryTagsByIdsInCache(ids []int64) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	conn := client.Get()
	for _, id := range ids {
		key := tagIdKey(id)
		r, err := conn.Do("GET", key)
		if err != nil {
			return nil, err
		}

		if r == nil {
			continue
		}

		tag := new(Tag)
		_ = json.Unmarshal(r.([]byte), tag)
		tags = append(tags, tag)
	}
	return tags, nil
}

func clearTagInCache(id int64) {
	conn := client.Get()
	_ = conn.Send("MUTLI")
	_ = conn.Send("SREM", tagSet, id)
	_ = conn.Send("DEL", tagIdKey(id))
	_ = conn.Send("SREM", tagSet, id)
	_ = conn.Send("EXEC")
	_ = conn.Flush()
}

func cacheTagsInCache(t *Tag) error {
	conn := client.Get()
	cnt, err := redis.Int(conn.Do("SADD", tagSet, t.Id))
	if err != nil {
		return err
	}

	if cnt > 0 {
		b, _ := json.Marshal(t)
		tagval := string(b)
		tagkey := tagIdKey(t.Id)
		_, err := conn.Do("SET", tagkey, tagval)
		if err != nil {
			return err
		}
	}

	return nil
}

func tagIdKey(id int64) string {
	tagkey :=
		fmt.Sprintf("%s%s", tagPrefix, strconv.FormatInt(id, 10))
	return tagkey
}
