package storage

import (
	"context"
	"crypto/tls"
	"sync/atomic"
	"time"

	redis "github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"

	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
)

type Config struct {
	Host                  string
	Port                  int
	Addrs                 []string
	MasterName            string
	Username              string
	Password              string
	Database              int
	MaxIdle               int
	MaxActive             int
	Timeout               int
	EnableCluster         bool
	UseSSL                bool
	SSLInsecureSkipVerify bool
}

type RedisOpts redis.UniversalOptions

var ErrRedisIsDown = errors.New("storage: Redis is either down or ws not configured")

var (
	singlePool      atomic.Value
	singleCachePool atomic.Value
	redisUp         atomic.Value
)

var disableRedis atomic.Value

// DisableRedis very handy when testsing it allows to dynamically enable/disable talking with redisW.
func DisableRedis(ok bool) {
	if ok {
		redisUp.Store(false)
		disableRedis.Store(true)

		return
	}
	redisUp.Store(true)
	disableRedis.Store(false)
}

func shouldConnect() bool {
	if v := disableRedis.Load(); v != nil {
		return !v.(bool)
	}

	return true
}

func Connected() bool {
	if v := redisUp.Load(); v != nil {
		return v.(bool)
	}

	return false
}

func singleton(cache bool) redis.UniversalClient {
	if cache {
		v := singleCachePool.Load()
		if v != nil {
			return v.(redis.UniversalClient)
		}

		return nil
	}

	if v := singlePool.Load(); v != nil {
		return v.(redis.UniversalClient)
	}

	return nil
}

// nolint: unparam
func connectSingleton(cache bool, config *Config) bool {
	if singleton(cache) == nil {
		log.Debug("Connecting to redis cluster")
		if cache {
			singleCachePool.Store(NewRedisClusterPool(cache, config))

			return true
		}
		singlePool.Store(NewRedisClusterPool(cache, config))

		return true
	}

	return true
}

type RedisCluster struct {
	KeyPrefix string
	HashKeys  bool
	IsCache   bool
}

func clusterConnectionIsOpen(cluster RedisCluster) bool {
	c := singleton(cluster.IsCache)
	testKey := "redis-test-" + uuid.Must(uuid.NewV4()).String()
	if err := c.Set(testKey, "test", time.Second).Err(); err != nil {
		//		log.Warnf("Error trying to set test key: %s", err.Error())

		return false
	}

	return true
}

func ConnectToRedis(ctx context.Context, config *Config) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()

	c := []RedisCluster{
		{}, {IsCache: true},
	}

	var ok bool
	for _, v := range c {
		if !connectSingleton(v.IsCache, config) {
			break
		}

		if !clusterConnectionIsOpen(v) {
			redisUp.Store(false)

			break
		}

		ok = true
	}
	redisUp.Store(ok)
again:
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if !shouldConnect() {
				continue
			}
			for _, v := range c {
				if !connectSingleton(v.IsCache, config) {
					redisUp.Store(false)

					goto again
				}

				if !clusterConnectionIsOpen(v) {
					redisUp.Store(false)

					goto again
				}
			}
			redisUp.Store(true)
		}
	}
}

func NewRedisClusterPool(isCache bool, config *Config) redis.UniversalClient {
	log.Debug("Creating new Redis connection pool")

	poolSize := 500
	if config.MaxActive > 0 {
		poolSize = config.MaxActive
	}

	timeout := 5 * time.Second

	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	var tlsConfig *tls.Config

	if config.UseSSL {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: config.SSLInsecureSkipVerify,
		}
	}

	var client redis.UniversalClient
	opts := &RedisOpts{
		Addrs:        getRedisAddrs(config),
		MasterName:   config.MasterName,
		Password:     config.Password,
		DB:           config.Database,
		DialTimeout:  timeout,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  240 * timeout,
		PoolSize:     poolSize,
		TLSConfig:    tlsConfig,
	}

	if opts.MasterName != "" {
		log.Info("--> [REDIS] Creating sentinel-backed failover client")
		client = redis.NewFailoverClient(opts.failover())
	} else if config.EnableCluster {
		log.Info("--> [REDIS] Creating cluster client")
		client = redis.NewClusterClient(opts.cluster())
	} else {
		log.Info("--> [REDIS] Creating single-node client")
		client = redis.NewClient(opts.simple())
	}

	return client
}

func (opt *RedisOpts) cluster() *redis.ClusterOptions {
	if len(opt.Addrs) == 0 {
		opt.Addrs = []string{"127.0.0.1:6379"}
	}

	return &redis.ClusterOptions{
		Addrs:     opt.Addrs,
		OnConnect: opt.OnConnect,

		Password: opt.Password,

		MaxRedirects:   opt.MaxRedirects,
		ReadOnly:       opt.ReadOnly,
		RouteByLatency: opt.RouteByLatency,
		RouteRandomly:  opt.RouteRandomly,

		MaxRetries:      opt.MaxRetries,
		MinRetryBackoff: opt.MinRetryBackoff,
		MaxRetryBackoff: opt.MaxRetryBackoff,

		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,

		TLSConfig: opt.TLSConfig,
	}
}

func (o *RedisOpts) simple() *redis.Options {
	addr := "127.0.0.1:6379"
	if len(o.Addrs) > 0 {
		addr = o.Addrs[0]
	}

	return &redis.Options{
		Addr:      addr,
		OnConnect: o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}

func (o *RedisOpts) failover() *redis.FailoverOptions {
	if len(o.Addrs) == 0 {
		o.Addrs = []string{"127.0.0.1:26379"}
	}

	return &redis.FailoverOptions{
		SentinelAddrs: o.Addrs,
		MasterName:    o.MasterName,
		OnConnect:     o.OnConnect,

		DB:       o.DB,
		Password: o.Password,

		MaxRetries:      o.MaxRetries,
		MinRetryBackoff: o.MinRetryBackoff,
		MaxRetryBackoff: o.MaxRetryBackoff,

		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,

		PoolSize:           o.PoolSize,
		MinIdleConns:       o.MinIdleConns,
		MaxConnAge:         o.MaxConnAge,
		PoolTimeout:        o.PoolTimeout,
		IdleTimeout:        o.IdleTimeout,
		IdleCheckFrequency: o.IdleCheckFrequency,

		TLSConfig: o.TLSConfig,
	}
}
