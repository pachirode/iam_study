package pump

import (
	"fmt"
	"time"

	goredislib "github.com/go-redis/redis/v7"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v7"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/pachirode/iam_study/internal/pump/analytics"
	"github.com/pachirode/iam_study/internal/pump/config"
	"github.com/pachirode/iam_study/internal/pump/options"
	"github.com/pachirode/iam_study/internal/pump/pumps"
	"github.com/pachirode/iam_study/internal/pump/storage"
	"github.com/pachirode/iam_study/internal/pump/storage/redis"
	"github.com/pachirode/iam_study/pkg/log"
)

var pumpList []pumps.Pump

type pumpServer struct {
	secInterval    int
	omitDetails    bool
	mutex          *redsync.Mutex
	analyticsStore storage.AnalyticsStore
	pumps          map[string]options.PumpConfig
}

type preparedPumpServer struct {
	*pumpServer
}

func createPumpServer(cfg *config.Config) (*pumpServer, error) {
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisOptions.Host, cfg.RedisOptions.Port),
		Username: cfg.RedisOptions.Username,
		Password: cfg.RedisOptions.Password,
	})

	rs := redsync.New(goredis.NewPool(client))

	server := &pumpServer{
		secInterval:    cfg.PurgeDelay,
		omitDetails:    cfg.OmitDetailedRecording,
		mutex:          rs.NewMutex("iam-pump", redsync.WithExpiry(10*time.Minute)),
		analyticsStore: &redis.RedisClusterStorageManager{},
		pumps:          cfg.Pumps,
	}

	if err := server.analyticsStore.Init(cfg.RedisOptions); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *pumpServer) PrepareRun() preparedPumpServer {
	s.initialize()

	return preparedPumpServer{s}
}

func (s preparedPumpServer) Run(stopCh <-chan struct{}) error {
	ticker := time.NewTicker(time.Duration(s.secInterval) * time.Second)
	defer ticker.Stop()

	log.Info("Now run loop to clean data from redis")
	for {
		select {
		case <-ticker.C:
			s.pump()
		case <-stopCh:
			log.Info("stop purge loop")

			return nil
		}
	}
}

// pump get authorization log from redis and write to pumps
func (s *pumpServer) pump() {
	if err := s.mutex.Lock(); err != nil {
		log.Info("There is already an iam-pump instance running.")

		return
	}

	defer func() {
		if _, err := s.mutex.Unlock(); err != nil {
			log.Errorf("Could not release iam-pump lock. err: %v", err)
		}
	}()

	analyticsValues := s.analyticsStore.GetAndDeleteSet(storage.AnalyticsKeyName)
	if len(analyticsValues) == 0 {
		return
	}

	keys := make([]interface{}, len(analyticsValues))

	for i, v := range analyticsValues {
		decoded := analytics.AnalyticsRecord{}
		err := msgpack.Unmarshal([]byte(v.(string)), &decoded)
		log.Debugf("Decoded Record: %v", decoded)
		if err != nil {
			log.Errorf("Couldn't unmarshal analytics data: %s", err.Error())
		} else {
			if s.omitDetails {
				decoded.Policies = ""
				decoded.Deciders = ""
			}
			keys[i] = interface{}(decoded)
		}
	}

	writeToPumps(keys, s.secInterval)
}

func (s *pumpServer) initialize() {
	pumpList = make([]pumps.Pump, len(s.pumps))
	i := 0
	for key, pump := range s.pumps {
		pumpTypeName := pump.Type
		if pumpTypeName == "" {
			pumpTypeName = key
		}

		pumpType, err := pumps.GetPumpByName(pumpTypeName)
		if err != nil {
			log.Errorf("Pump load error (skipping): %s", err.Error())
		} else {
			pumpIns := pumpType.New()
			initErr := pumpIns.Init(pump.Meta)
			if initErr != nil {
				log.Errorf("Pump init error (skipping): %s", initErr.Error())
			} else {
				log.Infof("Init Pump: %s", pumpIns.GetName())
				pumpIns.SetFilters(pump.Filters)
				pumpIns.SetTimeout(pump.Timeout)
				pumpIns.SetOmitDetailedRecording(pump.OmitDetailedRecording)
				pumpList[i] = pumpIns
			}
		}
		i++
	}
}
