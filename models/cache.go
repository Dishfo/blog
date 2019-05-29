package models

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	client *redis.Pool
)

/**
对于 tags 的存储,这是一个需要对外显示的字段
先使用json 转化为string 在存储
*/

func InitRedis() {
	host := beego.AppConfig.String("redis.server")
	pass := beego.AppConfig.String("redis.pass")
	client = &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", pass); err != nil {
				_ = c.Close()
				return nil, err
			}
			return c, nil
		},
		MaxIdle: 30,
	}

	r, err := client.Get().Do("PING")
	if err != nil {
		beego.BeeLogger.Error("%s", err.Error())
		return
	}
	_, _ = client.Get().Do("flushall")
	str := r.(string)
	beego.BeeLogger.Info("%s init redis connection pool", str)

	err = loadTags()
	if err != nil {
		beego.BeeLogger.Error("%s when load tags ", err.Error())
		return
	}

	//定期移除过期的token
	go func() {
		for {
			time.Sleep(time.Second * 120)
			cleanExpiredToken()
			clearAccessRecord()
			clearAccessSeq()
			clearArticles()
		}
	}()
}
