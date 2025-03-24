package arangodbaas

import (
	"context"

	"github.com/netcracker/qubership-core-lib-go/v3/logging"
	"github.com/netcracker/qubership-core-lib-go-dbaas-arangodb-client/v4/model"
	dbaasbase "github.com/netcracker/qubership-core-lib-go-dbaas-base-client/v3"
	"github.com/netcracker/qubership-core-lib-go-dbaas-base-client/v3/cache"
)

var logger logging.Logger

func init() {
	logger = logging.GetLogger("arangodbaas")
}

const (
	propMicroserviceName = "microservice.name"

	dbIdParam = "dbaas.arangodb.dbId"

	DbType = "arangodb"

	DbIdKey = "dbId"

	DefaultDbId = "default"

	TlsArangodbPort = 8530
)

type DbaaSArangoDbClient struct {
	arangoClientCache cache.DbaaSCache
	pool              *dbaasbase.DbaaSPool
}

func NewClient(pool *dbaasbase.DbaaSPool) *DbaaSArangoDbClient {
	localCache := cache.DbaaSCache{
		LogicalDbCache: make(map[cache.Key]interface{}),
	}
	return &DbaaSArangoDbClient{
		arangoClientCache: localCache,
		pool:              pool,
	}
}

func (d *DbaaSArangoDbClient) ServiceDatabase(params ...model.DbParams) Database {
	return &arangoDatabase{
		params:      d.buildServiceDbParams(params),
		dbaasPool:   d.pool,
		arangoCache: &d.arangoClientCache,
	}
}

func (d *DbaaSArangoDbClient) TenantDatabase(params ...model.DbParams) Database {
	return &arangoDatabase{
		params:      d.buildTenantDbParams(params),
		dbaasPool:   d.pool,
		arangoCache: &d.arangoClientCache,
	}
}

func ServiceClassifier(ctx context.Context) map[string]interface{} {
	classifier := dbaasbase.BaseServiceClassifier(ctx)
	if dbId, ok := ctx.Value(dbIdParam).(string); ok {
		classifier[DbIdKey] = dbId
	} else {
		classifier[DbIdKey] = DefaultDbId
	}
	return classifier
}

func TenantClassifier(ctx context.Context) map[string]interface{} {
	classifier := dbaasbase.BaseTenantClassifier(ctx)
	if dbId, ok := ctx.Value(dbIdParam).(string); ok {
		classifier[DbIdKey] = dbId
	} else {
		classifier[DbIdKey] = DefaultDbId
	}
	return classifier
}

func (d *DbaaSArangoDbClient) buildServiceDbParams(params []model.DbParams) model.DbParams {
	localParams := model.DbParams{}
	if params != nil {
		localParams = params[0]
	}
	if localParams.Classifier == nil {
		localParams.Classifier = ServiceClassifier
	}
	return localParams
}

func (d *DbaaSArangoDbClient) buildTenantDbParams(params []model.DbParams) model.DbParams {
	localParams := model.DbParams{}
	if params != nil {
		localParams = params[0]
	}
	if localParams.Classifier == nil {
		localParams.Classifier = TenantClassifier
	}
	return localParams
}
