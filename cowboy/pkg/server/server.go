package server

import (
	"context"
	"fmt"
	"strconv"

	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/dao"
	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/k8s"
	"github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/pkg/model"
	pb "github.com/evalsocket/envoy-kratos-grpc-auth/cowboy/proto"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/metadata"

	"sync"
)

// Backend implements the protobuf interface
type Backend struct {
	mu *sync.RWMutex
	k8 *k8s.K8sClient
}

// New initializes a new Backend struct.
func New() *Backend {

	k8 := k8s.NewK8s()
	return &Backend{
		mu: &sync.RWMutex{},
		k8: k8,
	}
}

// CreateApplication get application
func (b *Backend) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.GeneralResponse, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if md, ok := metadata.FromIncomingContext(ctx); ok || !ok {
		clientId := md["x-current-user"]
		fmt.Println(clientId)
		// TODO:
		// Create Resource and service
		// Deploy them in specific ns
		// create Entry in db
		tp, err := strconv.Atoi(req.Replica)
		if err != nil {
			fmt.Println("error in converting port")
		}
		rep := int32(tp)
		var appconfig = model.AppConfig{
			Image:   req.Image,
			Name:    req.Name,
			Ports:   req.Ports,
			Replica: &rep,
			Org:     clientId[0],
			Region:  req.Region,
		}
		fmt.Println(appconfig)
		daoResource, err := dao.NewDao()
		if err != nil {
			fmt.Println(err)
		}
		deployment, service, err := b.k8.CreateApplication(appconfig)
		if err != nil {
			fmt.Println(err)
			return &pb.GeneralResponse{
				Succeed: false,
				Message: "Failed to create application",
			}, nil
		}
		dbmodel := &model.DBModel{
			Deployment: *deployment,
			Service:    *service,
			Org:        appconfig.Org,
			ID:         uuid.Must(uuid.NewV4()).String(),
		}
		fmt.Println(dbmodel)
		if err := daoResource.CreateApplication(dbmodel); err != nil {
			return &pb.GeneralResponse{
				Succeed: false,
				Message: "Failed to create deployment",
			}, nil
		}
		return &pb.GeneralResponse{
			Succeed: true,
			Message: uuid.Must(uuid.NewV4()).String(),
		}, nil
	}
	return &pb.GeneralResponse{
		Succeed: true,
		Message: uuid.Must(uuid.NewV4()).String(),
	}, nil
}
