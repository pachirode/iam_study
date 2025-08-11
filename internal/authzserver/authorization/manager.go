package authorization

import (
	"context"

	"github.com/ory/ladon"

	"github.com/pachirode/iam_study/pkg/errors"
)

type PolicyManager struct {
	client AuthorizationInterface
}

func NewPolicyManager(client AuthorizationInterface) ladon.Manager {
	return &PolicyManager{
		client: client,
	}
}

func (*PolicyManager) Create(ctx context.Context, policy ladon.Policy) error {
	return nil
}

func (*PolicyManager) Update(ctx context.Context, policy ladon.Policy) error {
	return nil
}

func (*PolicyManager) Get(ctx context.Context, id string) (ladon.Policy, error) {
	return nil, nil
}

func (*PolicyManager) Delete(ctx context.Context, id string) error {
	return nil
}

func (*PolicyManager) GetAll(ctx context.Context, limit, offset int64) (ladon.Policies, error) {
	return nil, nil
}

func (m *PolicyManager) FindRequestCandidates(ctx context.Context, r *ladon.Request) (ladon.Policies, error) {
	username := ""

	if user, ok := r.Context["username"].(string); ok {
		username = user
	}

	policies, err := m.client.List(username)
	if err != nil {
		return nil, errors.Wrap(err, "list policies failed")
	}

	ret := make([]ladon.Policy, 0, len(policies))
	for _, policy := range policies {
		ret = append(ret, policy)
	}

	return ret, nil
}

func (m *PolicyManager) FindPoliciesForSubject(ctx context.Context, subject string) (ladon.Policies, error) {
	return nil, nil
}

func (m *PolicyManager) FindPoliciesForResource(ctx context.Context, resource string) (ladon.Policies, error) {
	return nil, nil
}
