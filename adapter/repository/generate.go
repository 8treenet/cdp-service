package repository

import (
	"errors"
	"fmt"
	"github.com/8treenet/crm-service/domain/po"
	"github.com/8treenet/freedom"
	"gorm.io/gorm"
	"strings"
	"time"
)

// GORMRepository .
type GORMRepository interface {
	db() *gorm.DB
	Worker() freedom.Worker
}

type saveObject interface {
	TableName() string
	Location() map[string]interface{}
	GetChanges() map[string]interface{}
}

// Builder .
type Builder interface {
	Execute(db *gorm.DB, object interface{}) error
}

// Pager .
type Pager struct {
	pageSize  int
	page      int
	totalPage int
	fields    []string
	orders    []string
}

// NewDescPager .
func NewDescPager(column string, columns ...string) *Pager {
	return newDefaultPager("desc", column, columns...)
}

// NewAscPager .
func NewAscPager(column string, columns ...string) *Pager {
	return newDefaultPager("asc", column, columns...)
}

// NewDescOrder .
func newDefaultPager(sort, field string, args ...string) *Pager {
	fields := []string{field}
	fields = append(fields, args...)
	orders := []string{}
	for index := 0; index < len(fields); index++ {
		orders = append(orders, sort)
	}
	return &Pager{
		fields: fields,
		orders: orders,
	}
}

// Order .
func (p *Pager) Order() interface{} {
	if len(p.fields) == 0 {
		return nil
	}
	args := []string{}
	for index := 0; index < len(p.fields); index++ {
		args = append(args, fmt.Sprintf("`%s` %s", p.fields[index], p.orders[index]))
	}

	return strings.Join(args, ",")
}

// TotalPage .
func (p *Pager) TotalPage() int {
	return p.totalPage
}

// SetPage .
func (p *Pager) SetPage(page, pageSize int) *Pager {
	p.page = page
	p.pageSize = pageSize
	return p
}

// Execute .
func (p *Pager) Execute(db *gorm.DB, object interface{}) (e error) {
	pageFind := false
	orderValue := p.Order()
	if orderValue != nil {
		db = db.Order(orderValue)
	}
	if p.page != 0 && p.pageSize != 0 {
		pageFind = true
		db = db.Offset((p.page - 1) * p.pageSize).Limit(p.pageSize)
	}

	resultDB := db.Find(object)
	if resultDB.Error != nil {
		return resultDB.Error
	}

	if !pageFind {
		return
	}

	var count64 int64
	e = resultDB.Offset(0).Limit(1).Count(&count64).Error
	count := int(count64)
	if e == nil && count != 0 {
		//Calculate the length of the pagination
		if count%p.pageSize == 0 {
			p.totalPage = count / p.pageSize
		} else {
			p.totalPage = count/p.pageSize + 1
		}
	}
	return
}

func ormErrorLog(repo GORMRepository, model, method string, e error, expression ...interface{}) {
	if e == nil || e == gorm.ErrRecordNotFound {
		return
	}
	repo.Worker().Logger().Errorf("error: %v, model: %s, method: %s", e, model, method)
}

// findCustomerTemplate .
func findCustomerTemplate(repo GORMRepository, result *po.CustomerTemplate, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplate", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplate", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerTemplateListByPrimarys .
func findCustomerTemplateListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplateListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplatesByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerTemplateByWhere .
func findCustomerTemplateByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplateByWhere", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplateByWhere", e, query, args)
	}()
	db := repo.db()
	if query != "" {
		db = db.Where(query, args...)
	}
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerTemplateByMap .
func findCustomerTemplateByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplateByMap", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplateByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerTemplateList .
func findCustomerTemplateList(repo GORMRepository, query po.CustomerTemplate, builders ...Builder) (results []po.CustomerTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplateList", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplates", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerTemplateListByWhere .
func findCustomerTemplateListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplateListByWhere", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplatesByWhere", e, query, args)
	}()
	db := repo.db()
	if query != "" {
		db = db.Where(query, args...)
	}

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerTemplateListByMap .
func findCustomerTemplateListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "findCustomerTemplateListByMap", e, now)
		ormErrorLog(repo, "CustomerTemplate", "findCustomerTemplatesByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerTemplate .
func createCustomerTemplate(repo GORMRepository, object *po.CustomerTemplate) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "createCustomerTemplate", e, now)
		ormErrorLog(repo, "CustomerTemplate", "createCustomerTemplate", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerTemplate .
func saveCustomerTemplate(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemplate", "saveCustomerTemplate", e, now)
		ormErrorLog(repo, "CustomerTemplate", "saveCustomerTemplate", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}
