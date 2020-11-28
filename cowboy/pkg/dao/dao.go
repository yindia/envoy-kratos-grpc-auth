package dao

import (
	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/dao/es"
	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/db"
	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/model"
)

// ICowboyDao : interface to dao operation
type ICowboyDao interface {
	CreateApplication(app *model.DBModel) error
}

// ToDo : can make it a factory if a more dynamic logic needed. eg., in case of multiple datastores.
// NewDao : returns dao impl
func NewDao() (ICowboyDao, error) {
	return newESDao()
}

// newESDao : returns es impl of contact dao
func newESDao() (ICowboyDao, error) {
	esclient, err := db.NewESClient()
	if err != nil {
		return nil, err
	}
	return es.NewResourceDao(esclient), nil
}
