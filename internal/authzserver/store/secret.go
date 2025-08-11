package store

import (
	pb "github.com/pachirode/iam_study/internal/pkg/api/proto/apiserver/v1"
)

type SecretStore interface {
	List() (map[string]*pb.SecretInfo, error)
}
