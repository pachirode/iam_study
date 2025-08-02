package code

// go:generate codegen -type=int

const (
	ErrUserNotFound int = iota + 110001
	ErrUserAlreadyExist
)

const (
	ErrReachMaxCount int = iota + 110101
	ErrSecretNotFound
)

const (
	ErrPolicyNotFound int = iota + 110201
)
