package cattle

import (
	"github.com/8treenet/freedom"
	"github.com/go-xorm/builder"
)

type Query struct {
	logger  Logger
	manager *Manager
	builder *builder.Builder
}

func (q *Query) init() {
	q.logger = freedom.Logger()
}

func (q *Query) SetLogger(l Logger) *Query {
	q.logger = l
	return q
}

func (q *Query) Do(dest interface{}) error {
	sql, err := q.builder.ToBoundSQL()
	if err != nil {
		return err
	}
	return q.manager.db.Raw(sql).Scan(dest).Error
}
