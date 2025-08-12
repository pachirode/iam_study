package storage

import (
	"strconv"

	genericOptions "github.com/pachirode/iam_study/internal/pkg/options"
)

func getRedisAddrs(config genericOptions.RedisOptions) (addrs []string) {
	if len(config.Addrs) != 0 {
		addrs = config.Addrs
	}

	if len(addrs) == 0 && config.Port != 0 {
		addr := config.Host + ":" + strconv.Itoa(config.Port)
		addrs = append(addrs, addr)
	}

	return addrs
}
