package store

var client Factory

type Factory interface {
	Policies() PolicyStore
	Secrets() SecretStore
}

func Client() Factory {
	return client
}

func SetClient(factory Factory) {
	client = factory
}
