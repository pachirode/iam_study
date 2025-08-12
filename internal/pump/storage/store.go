package storage

const (
	AnalyticsKeyName string = "iam-system-analytics"
)

type AnalyticsStore interface {
	Init(config interface{}) error
	GetName() string
	Connect() bool
	GetAndDeleteSet(string) []interface{}
}
