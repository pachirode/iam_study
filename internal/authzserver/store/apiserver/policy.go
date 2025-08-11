package apiserver

import (
	"context"
	"encoding/json"

	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go"
	"github.com/ory/ladon"

	pb "github.com/pachirode/iam_study/internal/pkg/api/proto/apiserver/v1"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
)

type policies struct {
	cli pb.CacheClient
}

func newPolicies(ds *dataStore) *policies {
	return &policies{cli: ds.cli}
}

func (p *policies) List() (map[string][]*ladon.DefaultPolicy, error) {
	pols := make(map[string][]*ladon.DefaultPolicy)

	log.Info("Loading policies")

	req := &pb.ListPoliciesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *pb.ListPoliciesResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = p.cli.ListPolicies(context.Background(), req)
			if listErr != nil {
				return listErr
			}
			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "List policies failed")
	}

	log.Infof("Policies found (%d total)[username:name]:", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s:%s", v.Username, v.Name)

		var policy ladon.DefaultPolicy
		if err := json.Unmarshal([]byte(v.PolicyShadow), &policy); err != nil {
			log.Warnf("Failed to load policy for %s, error: %s", v.Name, err.Error())

			continue
		}

		pols[v.Username] = append(pols[v.Username], &policy)
	}

	return pols, nil
}
