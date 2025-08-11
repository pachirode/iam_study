package load

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/go-redis/redis/v7"

	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
	"github.com/pachirode/iam_study/pkg/storage"
)

type NotificationCommand string

var (
	reloadedQueue = make(chan func())
	requestLock   sync.Mutex
	requeue       []func()
)

const (
	RedisPubSubChannel                      = "iam.cluster.notifications"
	NoticePolicyChanged NotificationCommand = "PolicyChanged"
	NoticeSecretChanged NotificationCommand = "SecretChanged"
)

type Notification struct {
	Command       NotificationCommand `json:"Command"`
	Payload       string              `json:"payload"`
	Signature     string              `json:"signature"`
	SignatureAlgo crypto.Hash         `json:"algorithm"`
}

type RedisNotifier struct {
	store   *storage.RedisCluster
	channel string
}

func (n *Notification) Sign() {
	n.SignatureAlgo = crypto.SHA256
	hash := sha256.Sum256([]byte(string(n.Command) + n.Payload))
	n.Signature = hex.EncodeToString(hash[:])
}

func handleRedisEvent(v interface{}, handled func(NotificationCommand), reloaded func()) {
	msg, ok := v.(*redis.Message)
	if !ok {
		return
	}

	notification := Notification{}
	if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
		log.Errorf("Unmarshalling message body failed, malformed: ", err)

		return
	}
	log.Infow("receive redis message", "command", notification.Command, "payload", msg.Payload)

	switch notification.Command {
	case NoticePolicyChanged, NoticeSecretChanged:
		log.Info("Reloading secrets and policies")
		reloadedQueue <- reloaded
	default:
		log.Warnf("Unknow notification command: %q", notification.Command)
		return
	}

	if handled != nil {
		handled(notification.Command)
	}
}

func (r *RedisNotifier) Notify(notification interface{}) bool {
	if n, ok := notification.(Notification); ok {
		n.Sign()
		notification = n
	}

	toSend, err := json.Marshal(notification)
	if err != nil {
		log.Errorf("Problem marshaling notification: %s", err.Error())

		return false
	}

	log.Debugf("Sending notification: %v", notification)

	if err := r.store.Publish(r.channel, string(toSend)); err != nil {
		if !errors.Is(err, storage.ErrRedisIsDown) {
			log.Errorf("Could not send notification: %s", err.Error())
		}

		return false
	}

	return true
}
