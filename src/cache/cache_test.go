package cache

import (
	"chainmaker_web/src/db"
	"testing"
)

func TestInitRedis(t *testing.T) {
	redisCfg, err := db.InitRedisContainer()
	if err != nil || redisCfg == nil {
		return
	}
	InitRedis(redisCfg)
}
