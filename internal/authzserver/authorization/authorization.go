package authorization

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ory/ladon"

	"github.com/pachirode/iam_study/internal/authzserver/analytics"
)

type Authorization struct {
	getter PolicyGetter
}

var _ AuthorizationInterface = (*Authorization)(nil)

func NewAuthorization(getter PolicyGetter) AuthorizationInterface {
	return &Authorization{getter: getter}
}

func (auth *Authorization) Create(policy *ladon.DefaultPolicy) error {
	return nil
}

func (auth *Authorization) Update(policy *ladon.DefaultPolicy) error {
	return nil
}

func (auth *Authorization) Delete(id string) error {
	return nil
}

func (auth *Authorization) DeleteCollection(idList []string) error {
	return nil
}

func (auth *Authorization) Get(id string) (*ladon.DefaultPolicy, error) {
	return &ladon.DefaultPolicy{}, nil
}

func (auth *Authorization) List(username string) ([]*ladon.DefaultPolicy, error) {
	return auth.getter.GetPolicy(username)
}

func (auth *Authorization) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	var conclusion string
	if len(d) > 1 {
		allowed := joinPoliciesNames(d[0 : len(d)-1])
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("Policies %s allow access, but policy %s forcefully denied it", allowed, denied)
	} else if len(d) == 1 {
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("Policy %s forcefully denied the access", denied)
	} else {
		conclusion = "No policy allowed access"
	}

	rString, pString, dString := convertToString(r, p, d)
	record := analytics.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.DenyAccess,
		Conclusion: conclusion,
		Request:    rString,
		Policies:   pString,
		Deciders:   dString,
	}

	record.SetExpiry(0)
	_ = analytics.GetAnalytics().RecordHit(&record)
}

func (auth *Authorization) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	conclusion := fmt.Sprintf("Policies %s allowed access", joinPoliciesNames(d))
	rString, dString, pString := convertToString(r, p, d)
	record := analytics.AnalyticsRecord{
		TimeStamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.AllowAccess,
		Conclusion: conclusion,
		Request:    rString,
		Policies:   pString,
		Deciders:   dString,
	}

	record.SetExpiry(0)
	_ = analytics.GetAnalytics().RecordHit(&record)
}

func joinPoliciesNames(policies ladon.Policies) string {
	names := []string{}
	for _, policy := range policies {
		names = append(names, policy.GetID())
	}

	return strings.Join(names, ", ")
}

func convertToString(r *ladon.Request, p ladon.Policies, d ladon.Policies) (string, string, string) {
	rBytes, _ := json.Marshal(r)
	pBytes, _ := json.Marshal(p)
	dBytes, _ := json.Marshal(d)

	return string(rBytes), string(pBytes), string(dBytes)
}
