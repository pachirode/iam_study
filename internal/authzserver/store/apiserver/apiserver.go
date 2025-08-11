package apiserver

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/pachirode/iam_study/internal/authzserver/store"
	pb "github.com/pachirode/iam_study/internal/pkg/api/proto/apiserver/v1"
	"github.com/pachirode/iam_study/pkg/log"
)

type dataStore struct {
	cli pb.CacheClient
}

func (ds *dataStore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *dataStore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

var (
	apiServerFactory store.Factory
	once             sync.Once
)

func GetAPIServerFactoryOrDie(address string, clientCA string) store.Factory {
	once.Do(func() {
		var (
			err   error
			conn  *grpc.ClientConn
			creds credentials.TransportCredentials
		)

		creds, err = credentials.NewClientTLSFromFile(clientCA, "")
		if err != nil {
			log.Panicf("Credentials.NewClientTLSFromFile err: ", err)
		}

		conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Panicf("Connect to grpc server failed, error: %s", err.Error())
		}

		apiServerFactory = &dataStore{pb.NewCacheClient(conn)}
		log.Infof("Connected to grpc server, address: %s", address)
	})

	if apiServerFactory == nil {
		log.Panicf("Failed to get apiserver store factory")
	}

	return apiServerFactory
}
