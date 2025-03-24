package arangodbaas

import (
	"context"

	"github.com/arangodb/go-driver/v2/connection"
	"github.com/netcracker/qubership-core-lib-go-dbaas-arangodb-client/v4/model"
	dbaasbase "github.com/netcracker/qubership-core-lib-go-dbaas-base-client/v3"
	"github.com/netcracker/qubership-core-lib-go-dbaas-base-client/v3/cache"
)

type Database interface {
	GetArangoDbClient(httpConfigs ...*connection.HttpConfiguration) (ArangoDbClient, error)
	GetConnectionProperties(ctx context.Context) (*model.ArangoConnProperties, error)
	FindConnectionProperties(ctx context.Context) (*model.ArangoConnProperties, error)

	GetConnectionPropertiesWithDbId(ctx context.Context, dbId string) (*model.ArangoConnProperties, error)
	FindConnectionPropertiesWithDbId(ctx context.Context, dbId string) (*model.ArangoConnProperties, error)
}

type arangoDatabase struct {
	dbaasPool   *dbaasbase.DbaaSPool
	params      model.DbParams
	arangoCache *cache.DbaaSCache
}

func (d arangoDatabase) GetArangoDbClient(httpConfigs ...*connection.HttpConfiguration) (ArangoDbClient, error) {
	httpConfig := &connection.HttpConfiguration{}
	if httpConfigs != nil {
		httpConfig = httpConfigs[0]
	}
	return &ArangoDbClientImpl{
		httpConfig:    httpConfig,
		dbaasClient:   d.dbaasPool.Client,
		arangodbCache: d.arangoCache,
		params:        d.params,
	}, nil
}

func (d arangoDatabase) GetConnectionPropertiesWithDbId(ctx context.Context, dbId string) (*model.ArangoConnProperties, error) {
	return d.GetConnectionProperties(context.WithValue(ctx, dbIdParam, dbId))
}

func (d arangoDatabase) FindConnectionPropertiesWithDbId(ctx context.Context, dbId string) (*model.ArangoConnProperties, error) {
	return d.FindConnectionProperties(context.WithValue(ctx, dbIdParam, dbId))
}

func (d arangoDatabase) GetConnectionProperties(ctx context.Context) (*model.ArangoConnProperties, error) {
	baseDbParams := d.params.BaseDbParams
	classifier := d.params.Classifier(ctx)

	agrLogicalDb, err := d.dbaasPool.GetOrCreateDb(ctx, DbType, classifier, baseDbParams)
	if err != nil {
		logger.Error("Error acquiring connection properties from DBaaS: %v", err)
		return nil, err
	}
	agrConnProperties := toArangoConnProperties(agrLogicalDb.ConnectionProperties)
	return &agrConnProperties, nil
}

func (d arangoDatabase) FindConnectionProperties(ctx context.Context) (*model.ArangoConnProperties, error) {
	classifier := d.params.Classifier(ctx)
	params := d.params.BaseDbParams
	responseBody, err := d.dbaasPool.GetConnection(ctx, DbType, classifier, params)
	if err != nil {
		logger.ErrorC(ctx, "Error finding connection properties from DBaaS: %v", err)
		return nil, err
	}
	logger.Info("Found connection to arango db with classifier %+v", classifier)
	agrConnProperties := toArangoConnProperties(responseBody)
	return &agrConnProperties, err
}

func toArangoConnProperties(connProperties map[string]interface{}) model.ArangoConnProperties {
	return model.ArangoConnProperties{
		DbName:   connProperties["dbName"].(string),
		Host:     connProperties["host"].(string),
		Port:     int(connProperties["port"].(float64)),
		Username: connProperties["username"].(string),
		Password: connProperties["password"].(string),
	}
}
