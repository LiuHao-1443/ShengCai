package service

import (
	"context"
	"os"
	v1 "shengcai/api/v1"
	"shengcai/internal/repository"
)

type ShengCaiService interface {
	List(ctx context.Context, req *v1.ShengCaiListRequest) (*v1.ShengCaiListResponse, error)
	GetMetaData(ctx context.Context, req *v1.ShengCaiGetMetaDataRequest) (*v1.ShengCaiGetMetaDataResponse, error)
	CreateData(ctx context.Context) error
}

func NewShengCaiService(
	service *Service,
	feiShuService FeiShuService,
	cellDataRepo repository.CellDataRepository,
) ShengCaiService {
	return &shengCaiService{
		Service:       service,
		FeiShuService: feiShuService,
		CellDataRepo:  cellDataRepo,
	}
}

type shengCaiService struct {
	*Service
	FeiShuService FeiShuService
	CellDataRepo  repository.CellDataRepository
}

func (s *shengCaiService) List(ctx context.Context, req *v1.ShengCaiListRequest) (*v1.ShengCaiListResponse, error) {
	if list, err := s.CellDataRepo.List(ctx, struct {
		SheetID string `json:"sheet_id" validate:"required"`
	}(struct{ SheetID string }{SheetID: req.SpreadsheetToken}), req.Page); err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (s *shengCaiService) GetMetaData(ctx context.Context, req *v1.ShengCaiGetMetaDataRequest) (*v1.ShengCaiGetMetaDataResponse, error) {
	if list, err := s.CellDataRepo.GetMetaData(ctx, req.SpreadsheetToken); err != nil {
		return nil, err
	} else {
		return list, nil
	}
}

func (s *shengCaiService) CreateData(ctx context.Context) error {
	appID := os.Getenv("app_id")
	if appID == "" {
		appID = s.conf.GetString("open.app_id")
	}

	appSecret := os.Getenv("app_secret")
	if appSecret == "" {
		appSecret = s.conf.GetString("open.app_secret")
	}

	apiKey := os.Getenv("api_key")
	if apiKey == "" {
		apiKey = s.conf.GetString("open.api_key")
	}

	spreadsheetToken := os.Getenv("spreadsheet_token")
	if spreadsheetToken == "" {
		spreadsheetToken = s.conf.GetString("feishu.spreadsheet_token")
	}

	return s.FeiShuService.SaveTableData(ctx, appID, appSecret, spreadsheetToken)
}
