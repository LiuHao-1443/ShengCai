package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
	v1 "shengcai/api/v1"
	"shengcai/internal/service"
)

type ShengCaiHandler struct {
	*Handler
	shengCaiService service.ShengCaiService
}

func NewShengCaiHandler(handler *Handler, shengCaiService service.ShengCaiService) *ShengCaiHandler {
	return &ShengCaiHandler{
		Handler:         handler,
		shengCaiService: shengCaiService,
	}
}

func (h *ShengCaiHandler) List(ctx *gin.Context) {
	req := new(v1.ShengCaiListRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.WithContext(ctx).Error("ShengCaiHandler.List!!! ctx.ShouldBindJSON error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.logger.WithContext(ctx).Error("ShengCaiHandler.List!!! validate.Struct error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if cellDataList, err := h.shengCaiService.List(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("ShengCaiHandler.List!!! shengCaiService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	} else {
		v1.HandleSuccess(ctx, cellDataList)
		return
	}
}

func (h *ShengCaiHandler) GetMetaData(ctx *gin.Context) {
	req := new(v1.ShengCaiGetMetaDataRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		h.logger.WithContext(ctx).Error("ShengCaiHandler.GetMetaData!!! ctx.ShouldBindJSON error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.logger.WithContext(ctx).Error("ShengCaiHandler.GetMetaData!!! validate.Struct error", zap.Error(err))
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if metaData, err := h.shengCaiService.GetMetaData(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("ShengCaiHandler.GetMetaData!!! shengCaiService.GetMetaData error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	} else {
		v1.HandleSuccess(ctx, metaData)
		return
	}
}

func (h *ShengCaiHandler) CreateData(ctx *gin.Context) {
	if err := h.shengCaiService.CreateData(ctx); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	} else {
		v1.HandleSuccess(ctx, nil)
		return
	}
}
