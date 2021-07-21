package entity

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/8treenet/cdp-service/domain/po"
	"github.com/8treenet/cdp-service/infra/cattle"
	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
)

type Report struct {
	po.AnalysisReport
	freedom.Entity
}

func (entity *Report) Identity() string {
	return fmt.Sprint(entity.ID)
}

func (entity *Report) RatioSingle(numerator *Report) (interface{}, error) {
	selfMap := map[string]interface{}{}
	numeratorMap := map[string]interface{}{}
	if e := json.Unmarshal(entity.Data, &selfMap); e != nil {
		return nil, e
	}
	if e := json.Unmarshal(numerator.Data, &numeratorMap); e != nil {
		return nil, e
	}
	if len(selfMap) != 1 {
		return nil, fmt.Errorf("未知的分母数据错误")
	}
	if len(numeratorMap) != 1 {
		return nil, fmt.Errorf("未知的分子数据错误")
	}

	var selfValue float64
	var numeratorValue float64

	for _, v := range selfMap {
		parseValue, err := strconv.ParseFloat(fmt.Sprint(v), 64)
		if err != nil {
			return 0.0, err
		}
		selfValue = parseValue
		break
	}

	for _, v := range numeratorMap {
		parseValue, err := strconv.ParseFloat(fmt.Sprint(v), 64)
		if err != nil {
			return 0.0, err
		}
		numeratorValue = parseValue
		break
	}

	result := map[string]interface{}{
		cattle.OutTypeRatioValue: numeratorValue / selfValue,
	}
	return result, nil
}

func (entity *Report) RatioMultiple(numerator *Report) (interface{}, error) {
	selfMap := map[string]interface{}{}
	numeratorMaps := []map[string]interface{}{}
	if e := json.Unmarshal(entity.Data, &selfMap); e != nil {
		return nil, e
	}
	if e := json.Unmarshal(numerator.Data, &numeratorMaps); e != nil {
		return nil, e
	}
	if len(selfMap) != 1 {
		return nil, fmt.Errorf("未知的分母数据错误")
	}

	var selfValue float64
	for _, v := range selfMap {
		parseValue, err := strconv.ParseFloat(fmt.Sprint(v), 64)
		if err != nil {
			return 0.0, err
		}
		selfValue = parseValue
		break
	}

	result := []map[string]interface{}{}
	for _, numeratorMap := range numeratorMaps {
		newMap := map[string]interface{}{}
		for key, v := range numeratorMap {
			if !utils.InSlice([]string{
				cattle.OutTypeAvgValue,
				cattle.OutTypeCountValue,
				cattle.OutTypeMaxValue,
				cattle.OutTypeMinValue,
				cattle.OutTypePeopleValue,
				cattle.OutTypeSumValue,
			}, key) {
				newMap[key] = v
				continue
			}

			parseValue, err := strconv.ParseFloat(fmt.Sprint(v), 64)
			if err != nil {
				return 0.0, err
			}
			newMap[cattle.OutTypeRatioValue] = parseValue / selfValue
		}
		result = append(result, newMap)
	}

	return result, nil
}
