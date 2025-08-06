package secret

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

// Get get an policy by the secret identifier.
func (s *SecretController) Get(c *gin.Context) {
	log.L(c).Info("get secret function called.")

	secret, err := s.serve.Secrets().Get(c, c.GetString(middleware.UsernameKey), c.Param("name"), metaV1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, secret)
}
