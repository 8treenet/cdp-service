package cattle

import "github.com/go-xorm/builder"

func ExplainSingleAnalysis(dsl *DSL) (*builder.Builder, error) {
	var where builder.Cond
	if dsl.FindCondition() != nil {
		cond, err := dsl.Condition(dsl.FindCondition())
		if err != nil {
			return nil, err
		}
		where = cond
	}

	var isPeople bool
	selectBuilder, err := dsl.SingleOut(&isPeople)
	if err != nil {
		return nil, err
	}

	if isPeople {
		tc := dsl.FindFrom().GetContent() + "." + ColumnUserId
		subSelect := builder.Select(tc).GroupBy(tc)
		if where != nil {
			subSelect = subSelect.Where(where)
		}
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

func ExplainMultipleAnalysis(dsl *DSL) (*builder.Builder, error) {
	var where builder.Cond
	if dsl.FindCondition() != nil {
		cond, err := dsl.Condition(dsl.FindCondition())
		if err != nil {
			return nil, err
		}
		where = cond
	}

	selectBuilder, err := dsl.MultipleOut()
	if err != nil {
		return nil, err
	}
	join, err := dsl.Join(dsl.From(selectBuilder))
	if err != nil {
		return nil, err
	}
	selectBuilder = join
	if where != nil {
		selectBuilder.Where(where)
	}
	return selectBuilder, nil
}
