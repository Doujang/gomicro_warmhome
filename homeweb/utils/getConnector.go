package utils

import (
	"encoding/json"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

func GetRedisConnector() (cache.Cache, error) {
	//连接redis创建句柄
	redis_config_map := map[string]string{
		"key":      G_server_name,
		"conn":     G_redis_addr + ":" + G_redis_port,
		"dbNum":    G_redis_dbnum,
		"password": G_redis_passwd,
	}
	//将map转化为json
	redis_config, _ := json.Marshal(redis_config_map)
	//连接redis
	bm, err := cache.NewCache("redis", string(redis_config))
	return bm, err
}
