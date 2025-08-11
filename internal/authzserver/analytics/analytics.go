package analytics

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/storage"
)

const analyticsKeyName = "iam-system-analytics"

const (
	recordsBufferForcedFlushInterval = 1 * time.Second
)

type AnalyticsRecord struct {
	TimeStamp  int64     `json:"timestamp"`
	Username   string    `json:"username"`
	Effect     string    `json:"effect"`
	Conclusion string    `json:"conclusion"`
	Request    string    `json:"request"`
	Policies   string    `json:"policies"`
	Deciders   string    `json:"deciders"`
	ExpireAt   time.Time `json:"expireAt"   bson:"expireAt"`
}

var analytics *Analytics

func (a *AnalyticsRecord) SetExpiry(expiresInSeconds int64) {
	expiry := time.Duration(expiresInSeconds) * time.Second
	if expiresInSeconds == 0 {
		expiry = 24 * 365 * 100 * time.Hour
	}

	t := time.Now()
	t2 := t.Add(expiry)
	a.ExpireAt = t2
}

type Analytics struct {
	store                      storage.AnalyticsHandler
	poolSize                   int
	recordsChan                chan *AnalyticsRecord
	workerBufferSize           uint64
	recordsBufferFlushInternal uint64
	shouldStop                 uint32
	poolWg                     sync.WaitGroup
}

func NewAnalytics(opts *AnalyticsOptions, store storage.AnalyticsHandler) *Analytics {
	ps := opts.PoolSize
	recordsBufferSize := opts.RecordsBufferSize
	workerBufferSize := recordsBufferSize / uint64(ps)
	log.Debug("Analytics pool worker buffer size", log.Uint64("workerBufferSize", workerBufferSize))

	recordsChan := make(chan *AnalyticsRecord, recordsBufferSize)

	analytics = &Analytics{
		store:                      store,
		poolSize:                   ps,
		recordsChan:                recordsChan,
		workerBufferSize:           workerBufferSize,
		recordsBufferFlushInternal: opts.FlushInterval,
	}

	return analytics
}

func GetAnalytics() *Analytics {
	return analytics
}

func (r *Analytics) Start() {
	r.store.Connect()

	atomic.SwapUint32(&r.shouldStop, 0)
	for i := 0; i < r.poolSize; i++ {
		r.poolWg.Add(1)
		go r.recordWorker()
	}
}

func (r *Analytics) Stop() {
	atomic.SwapUint32(&r.shouldStop, 1)
	close(r.recordsChan)
	r.poolWg.Wait()
}

// RecordHit will store AnalyticeRecord in redis
func (r *Analytics) RecordHit(record *AnalyticsRecord) error {
	if atomic.LoadUint32(&r.shouldStop) > 0 {
		return nil
	}

	r.recordsChan <- record

	return nil
}

func (r *Analytics) recordWorker() {
	defer r.poolWg.Done()

	// buffer to send pipelined command to redis
	recordsBuffer := make([][]byte, 0, r.workerBufferSize)

	lastSentTS := time.Now()
	for {
		var readyToSend bool
		select {
		case record, ok := <-r.recordsChan:
			if !ok {
				// send left in buffer
				r.store.AppendToSetPipelined(analyticsKeyName, recordsBuffer)

				return
			}

			if encoded, err := msgpack.Marshal(record); err != nil {
				log.Errorf("Error encoding analytics data: %s", err.Error())
			} else {
				recordsBuffer = append(recordsBuffer, encoded)
			}

			// identify buffer is ready to send
			readyToSend = uint64(len(recordsBuffer)) == r.workerBufferSize
		case <-time.After(time.Duration(r.recordsBufferFlushInternal) * time.Millisecond):
			readyToSend = true
		}
		if len(recordsBuffer) > 0 && (readyToSend || time.Since(lastSentTS) >= recordsBufferForcedFlushInterval) {
			r.store.AppendToSetPipelined(analyticsKeyName, recordsBuffer)
			recordsBuffer = recordsBuffer[:0]
			lastSentTS = time.Now()
		}
	}
}
