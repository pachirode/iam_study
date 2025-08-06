package secret

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/iam_study/internal/pkg/middleware"
	"github.com/pachirode/iam_study/pkg/core"
	"github.com/pachirode/iam_study/pkg/log"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

// Delete delete a secret by the secret identifier.
func (s *SecretController) Delete(c *gin.Context) {
	log.L(c).Info("delete secret function called.")
	opts := metaV1.DeleteOptions{Unscoped: true}
	if err := s.serve.Secrets().Delete(c, c.GetString(middleware.UsernameKey), c.Param("name"), opts); err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
