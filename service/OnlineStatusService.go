package service

import (
	"context"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"user/pkg/redis"
)

const (
	onlineKeyPrefix = "user:online:"
	onlineTTL       = 5 * time.Minute
)

func heartbeatUser(uid int) {
	if !redis.Enabled() {
		return
	}
	redis.Client().Set(context.Background(), onlineKeyPrefix+strconv.Itoa(uid), "1", onlineTTL)
}

func IsUserOnline(uid int) bool {
	if !redis.Enabled() {
		return false
	}
	n, _ := redis.Client().Exists(context.Background(), onlineKeyPrefix + strconv.Itoa(uid)).Result()
	return n > 0
}

func BatchIsOnline(uids []int) map[int]bool {
	result := make(map[int]bool, len(uids))
	if !redis.Enabled() || len(uids) == 0 {
		return result
	}
	ctx := context.Background()
	pipe := redis.Client().Pipeline()
	cmds := make([]*goredis.IntCmd, len(uids))
	for i, uid := range uids {
		cmds[i] = pipe.Exists(ctx, onlineKeyPrefix+strconv.Itoa(uid))
	}
	pipe.Exec(ctx)
	for i, uid := range uids {
		result[uid] = cmds[i].Val() > 0
	}
	return result
}
