package authzserver

import (
	"github.com/pachirode/iam_study/internal/authzserver/load/cache"
	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/internal/pkg/middleware/auth"
	"github.com/pachirode/iam_study/pkg/errors"
)

func newCacheAuth() middleware.AuthStrategy {
	return auth.NewCacheStrategy(getSecretFunc())
}

func getSecretFunc() func(string) (auth.Secret, error) {
	return func(kid string) (auth.Secret, error) {
		cli, err := cache.GetCacheInsOr(nil)
		if err != nil || cli == nil {
			return auth.Secret{}, errors.Wrap(err, "get cache instance failed")
		}

		secret, err := cli.GetSecret(kid)
		if err != nil {
			return auth.Secret{}, err
		}

		return auth.Secret{
			Username: secret.Username,
			ID:       secret.SecretId,
			Key:      secret.SecretKey,
			Expires:  secret.Expires,
		}, nil
	}
}
