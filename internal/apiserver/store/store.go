package store

var client Factory

type Factory interface {
	Users() UserStore
	Secrets() SecretStore
	Policies() PolicyStore
	PolicyAudits() PolicyAuditStore
	Close() error
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
