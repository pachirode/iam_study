package storage

import (
	"fmt"

	"github.com/go-redis/redis/v7"

	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
)

// GetKey will retrieve a key from the database.
func (r *RedisCluster) GetKey(keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}

	cluster := r.singleton()

	value, err := cluster.Get(r.fixKey(keyName)).Result()
	if err != nil {
		log.Debugf("Error trying to get value: %s", err.Error())

		return "", ErrKeyNotFound
	}

	return value, nil
}

// GetMultiKey gets multiple keys from the database.
func (r *RedisCluster) GetMultiKey(keys []string) ([]string, error) {
	if err := r.up(); err != nil {
		return nil, err
	}
	cluster := r.singleton()
	keyNames := make([]string, len(keys))
	copy(keyNames, keys)
	for index, val := range keyNames {
		keyNames[index] = r.fixKey(val)
	}

	result := make([]string, 0)

	switch v := cluster.(type) {
	case *redis.ClusterClient:
		{
			getCmds := make([]*redis.StringCmd, 0)
			pipe := v.Pipeline()
			for _, key := range keyNames {
				getCmds = append(getCmds, pipe.Get(key))
			}
			_, err := pipe.Exec()
			if err != nil && !errors.Is(err, redis.Nil) {
				log.Debugf("Error trying to get value: %s", err.Error())

				return nil, ErrKeyNotFound
			}
			for _, cmd := range getCmds {
				result = append(result, cmd.Val())
			}
		}
	case *redis.Client:
		{
			values, err := cluster.MGet(keyNames...).Result()
			if err != nil {
				log.Debugf("Error trying to get value: %s", err.Error())

				return nil, ErrKeyNotFound
			}
			for _, val := range values {
				strVal := fmt.Sprint(val)
				if strVal == "<nil>" {
					strVal = ""
				}
				result = append(result, strVal)
			}
		}
	}

	for _, val := range result {
		if val != "" {
			return result, nil
		}
	}

	return nil, ErrKeyNotFound
}

// GetKeyTTL return ttl of the given key.
func (r *RedisCluster) GetKeyTTL(keyName string) (ttl int64, err error) {
	if err = r.up(); err != nil {
		return 0, err
	}
	duration, err := r.singleton().TTL(r.fixKey(keyName)).Result()

	return int64(duration.Seconds()), err
}

// GetRawKey return the value of the given key.
func (r *RedisCluster) GetRawKey(keyName string) (string, error) {
	if err := r.up(); err != nil {
		return "", err
	}
	value, err := r.singleton().Get(keyName).Result()
	if err != nil {
		log.Debugf("Error trying to get value: %s", err.Error())

		return "", ErrKeyNotFound
	}

	return value, nil
}

// GetExp return the expiry of the given key.
func (r *RedisCluster) GetExp(keyName string) (int64, error) {
	log.Debugf("Getting exp for key: %s", r.fixKey(keyName))
	if err := r.up(); err != nil {
		return 0, err
	}

	value, err := r.singleton().TTL(r.fixKey(keyName)).Result()
	if err != nil {
		log.Errorf("Error trying to get TTL: ", err.Error())

		return 0, ErrKeyNotFound
	}

	return int64(value.Seconds()), nil
}

// GetKeys will return all keys according to the filter (filter is a prefix - e.g. tyk.keys.*).
func (r *RedisCluster) GetKeys(filter string) []string {
	if err := r.up(); err != nil {
		return nil
	}
	client := r.singleton()

	filterHash := ""
	if filter != "" {
		filterHash = r.hashKey(filter)
	}
	searchStr := r.KeyPrefix + filterHash + "*"
	log.Debugf("[STORE] Getting list by: %s", searchStr)

	fnFetchKeys := func(client *redis.Client) ([]string, error) {
		values := make([]string, 0)

		iter := client.Scan(0, searchStr, 0).Iterator()
		for iter.Next() {
			values = append(values, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return nil, err
		}

		return values, nil
	}

	var err error
	var values []string
	sessions := make([]string, 0)

	switch v := client.(type) {
	case *redis.ClusterClient:
		ch := make(chan []string)

		go func() {
			err = v.ForEachMaster(func(client *redis.Client) error {
				values, err = fnFetchKeys(client)
				if err != nil {
					return err
				}

				ch <- values

				return nil
			})
			close(ch)
		}()

		for res := range ch {
			sessions = append(sessions, res...)
		}
	case *redis.Client:
		sessions, err = fnFetchKeys(v)
	}

	if err != nil {
		log.Errorf("Error while fetching keys: %s", err)

		return nil
	}

	for i, v := range sessions {
		sessions[i] = r.cleanKey(v)
	}

	return sessions
}

// GetKeysAndValuesWithFilter will return all keys and their values with a filter.
func (r *RedisCluster) GetKeysAndValuesWithFilter(filter string) map[string]string {
	if err := r.up(); err != nil {
		return nil
	}
	keys := r.GetKeys(filter)
	if keys == nil {
		log.Error("Error trying to get filtered client keys")

		return nil
	}

	if len(keys) == 0 {
		return nil
	}

	for i, v := range keys {
		keys[i] = r.KeyPrefix + v
	}

	client := r.singleton()
	values := make([]string, 0)

	switch v := client.(type) {
	case *redis.ClusterClient:
		{
			getCmds := make([]*redis.StringCmd, 0)
			pipe := v.Pipeline()
			for _, key := range keys {
				getCmds = append(getCmds, pipe.Get(key))
			}
			_, err := pipe.Exec()
			if err != nil && !errors.Is(err, redis.Nil) {
				log.Errorf("Error trying to get client keys: %s", err.Error())

				return nil
			}

			for _, cmd := range getCmds {
				values = append(values, cmd.Val())
			}
		}
	case *redis.Client:
		{
			result, err := v.MGet(keys...).Result()
			if err != nil {
				log.Errorf("Error trying to get client keys: %s", err.Error())

				return nil
			}

			for _, val := range result {
				strVal := fmt.Sprint(val)
				if strVal == "<nil>" {
					strVal = ""
				}
				values = append(values, strVal)
			}
		}
	}

	m := make(map[string]string)
	for i, v := range keys {
		m[r.cleanKey(v)] = values[i]
	}

	return m
}

// GetKeysAndValues will return all keys and their values - not to be used lightly.
func (r *RedisCluster) GetKeysAndValues() map[string]string {
	return r.GetKeysAndValuesWithFilter("")
}
