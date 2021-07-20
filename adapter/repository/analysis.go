package repository

import (
	"encoding/base64"

	"github.com/8treenet/cdp-service/domain/entity"
	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra"
	"github.com/8treenet/cdp-service/infra/cattle"
	"github.com/8treenet/freedom"
	"gorm.io/gorm"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *AnalysisRepository {
			return &AnalysisRepository{detailedCacheKey: "cdp_detailed_%d"}
		})
	})
}

// AnalysisRepository 统计仓库.
type AnalysisRepository struct {
	freedom.Repository
	Common           *infra.CommonRequest
	Manager          *cattle.Manager
	detailedCacheKey string
}

func (repo *AnalysisRepository) FindByName(name string) (result *entity.Analysis, e error) {
	defer func() {
		if e != nil {
			return
		}
		result.XMLBytes, e = base64.StdEncoding.DecodeString(result.XMLData)
	}()
	if name == "" {
		e = gorm.ErrRecordNotFound
		return
	}

	result = &entity.Analysis{}
	result.Name = name
	repo.InjectBaseEntity(result)
	e = findAnalysis(repo, &result.Analysis)
	return
}

func (repo *AnalysisRepository) NewAnalysisEntity() (result *entity.Analysis) {
	result = &entity.Analysis{}
	repo.InjectBaseEntity(result)
	return
}

func (repo *AnalysisRepository) SaveAnalysisEntity(entity *entity.Analysis) error {
	entity.XMLData = base64.StdEncoding.EncodeToString(entity.XMLBytes)
	if entity.ID == 0 {
		_, e := createAnalysis(repo, &entity.Analysis)
		return e
	}

	_, e := saveAnalysis(repo, entity)
	return e
}

// GetAnalysisByPage .
func (repo *AnalysisRepository) GetAnalysisByPage() (result []*entity.Analysis, totalPage int, e error) {
	result = make([]*entity.Analysis, 0)

	page, pageSize := repo.Common.GetPage()
	pager := NewDescPager("id").SetPage(page, pageSize)
	list, e := findAnalysisList(repo, po.Analysis{}, pager)
	if e != nil {
		return
	}

	for i := 0; i < len(list); i++ {
		entity := &entity.Analysis{Analysis: list[0]}
		entity.XMLBytes, e = base64.StdEncoding.DecodeString(entity.XMLData)
		if e != nil {
			return
		}

		result = append(result, entity)
	}
	totalPage = pager.TotalPage()
	repo.InjectBaseEntitys(result)
	return
}

// GetAllAnalysis
func (repo *AnalysisRepository) GetAllAnalysis() (result []*entity.Analysis, e error) {
	result = make([]*entity.Analysis, 0)
	list, e := findAnalysisList(repo, po.Analysis{})
	if e != nil {
		return
	}

	for i := 0; i < len(list); i++ {
		entity := &entity.Analysis{Analysis: list[0]}
		entity.XMLBytes, e = base64.StdEncoding.DecodeString(entity.XMLData)
		if e != nil {
			return
		}

		result = append(result, entity)
	}
	repo.InjectBaseEntitys(result)
	return
}

// db .
func (repo *AnalysisRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
