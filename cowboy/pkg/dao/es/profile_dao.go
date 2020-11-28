package es

import (
	"context"

	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/model"

	"fmt"

	elastic "github.com/olivere/elastic/v7"
)

const typeNum = "resources"

// ResourceDao : holds the elasticsearch client
type ResourceDao struct {
	client *elastic.Client
}

// NewResourceDao : returns instance of Resource Dao
func NewResourceDao(client *elastic.Client) *ResourceDao {
	return &ResourceDao{
		client: client,
	}
}

func (e *ResourceDao) CreateApplication(app *model.DBModel) error {
	index := "test"
	typ := typeNum
	fmt.Println(e.client)

	resource, err := e.client.Index().Index(index).Type(typ).Id(app.ID).BodyJson(app).Do(context.Background())
	if err != nil {
		if eError, ok := err.(*elastic.Error); ok == true {
			fmt.Println("Index not found : ")
			fmt.Println(eError.Details)
			return err
		}
	}
	fmt.Println(resource)
	return nil
}
