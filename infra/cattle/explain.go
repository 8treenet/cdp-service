package cattle

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/8treenet/cdp-service/utils"
	"github.com/go-xorm/builder"
)

func ExplainSingleAnalysis(dsl *DSL, beginTime, endTime time.Time) (*builder.Builder, error) {
	from := dsl.FindFromNode()
	if from == nil {
		return nil, errors.New("未找到from标签")
	}
	where := builder.And(builder.Between{Col: from.GetContent() + "." + ColumnCreateTime, LessVal: utils.DateTimeFormat(beginTime), MoreVal: utils.DateTimeFormat(endTime)})
	if dsl.FindConditionNode() != nil {
		cond, err := dsl.Condition(dsl.FindConditionNode())
		if err != nil {
			return nil, err
		}
		where = where.And(cond)
	}

	var isPeople bool
	selectBuilder, err := dsl.SingleOut(&isPeople)
	if err != nil {
		return nil, err
	}

	if isPeople {
		tc := dsl.FindFromNode().GetContent() + "." + ColumnUserId
		subSelect := builder.Select(tc).GroupBy(tc).Where(where)
		subSelect = dsl.From(subSelect)
		join, err := dsl.Join(subSelect)
		if err != nil {
			panic(err)
		}
		selectBuilder = selectBuilder.From(join, "people")
		return selectBuilder, nil
	}

	join, err := dsl.Join(dsl.From(selectBuilder))
	if err != nil {
		return nil, err
	}
	selectBuilder = join.Where(where)
	return selectBuilder, nil
}

func ExplainMultipleAnalysis(dsl *DSL, beginTime, endTime time.Time) (*builder.Builder, error) {
	from := dsl.FindFromNode()
	if from == nil {
		return nil, errors.New("未找到from标签")
	}

	where := builder.And(builder.Between{Col: from.GetContent() + "." + ColumnCreateTime, LessVal: utils.DateTimeFormat(beginTime), MoreVal: utils.DateTimeFormat(endTime)})
	if dsl.FindConditionNode() != nil {
		cond, err := dsl.Condition(dsl.FindConditionNode())
		if err != nil {
			return nil, err
		}
		where = where.And(cond)
	}

	selectBuilder, err := dsl.MultipleOut()
	if err != nil {
		return nil, err
	}
	join, err := dsl.Join(dsl.From(selectBuilder))
	if err != nil {
		return nil, err
	}
	return join.Where(where), nil
}

func ExplainPersonasAnalysis(dsl *DSL, userIds []string, beginTimes ...time.Time) (*builder.Builder, error) {
	from := dsl.FindFromNode()
	if from == nil {
		return nil, errors.New("未找到from标签")
	}
	list, err := utils.ToInterfaces(userIds)
	if err != nil {
		return nil, err
	}

	where := builder.And(builder.In(fmt.Sprintf("%s.%s", from.GetContent(), ColumnUserId), list...))
	if dsl.FindConditionNode() != nil {
		cond, err := dsl.Condition(dsl.FindConditionNode())
		if err != nil {
			return nil, err
		}
		where = where.And(cond)
	}

	selectBuilder, err := dsl.PersonasOut()
	if err != nil {
		return nil, err
	}
	join, err := dsl.Join(dsl.From(selectBuilder))
	if err != nil {
		return nil, err
	}

	beginTime, err := DayToTime(dsl.FindPersonasNode().GetAttribute(attributeDay))
	if err != nil {
		err = fmt.Errorf("PersonasNode.GetAttribute(day) error: %w", err)
		return nil, err
	}
	if len(beginTimes) > 0 {
		beginTime = &beginTimes[0]
	}
	where = where.And(builder.Gte{from.GetContent() + "." + ColumnCreateTime: utils.DateTimeFormat(*beginTime)})
	return join.Where(where), nil
}

func DayToTime(day string) (t *time.Time, e error) {
	if day == "0.5" {
		now := time.Now().Add(-(12 * time.Hour))
		t = &now
		return
	}

	i, e := strconv.Atoi(day)
	if e != nil {
		return
	}
	now := time.Now().AddDate(0, 0, -i)
	t = &now
	return
}
