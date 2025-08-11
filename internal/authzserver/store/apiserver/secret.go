package apiserver

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go"

	pb "github.com/pachirode/iam_study/internal/pkg/api/proto/apiserver/v1"
	"github.com/pachirode/iam_study/pkg/errors"
	"github.com/pachirode/iam_study/pkg/log"
)

type secrets struct {
	cli pb.CacheClient
}

func newSecrets(ds *dataStore) *secrets {
	return &secrets{ds.cli}
}

func (s *secrets) List() (map[string]*pb.SecretInfo, error) {
	secrets := make(map[string]*pb.SecretInfo)

	log.Info("Loading secrets")

	req := &pb.ListSecretsRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *pb.ListSecretsResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = s.cli.ListSecrets(context.Background(), req)

			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3))
	if err != nil {
		return nil, errors.Wrap(err, "List secrets failed")
	}

	log.Infof("Secrets found (%d total):", len(resp.Items))

	for _, v := range resp.Items {
		log.Infof(" - %s:%s", v.Username, v.SecretId)
		secrets[v.SecretId] = v
	}

	return secrets, nil
}
