package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"shengcai/internal/model"
)

type SheetInfoRepository interface {
	Create(ctx context.Context, sheetInfo *model.SheetInfo) error
}

func NewSheetInfoRepository(
	r *Repository,
) SheetInfoRepository {
	return &sheetInfoRepository{
		Repository: r,
	}
}

type sheetInfoRepository struct {
	*Repository
}

func (r *sheetInfoRepository) Create(ctx context.Context, sheetInfo *model.SheetInfo) error {
	// 检查数据库中是否已存在相同的 sheet_id
	var existingSheetInfo model.SheetInfo
	err := r.DB(ctx).Where("sheet_id = ?", sheetInfo.SheetID).First(&existingSheetInfo).Error

	if err == nil {
		// 如果找到了相同的 sheet_id，则返回自定义错误或跳过创建
		fmt.Printf("sheet with sheet_id '%s' already exists\n", sheetInfo.SheetID)
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果查询出错且错误不是 "记录未找到"，则返回查询错误
		return err
	}

	// 如果没有找到相同的 sheet_id，则进行创建
	if err = r.DB(ctx).Create(sheetInfo).Error; err != nil {
		return err
	}
	return nil
}
