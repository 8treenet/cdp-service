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

// findCustomerExtensionTemplate .
func findCustomerExtensionTemplate(repo GORMRepository, result *po.CustomerExtensionTemplate, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplate", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplate", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerExtensionTemplateListByPrimarys .
func findCustomerExtensionTemplateListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerExtensionTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplateListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplatesByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerExtensionTemplateByWhere .
func findCustomerExtensionTemplateByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerExtensionTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplateByWhere", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplateByWhere", e, query, args)
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

// findCustomerExtensionTemplateByMap .
func findCustomerExtensionTemplateByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerExtensionTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplateByMap", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplateByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerExtensionTemplateList .
func findCustomerExtensionTemplateList(repo GORMRepository, query po.CustomerExtensionTemplate, builders ...Builder) (results []po.CustomerExtensionTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplateList", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplates", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerExtensionTemplateListByWhere .
func findCustomerExtensionTemplateListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerExtensionTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplateListByWhere", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplatesByWhere", e, query, args)
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

// findCustomerExtensionTemplateListByMap .
func findCustomerExtensionTemplateListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerExtensionTemplate, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "findCustomerExtensionTemplateListByMap", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "findCustomerExtensionTemplatesByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerExtensionTemplate .
func createCustomerExtensionTemplate(repo GORMRepository, object *po.CustomerExtensionTemplate) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "createCustomerExtensionTemplate", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "createCustomerExtensionTemplate", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerExtensionTemplate .
func saveCustomerExtensionTemplate(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionTemplate", "saveCustomerExtensionTemplate", e, now)
		ormErrorLog(repo, "CustomerExtensionTemplate", "saveCustomerExtensionTemplate", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerExtension .
func findCustomerExtension(repo GORMRepository, result *po.CustomerExtension, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtension", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtension", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerExtensionListByPrimarys .
func findCustomerExtensionListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerExtension, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtensionListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtensionsByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerExtensionByWhere .
func findCustomerExtensionByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerExtension, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtensionByWhere", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtensionByWhere", e, query, args)
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

// findCustomerExtensionByMap .
func findCustomerExtensionByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerExtension, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtensionByMap", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtensionByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerExtensionList .
func findCustomerExtensionList(repo GORMRepository, query po.CustomerExtension, builders ...Builder) (results []po.CustomerExtension, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtensionList", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtensions", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerExtensionListByWhere .
func findCustomerExtensionListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerExtension, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtensionListByWhere", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtensionsByWhere", e, query, args)
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

// findCustomerExtensionListByMap .
func findCustomerExtensionListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerExtension, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "findCustomerExtensionListByMap", e, now)
		ormErrorLog(repo, "CustomerExtension", "findCustomerExtensionsByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerExtension .
func createCustomerExtension(repo GORMRepository, object *po.CustomerExtension) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "createCustomerExtension", e, now)
		ormErrorLog(repo, "CustomerExtension", "createCustomerExtension", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerExtension .
func saveCustomerExtension(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtension", "saveCustomerExtension", e, now)
		ormErrorLog(repo, "CustomerExtension", "saveCustomerExtension", e, object)
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
