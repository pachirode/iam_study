package pump

import (
	"context"
	"sync"
	"time"

	"github.com/pachirode/iam_study/internal/pump/analytics"
	"github.com/pachirode/iam_study/internal/pump/pumps"
	"github.com/pachirode/iam_study/pkg/log"
)

func writeToPumps(keys []interface{}, purgeDelay int) {
	if pumpList != nil {
		var wg sync.WaitGroup
		wg.Add(len(pumpList))
		for _, pump := range pumpList {
			go execPumpWriting(&wg, pump, &keys, purgeDelay)
		}
		wg.Wait()
	} else {
		log.Warn("No pumps defined!")
	}
}

func filterData(pump pumps.Pump, keys []interface{}) []interface{} {
	filters := pump.GetFilters()
	if !filters.HasFilter() && !pump.GetOmitDetailedRecording() {
		return keys
	}

	filteredKeys := keys[:]
	newLength := 0

	for _, key := range filteredKeys {
		decoded, _ := key.(analytics.AnalyticsRecord)
		if pump.GetOmitDetailedRecording() {
			decoded.Policies = ""
			decoded.Deciders = ""
		}
		if filters.ShouldFilter(decoded) {
			continue
		}
		filteredKeys[newLength] = decoded
		newLength++
	}
	filteredKeys = filteredKeys[:newLength]

	return filteredKeys
}

func execPumpWriting(wg *sync.WaitGroup, pump pumps.Pump, keys *[]interface{}, purgeDelay int) {
	timer := time.AfterFunc(time.Duration(purgeDelay)*time.Second, func() {
		if pump.GetTimeout() == 0 {
			log.Warnf(
				"Pump %s is taking more time than the value configured of purge_delay. You should try to set a timeout for this pump",
				pump.GetName(),
			)
		} else if pump.GetTimeout() > purgeDelay {
			log.Warnf("Pump %s is taking more time than the value configured of purge_delay. You should try lowering the timeout configured for this pump", pump.GetName())
		}
	})

	defer timer.Stop()
	defer wg.Done()

	log.Debugf("Writing to: %s", pump.GetName())

	ch := make(chan error, 1)
	var ctx context.Context
	var cancel context.CancelFunc

	if tm := pump.GetTimeout(); tm > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(tm)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	defer cancel()

	go func(ch chan error, ctx context.Context, pump pumps.Pump, keys *[]interface{}) {
		filteredKeys := filterData(pump, *keys)

		ch <- pump.WriteData(ctx, filteredKeys)
	}(ch, ctx, pump, keys)

	select {
	case err := <-ch:
		if err != nil {
			log.Warnf("Error Writing to: %s - Error: %s", pump.GetName(), err.Error())
		}
	case <-ctx.Done():
		switch ctx.Err() {
		case context.Canceled:
			log.Warnf("The writing to %s has got canceled.", pump.GetName())
		case context.DeadlineExceeded:
			log.Warnf("Timeout writing to: %s", pump.GetName())
		}
	}
}
