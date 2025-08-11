package storage

import (
	"time"

	"github.com/spf13/viper"

	"github.com/pachirode/iam_study/pkg/log"
)

// AppendToSetPipelined append values to redis pipeline.
func (r *RedisCluster) AppendToSetPipelined(key string, values [][]byte) {
	if len(values) == 0 {
		return
	}

	fixedKey := r.fixKey(key)
	if err := r.up(); err != nil {
		log.Debug(err.Error())

		return
	}
	client := r.singleton()

	pipe := client.Pipeline()
	for _, val := range values {
		pipe.RPush(fixedKey, val)
	}

	if _, err := pipe.Exec(); err != nil {
		log.Errorf("Error trying to append to set keys: %s", err.Error())
	}

	// if we need to set an expiration time
	if storageExpTime := int64(viper.GetDuration("analytics.storage-expiration-time")); storageExpTime != int64(-1) {
		// If there is no expiry on the analytics set, we should set it.
		exp, _ := r.GetExp(key)
		if exp == -1 {
			_ = r.SetExp(key, time.Duration(storageExpTime)*time.Second)
		}
	}
}
