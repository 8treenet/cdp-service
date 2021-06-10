package repository

import (
	"errors"
	"fmt"
	"github.com/8treenet/cdp-service/domain/po"
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

// findCustomerExtendTemplate .
func findCustomerExtendTemplate(repo GORMRepository, result *po.CustomerExtendTemplate, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplate", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplate", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerExtendTemplateListByPrimarys .
func findCustomerExtendTemplateListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerExtendTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplateListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplatesByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerExtendTemplateByWhere .
func findCustomerExtendTemplateByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerExtendTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplateByWhere", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplateByWhere", e, query, args)
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

// findCustomerExtendTemplateByMap .
func findCustomerExtendTemplateByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerExtendTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplateByMap", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplateByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerExtendTemplateList .
func findCustomerExtendTemplateList(repo GORMRepository, query po.CustomerExtendTemplate, builders ...Builder) (results []po.CustomerExtendTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplateList", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplates", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerExtendTemplateListByWhere .
func findCustomerExtendTemplateListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerExtendTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplateListByWhere", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplatesByWhere", e, query, args)
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

// findCustomerExtendTemplateListByMap .
func findCustomerExtendTemplateListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerExtendTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "findCustomerExtendTemplateListByMap", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "findCustomerExtendTemplatesByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerExtendTemplate .
func createCustomerExtendTemplate(repo GORMRepository, object *po.CustomerExtendTemplate) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "createCustomerExtendTemplate", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "createCustomerExtendTemplate", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerExtendTemplate .
func saveCustomerExtendTemplate(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtendTemplate", "saveCustomerExtendTemplate", e, now)
		ormErrorLog(repo, "CustomerExtendTemplate", "saveCustomerExtendTemplate", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerExtend .
func findCustomerExtend(repo GORMRepository, result *po.CustomerExtend, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtend", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtend", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerExtendListByPrimarys .
func findCustomerExtendListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerExtend, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtendListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtendsByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerExtendByWhere .
func findCustomerExtendByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerExtend, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtendByWhere", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtendByWhere", e, query, args)
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

// findCustomerExtendByMap .
func findCustomerExtendByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerExtend, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtendByMap", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtendByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerExtendList .
func findCustomerExtendList(repo GORMRepository, query po.CustomerExtend, builders ...Builder) (results []po.CustomerExtend, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtendList", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtends", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerExtendListByWhere .
func findCustomerExtendListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerExtend, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtendListByWhere", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtendsByWhere", e, query, args)
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

// findCustomerExtendListByMap .
func findCustomerExtendListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerExtend, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "findCustomerExtendListByMap", e, now)
		ormErrorLog(repo, "CustomerExtend", "findCustomerExtendsByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerExtend .
func createCustomerExtend(repo GORMRepository, object *po.CustomerExtend) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "createCustomerExtend", e, now)
		ormErrorLog(repo, "CustomerExtend", "createCustomerExtend", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerExtend .
func saveCustomerExtend(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtend", "saveCustomerExtend", e, now)
		ormErrorLog(repo, "CustomerExtend", "saveCustomerExtend", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomer .
func findCustomer(repo GORMRepository, result *po.Customer, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomer", e, now)
		ormErrorLog(repo, "Customer", "findCustomer", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerListByPrimarys .
func findCustomerListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.Customer, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomerListByPrimarys", e, now)
		ormErrorLog(repo, "Customer", "findCustomersByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerByWhere .
func findCustomerByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.Customer, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomerByWhere", e, now)
		ormErrorLog(repo, "Customer", "findCustomerByWhere", e, query, args)
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

// findCustomerByMap .
func findCustomerByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.Customer, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomerByMap", e, now)
		ormErrorLog(repo, "Customer", "findCustomerByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerList .
func findCustomerList(repo GORMRepository, query po.Customer, builders ...Builder) (results []po.Customer, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomerList", e, now)
		ormErrorLog(repo, "Customer", "findCustomers", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerListByWhere .
func findCustomerListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.Customer, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomerListByWhere", e, now)
		ormErrorLog(repo, "Customer", "findCustomersByWhere", e, query, args)
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

// findCustomerListByMap .
func findCustomerListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.Customer, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "findCustomerListByMap", e, now)
		ormErrorLog(repo, "Customer", "findCustomersByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomer .
func createCustomer(repo GORMRepository, object *po.Customer) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "createCustomer", e, now)
		ormErrorLog(repo, "Customer", "createCustomer", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomer .
func saveCustomer(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Customer", "saveCustomer", e, now)
		ormErrorLog(repo, "Customer", "saveCustomer", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}
