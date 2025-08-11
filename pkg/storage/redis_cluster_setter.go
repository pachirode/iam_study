package storage

import (
	"time"

	"github.com/pachirode/iam_study/pkg/log"
)

// SetExp set expiry of the given key.
func (r *RedisCluster) SetExp(keyName string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Expire(r.fixKey(keyName), timeout).Err()
	if err != nil {
		log.Errorf("Could not EXPIRE key: %s", err.Error())
	}

	return err
}

// SetKey will create (or update) a key value in the store.
func (r *RedisCluster) SetKey(keyName, session string, timeout time.Duration) error {
	log.Debugf("[STORE] SET Raw key is: %s", keyName)
	log.Debugf("[STORE] Setting key: %s", r.fixKey(keyName))

	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Set(r.fixKey(keyName), session, timeout).Err()
	if err != nil {
		log.Errorf("Error trying to set value: %s", err.Error())

		return err
	}

	return nil
}

// SetRawKey set the value of the given key.
func (r *RedisCluster) SetRawKey(keyName, session string, timeout time.Duration) error {
	if err := r.up(); err != nil {
		return err
	}
	err := r.singleton().Set(keyName, session, timeout).Err()
	if err != nil {
		log.Errorf("Error trying to set value: %s", err.Error())

		return err
	}

	return nil
}
