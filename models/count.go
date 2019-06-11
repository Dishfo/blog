package models

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"strings"
	"time"
)

/**
提供	对接口访问的计数
记录	在各个时间段内的访问次数

know:path set {

}

1:base {
	cap:count
}

需要一个集合用于记录
定期把内存中的记录写向disk
也就是移除某些区间的记录写往持久层
一般来说这些被移除的记录都是不可能再变化的

*/

const (
	knowSet = "know"
)

var (
	precisions = []int{
		5, 60, 3600,
	}
	cleanExeTimes = 0
)

//参	数中的timeStamp 是
func AccessSite(base string, timestamp int64) error {
	//修改记录
	for _, prec := range precisions {
		conn := client.Get()
		countKey := fmt.Sprintf("%d:%s", prec, base)
		i := int64(prec) * (timestamp / int64(prec))
		_, _ = conn.Do("zadd", knowSet, countKey, 0)
		_, _ = conn.Do("hincrby", "count:"+countKey, i, 1)
	}
	return nil
}

func GetAccessCount(site string, prec int) (map[int64]int, error) {
	conn := client.Get()
	key := fmt.Sprintf("count:%d:%s", prec, site)
	counts, err := redis.IntMap(conn.Do("hgetall", key))
	if err != nil {
		return nil, err
	}

	if counts == nil {
		return nil, nil
	}

	result := make(map[int64]int)
	for k, v := range counts {
		s, _ := strconv.ParseInt(k, 10, 64)
		result[s] = v
	}

	return result, nil
}

//清除过久的数据部分 考虑把这部分数据写往磁盘,
//保证始终只有一个协程在运行这个函数
func cleanCount() {
	defer func() {
		cleanExeTimes++
	}()
	conn := client.Get()
	keys, err := redis.Strings(conn.Do("zrange", knowSet, 0, -1))
	if err != nil {
		return
	}

	for _, key := range keys {
		precStr := strings.Split(key, ":")[0]
		prec, _ := strconv.ParseInt(precStr, 10, 64)
		bprec := prec / 60
		if cleanExeTimes%int(bprec) == 0 {
			continue
		}

		fields, err := redis.Strings(conn.Do("hkeys",
			"count:"+key))

		if err != nil {
			return
		}

		toDel := checkDelCap(fields, int(prec))
		args := redis.Args{}
		for _, s := range toDel {
			args = args.Add(s)
		}

		vals, err := redis.Ints(conn.Do("hmget", "count:"+key,
			args))
		if err != nil {
			return
		}

		_, _ = conn.Do("hdel", "count:"+key, args)

		m := make(map[string]int)
		for i, v := range vals {
			m[toDel[i]] = v
		}

		afterClean(m, key)
	}
}

func checkDelCap(fields []string, prec int) []string {
	toDel := make([]string, 0)
	timeStamp := time.Now().Second()
	end := timeStamp / prec * prec
	for _, s := range fields {
		i, _ := strconv.Atoi(s)
		if i <= end {
			toDel = append(toDel, s)
		}
	}
	return toDel
}

/*
在移除访问计数时调用
把移除的计数保存到本地

*/
func afterClean(counts map[string]int, key string) {

}
