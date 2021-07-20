package entity

import (
	"github.com/8treenet/freedom"
)

type Report struct {
	freedom.Entity
	AnalysisID int               `json:"analysisID"` //分析实体id
	Detail     map[string]string `json:"detail"`     //结果详情
}
