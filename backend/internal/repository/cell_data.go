package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	v1 "shengcai/api/v1"
	"shengcai/internal/model"
)

type CellDataRepository interface {
	Create(ctx context.Context, cellData *model.CellData) error
	List(ctx context.Context, filter struct {
		SheetID string `json:"sheet_id" validate:"required"`
	}, page int) (*v1.ShengCaiListResponse, error)
	GetMetaData(ctx context.Context, sheetId string) (*v1.ShengCaiGetMetaDataResponse, error)
}

func NewCellDataRepository(
	r *Repository,
	aiRepository AIRepository,
) CellDataRepository {
	return &cellDataRepository{
		Repository:   r,
		aiRepository: aiRepository,
	}
}

type cellDataRepository struct {
	*Repository
	aiRepository AIRepository
}

func (r *cellDataRepository) Create(ctx context.Context, cellData *model.CellData) error {
	var existingCellData model.CellData

	// 检查数据库中是否已存在相同的 Link
	err := r.DB(ctx).Where("sheet_id =?", cellData.SheetID).Where("link = ?", cellData.Link).First(&existingCellData).Error
	fmt.Println("===================================")
	fmt.Println(err)
	fmt.Println("===================================")

	if err == nil {
		// 如果找到了相同的 Link，比较 ReleaseDate
		if existingCellData.ReleaseDate == cellData.ReleaseDate && existingCellData.SortNumber == cellData.SortNumber {
			// 如果 ReleaseDate 相同，则不进行任何操作
			return nil
		} else {
			// 如果 ReleaseDate 不同，则更新 Content 和 ReleaseDate
			existingCellData.Content = cellData.Content
			existingCellData.ReleaseDate = cellData.ReleaseDate
			existingCellData.SortNumber = cellData.SortNumber

			abstract, keyword, err := r.aiRepository.GenerateAbstractAndKeyword(ctx, existingCellData.Content)
			if err != nil {
				fmt.Println("+++++++++++++++++++++++++++++++++++")
				fmt.Println(err)
				fmt.Println("+++++++++++++++++++++++++++++++++++")
			}

			existingCellData.Abstract = abstract
			existingCellData.Keyword = keyword
			if err := r.DB(ctx).Save(&existingCellData).Error; err != nil {
				fmt.Println("-----------------------------------")
				fmt.Println(err)
				fmt.Println("-----------------------------------")
				return err
			}
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		abstract, keyword, err := r.aiRepository.GenerateAbstractAndKeyword(ctx, cellData.Content)
		if err != nil {
			fmt.Println("+++++++++++++++++++++++++++++++++++")
			fmt.Println(err)
			fmt.Println("+++++++++++++++++++++++++++++++++++")
		}
		cellData.Abstract = abstract
		cellData.Keyword = keyword

		// 如果没有找到相同的 Link，则进行新建
		if err := r.DB(ctx).Create(cellData).Error; err != nil {
			fmt.Println("***********************************")
			fmt.Println(err)
			fmt.Println("***********************************")
			return err
		}
	} else {
		// 如果查询出错且错误不是 "记录未找到"，则返回查询错误
		return err
	}

	return nil
}

func (r *cellDataRepository) List(ctx context.Context, filter struct {
	SheetID string `json:"sheet_id" validate:"required"`
}, page int) (*v1.ShengCaiListResponse, error) {
	// 构建查询条件
	query := r.DB(ctx)
	query = query.Where("sheet_id = ?", filter.SheetID)

	// 获取总记录数
	var totalCount int64
	if err := query.Model(&model.CellData{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 执行分页查询
	var list []*model.CellData
	err := query.Order("sort_number ASC").Offset((page - 1) * 10).Limit(10).Find(&list).Error
	if err != nil {
		return nil, err
	}

	// 构建返回结果
	result := &v1.ShengCaiListResponse{
		TotalCount: int(totalCount),
		List: make([]struct {
			SheetID     string `json:"sheet_id"`
			Title       string `json:"title"`
			Link        string `json:"link"`
			ReleaseDate string `json:"release_date"`
			Abstract    string `json:"abstract"`
			Keyword     string `json:"keyword"`
		}, len(list)),
	}

	for i, item := range list {
		result.List[i] = struct {
			SheetID     string `json:"sheet_id"`
			Title       string `json:"title"`
			Link        string `json:"link"`
			ReleaseDate string `json:"release_date"`
			Abstract    string `json:"abstract"`
			Keyword     string `json:"keyword"`
		}{
			SheetID:     item.SheetID,
			Title:       item.Title,
			Link:        item.Link,
			ReleaseDate: item.ReleaseDate,
			Abstract:    item.Abstract,
			Keyword:     item.Keyword,
		}
	}

	return result, nil
}

func (r *cellDataRepository) GetMetaData(ctx context.Context, sheetId string) (*v1.ShengCaiGetMetaDataResponse, error) {
	var result v1.ShengCaiGetMetaDataResponse
	var sheetInfo struct {
		SheetName string `gorm:"column:sheet_name"`
		UpdateLog string `gorm:"column:update_log"`
	}

	// 查询 sheet_info 表
	if err := r.db.WithContext(ctx).Table("sheet_info").Where("sheet_id = ?", sheetId).Select("sheet_name, update_log").First(&sheetInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.SheetName = ""
			result.UpdateLog = ""
		}
		// 处理其他错误
		return nil, err
	} else {
		// 将查询结果赋值给响应结构体
		result.SheetName = sheetInfo.SheetName
		result.UpdateLog = sheetInfo.UpdateLog
	}

	return &result, nil
}
