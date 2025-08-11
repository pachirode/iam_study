package cache

import (
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/ory/ladon"

	"github.com/pachirode/iam_study/internal/authzserver/store"
	pb "github.com/pachirode/iam_study/internal/pkg/api/proto/apiserver/v1"
	"github.com/pachirode/iam_study/pkg/errors"
)

type Cache struct {
	lock     *sync.RWMutex
	cli      store.Factory
	secrets  *ristretto.Cache
	policies *ristretto.Cache
}

var (
	ErrSecretNotFound = errors.New("Secret not found")
	ErrPolicyNotFound = errors.New("Policy not found")
)

var (
	onceCache sync.Once
	cacheIns  *Cache
)

func GetCacheInsOr(cli store.Factory) (*Cache, error) {
	var err error
	if cli != nil {
		var (
			secretCache *ristretto.Cache
			policyCache *ristretto.Cache
		)

		onceCache.Do(func() {
			config := &ristretto.Config{
				NumCounters: 1e7,     // track frequency 10M
				MaxCost:     1 << 30, // 1GB
				BufferItems: 64,
				Cost:        nil,
			}

			secretCache, err = ristretto.NewCache(config)
			if err != nil {
				return
			}

			policyCache, err = ristretto.NewCache(config)
			if err != nil {
				return
			}

			cacheIns = &Cache{
				cli:      cli,
				lock:     new(sync.RWMutex),
				secrets:  secretCache,
				policies: policyCache,
			}
		})
	}

	return cacheIns, err
}

func (c *Cache) GetSecret(key string) (*pb.SecretInfo, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.secrets.Get(key)
	if !ok {
		return nil, ErrSecretNotFound
	}

	return value.(*pb.SecretInfo), nil
}

func (c *Cache) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.policies.Get(key)
	if !ok {
		return nil, ErrPolicyNotFound
	}

	return value.([]*ladon.DefaultPolicy), nil
}

func (c *Cache) Reload() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	secrets, err := c.cli.Secrets().List()
	if err != nil {
		return errors.Wrap(err, "List secrets failed")
	}

	c.secrets.Clear()
	for key, val := range secrets {
		c.secrets.Set(key, val, 1)
	}

	policies, err := c.cli.Policies().List()
	if err != nil {
		return errors.Wrap(err, "List policies failed")
	}

	c.policies.Clear()
	for key, val := range policies {
		c.policies.Set(key, val, 1)
	}

	return nil
}
