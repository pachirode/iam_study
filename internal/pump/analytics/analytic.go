package analytics

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AnalyticsRecord struct {
	TimeStamp  int64     `json:"timestamp"`
	Username   string    `json:"username"`
	Effect     string    `json:"effect"`
	Conclusion string    `json:"conclusion"`
	Request    string    `json:"request"`
	Policies   string    `json:"policies"`
	Deciders   string    `json:"deciders"`
	ExpireAt   time.Time `json:"expireAt" bson:"expireAt"`
}

func (a *AnalyticsRecord) GetFieldNames() []string {
	val := reflect.ValueOf(a).Elem()
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fields = append(fields, typeField.Name)
	}

	return fields
}

func (a *AnalyticsRecord) GetLineValues() []string {
	val := reflect.ValueOf(a).Elem()
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		var thisVal string
		switch typeField.Type.String() {
		case "int":
			thisVal = strconv.Itoa(int(valueField.Int()))
		case "int64":
			thisVal = strconv.Itoa(int(valueField.Int()))
		case "[]string":
			tmpVal, _ := valueField.Interface().([]string)
			thisVal = strings.Join(tmpVal, ";")
		case "time.Time":
			tmpVal, _ := valueField.Interface().(time.Time)
			thisVal = tmpVal.String()
		case "time.Month":
			tmpVal, _ := valueField.Interface().(time.Month)
			thisVal = tmpVal.String()
		default:
			thisVal = valueField.String()
		}

		fields = append(fields, thisVal)
	}

	return fields
}
