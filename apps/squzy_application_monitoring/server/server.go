package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	apiPb "github.com/squzy/squzy_generated/generated/proto/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"squzy/apps/squzy_application_monitoring/config"
	"squzy/apps/squzy_application_monitoring/database"
	"squzy/internal/helpers"
)

type server struct {
	db      database.Database
	config  config.Config
	storage apiPb.StorageClient
}

func (s *server) updateStatus(ctx context.Context, applicationId primitive.ObjectID, status apiPb.ApplicationStatus) (*apiPb.Application, error) {
	err := s.db.SetStatus(ctx, applicationId, status)

	if err != nil {
		return nil, err
	}

	app, err := s.db.FindApplicationById(ctx, applicationId)

	if err != nil {
		return nil, err
	}

	return transformDbApplication(app), nil
}

func (s *server) ArchiveApplicationById(ctx context.Context, request *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	applicationId, err := primitive.ObjectIDFromHex(request.ApplicationId)
	if err != nil {
		return nil, err
	}
	return s.updateStatus(ctx, applicationId, apiPb.ApplicationStatus_APPLICATION_STATUS_ARCHIVED)
}

func (s *server) EnableApplicationById(ctx context.Context, request *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	applicationId, err := primitive.ObjectIDFromHex(request.ApplicationId)
	if err != nil {
		return nil, err
	}
	return s.updateStatus(ctx, applicationId, apiPb.ApplicationStatus_APPLICATION_STATUS_ENABLED)
}

func (s *server) DisableApplicationById(ctx context.Context, request *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	applicationId, err := primitive.ObjectIDFromHex(request.ApplicationId)
	if err != nil {
		return nil, err
	}
	return s.updateStatus(ctx, applicationId, apiPb.ApplicationStatus_APPLICATION_STATUS_DISABLED)
}

func transformDbApplication(dbApp *database.Application) *apiPb.Application {
	return &apiPb.Application{
		Id:       dbApp.Id.Hex(),
		Name:     dbApp.Name,
		HostName: dbApp.Host,
		Status:   dbApp.Status,
	}
}

func (s *server) GetApplicationById(ctx context.Context, request *apiPb.ApplicationByIdReuqest) (*apiPb.Application, error) {
	applicationId, err := primitive.ObjectIDFromHex(request.ApplicationId)
	if err != nil {
		return nil, err
	}

	app, err := s.db.FindApplicationById(ctx, applicationId)
	if err != nil {
		return nil, err
	}

	return transformDbApplication(app), nil
}

func (s *server) GetApplicationList(ctx context.Context, e *empty.Empty) (*apiPb.GetApplicationListResponse, error) {
	list, err := s.db.FindAllApplication(ctx)
	if err != nil {
		return nil, err
	}

	appList := []*apiPb.Application{}

	for _, v := range list {
		appList = append(appList, transformDbApplication(v))
	}

	return &apiPb.GetApplicationListResponse{
		Applications: appList,
	}, nil
}

var (
	errMissingName = errors.New("missing application name")
)

func (s *server) InitializeApplication(ctx context.Context, req *apiPb.ApplicationInfo) (*apiPb.InitializeApplicationResponse, error) {
	if req.Name == "" {
		return nil, errMissingName
	}

	application, err := s.db.FindOrCreate(ctx, req.Name, req.HostName)
	if err != nil {
		return nil, err
	}
	return &apiPb.InitializeApplicationResponse{
		ApplicationId: application.Id.Hex(),
		TracingHeader: s.config.GetTracingHeader(),
	}, nil
}

func (s *server) SaveTransaction(ctx context.Context, req *apiPb.TransactionInfo) (*empty.Empty, error) {
	applicationId, err := primitive.ObjectIDFromHex(req.ApplicationId)
	if err != nil {
		return nil, err
	}

	app, err := s.db.FindApplicationById(ctx, applicationId)
	if err != nil {
		return nil, err
	}

	if app.Status != apiPb.ApplicationStatus_APPLICATION_STATUS_ENABLED {
		return &empty.Empty{}, nil
	}

	go func() {
		reqCtx, cancel := helpers.TimeoutContext(context.Background(), s.config.GetStorageTimeout())
		defer cancel()
		_, _ = s.storage.SaveTransaction(reqCtx, req)
	}()

	return &empty.Empty{}, nil
}

func New(db database.Database, config config.Config, storage apiPb.StorageClient) apiPb.ApplicationMonitoringServer {
	return &server{
		db:      db,
		config:  config,
		storage: storage,
	}
}
