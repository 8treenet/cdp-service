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

// findSystemConfig .
func findSystemConfig(repo GORMRepository, result *po.SystemConfig, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfig", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfig", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findSystemConfigListByPrimarys .
func findSystemConfigListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.SystemConfig, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfigListByPrimarys", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfigsByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findSystemConfigByWhere .
func findSystemConfigByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.SystemConfig, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfigByWhere", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfigByWhere", e, query, args)
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

// findSystemConfigByMap .
func findSystemConfigByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.SystemConfig, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfigByMap", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfigByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findSystemConfigList .
func findSystemConfigList(repo GORMRepository, query po.SystemConfig, builders ...Builder) (results []po.SystemConfig, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfigList", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfigs", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findSystemConfigListByWhere .
func findSystemConfigListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.SystemConfig, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfigListByWhere", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfigsByWhere", e, query, args)
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

// findSystemConfigListByMap .
func findSystemConfigListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.SystemConfig, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "findSystemConfigListByMap", e, now)
		ormErrorLog(repo, "SystemConfig", "findSystemConfigsByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createSystemConfig .
func createSystemConfig(repo GORMRepository, object *po.SystemConfig) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "createSystemConfig", e, now)
		ormErrorLog(repo, "SystemConfig", "createSystemConfig", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveSystemConfig .
func saveSystemConfig(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("SystemConfig", "saveSystemConfig", e, now)
		ormErrorLog(repo, "SystemConfig", "saveSystemConfig", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findSource .
func findSource(repo GORMRepository, result *po.Source, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSource", e, now)
		ormErrorLog(repo, "Source", "findSource", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findSourceListByPrimarys .
func findSourceListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.Source, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSourceListByPrimarys", e, now)
		ormErrorLog(repo, "Source", "findSourcesByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findSourceByWhere .
func findSourceByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.Source, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSourceByWhere", e, now)
		ormErrorLog(repo, "Source", "findSourceByWhere", e, query, args)
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

// findSourceByMap .
func findSourceByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.Source, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSourceByMap", e, now)
		ormErrorLog(repo, "Source", "findSourceByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findSourceList .
func findSourceList(repo GORMRepository, query po.Source, builders ...Builder) (results []po.Source, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSourceList", e, now)
		ormErrorLog(repo, "Source", "findSources", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findSourceListByWhere .
func findSourceListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.Source, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSourceListByWhere", e, now)
		ormErrorLog(repo, "Source", "findSourcesByWhere", e, query, args)
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

// findSourceListByMap .
func findSourceListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.Source, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "findSourceListByMap", e, now)
		ormErrorLog(repo, "Source", "findSourcesByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createSource .
func createSource(repo GORMRepository, object *po.Source) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "createSource", e, now)
		ormErrorLog(repo, "Source", "createSource", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveSource .
func saveSource(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Source", "saveSource", e, now)
		ormErrorLog(repo, "Source", "saveSource", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findIPAddr .
func findIPAddr(repo GORMRepository, result *po.IPAddr, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddr", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddr", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findIPAddrListByPrimarys .
func findIPAddrListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.IPAddr, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddrListByPrimarys", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddrsByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findIPAddrByWhere .
func findIPAddrByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.IPAddr, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddrByWhere", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddrByWhere", e, query, args)
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

// findIPAddrByMap .
func findIPAddrByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.IPAddr, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddrByMap", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddrByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findIPAddrList .
func findIPAddrList(repo GORMRepository, query po.IPAddr, builders ...Builder) (results []po.IPAddr, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddrList", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddrs", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findIPAddrListByWhere .
func findIPAddrListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.IPAddr, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddrListByWhere", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddrsByWhere", e, query, args)
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

// findIPAddrListByMap .
func findIPAddrListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.IPAddr, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "findIPAddrListByMap", e, now)
		ormErrorLog(repo, "IPAddr", "findIPAddrsByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createIPAddr .
func createIPAddr(repo GORMRepository, object *po.IPAddr) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "createIPAddr", e, now)
		ormErrorLog(repo, "IPAddr", "createIPAddr", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveIPAddr .
func saveIPAddr(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("IPAddr", "saveIPAddr", e, now)
		ormErrorLog(repo, "IPAddr", "saveIPAddr", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerWechat .
func findCustomerWechat(repo GORMRepository, result *po.CustomerWechat, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechat", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechat", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerWechatListByPrimarys .
func findCustomerWechatListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerWechat, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechatListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechatsByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerWechatByWhere .
func findCustomerWechatByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerWechat, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechatByWhere", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechatByWhere", e, query, args)
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

// findCustomerWechatByMap .
func findCustomerWechatByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerWechat, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechatByMap", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechatByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerWechatList .
func findCustomerWechatList(repo GORMRepository, query po.CustomerWechat, builders ...Builder) (results []po.CustomerWechat, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechatList", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechats", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerWechatListByWhere .
func findCustomerWechatListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerWechat, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechatListByWhere", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechatsByWhere", e, query, args)
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

// findCustomerWechatListByMap .
func findCustomerWechatListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerWechat, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "findCustomerWechatListByMap", e, now)
		ormErrorLog(repo, "CustomerWechat", "findCustomerWechatsByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerWechat .
func createCustomerWechat(repo GORMRepository, object *po.CustomerWechat) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "createCustomerWechat", e, now)
		ormErrorLog(repo, "CustomerWechat", "createCustomerWechat", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerWechat .
func saveCustomerWechat(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerWechat", "saveCustomerWechat", e, now)
		ormErrorLog(repo, "CustomerWechat", "saveCustomerWechat", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerTemporary .
func findCustomerTemporary(repo GORMRepository, result *po.CustomerTemporary, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporary", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporary", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerTemporaryListByPrimarys .
func findCustomerTemporaryListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerTemporary, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporaryListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporarysByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerTemporaryByWhere .
func findCustomerTemporaryByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerTemporary, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporaryByWhere", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporaryByWhere", e, query, args)
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

// findCustomerTemporaryByMap .
func findCustomerTemporaryByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerTemporary, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporaryByMap", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporaryByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerTemporaryList .
func findCustomerTemporaryList(repo GORMRepository, query po.CustomerTemporary, builders ...Builder) (results []po.CustomerTemporary, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporaryList", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporarys", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerTemporaryListByWhere .
func findCustomerTemporaryListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerTemporary, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporaryListByWhere", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporarysByWhere", e, query, args)
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

// findCustomerTemporaryListByMap .
func findCustomerTemporaryListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerTemporary, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "findCustomerTemporaryListByMap", e, now)
		ormErrorLog(repo, "CustomerTemporary", "findCustomerTemporarysByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerTemporary .
func createCustomerTemporary(repo GORMRepository, object *po.CustomerTemporary) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "createCustomerTemporary", e, now)
		ormErrorLog(repo, "CustomerTemporary", "createCustomerTemporary", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerTemporary .
func saveCustomerTemporary(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerTemporary", "saveCustomerTemporary", e, now)
		ormErrorLog(repo, "CustomerTemporary", "saveCustomerTemporary", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerPhone .
func findCustomerPhone(repo GORMRepository, result *po.CustomerPhone, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhone", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhone", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerPhoneListByPrimarys .
func findCustomerPhoneListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerPhone, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhoneListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhonesByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerPhoneByWhere .
func findCustomerPhoneByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerPhone, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhoneByWhere", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhoneByWhere", e, query, args)
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

// findCustomerPhoneByMap .
func findCustomerPhoneByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerPhone, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhoneByMap", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhoneByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerPhoneList .
func findCustomerPhoneList(repo GORMRepository, query po.CustomerPhone, builders ...Builder) (results []po.CustomerPhone, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhoneList", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhones", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerPhoneListByWhere .
func findCustomerPhoneListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerPhone, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhoneListByWhere", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhonesByWhere", e, query, args)
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

// findCustomerPhoneListByMap .
func findCustomerPhoneListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerPhone, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "findCustomerPhoneListByMap", e, now)
		ormErrorLog(repo, "CustomerPhone", "findCustomerPhonesByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerPhone .
func createCustomerPhone(repo GORMRepository, object *po.CustomerPhone) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "createCustomerPhone", e, now)
		ormErrorLog(repo, "CustomerPhone", "createCustomerPhone", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerPhone .
func saveCustomerPhone(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerPhone", "saveCustomerPhone", e, now)
		ormErrorLog(repo, "CustomerPhone", "saveCustomerPhone", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerKey .
func findCustomerKey(repo GORMRepository, result *po.CustomerKey, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKey", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKey", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerKeyListByPrimarys .
func findCustomerKeyListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerKey, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKeyListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKeysByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerKeyByWhere .
func findCustomerKeyByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerKey, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKeyByWhere", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKeyByWhere", e, query, args)
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

// findCustomerKeyByMap .
func findCustomerKeyByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerKey, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKeyByMap", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKeyByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerKeyList .
func findCustomerKeyList(repo GORMRepository, query po.CustomerKey, builders ...Builder) (results []po.CustomerKey, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKeyList", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKeys", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerKeyListByWhere .
func findCustomerKeyListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerKey, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKeyListByWhere", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKeysByWhere", e, query, args)
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

// findCustomerKeyListByMap .
func findCustomerKeyListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerKey, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "findCustomerKeyListByMap", e, now)
		ormErrorLog(repo, "CustomerKey", "findCustomerKeysByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerKey .
func createCustomerKey(repo GORMRepository, object *po.CustomerKey) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "createCustomerKey", e, now)
		ormErrorLog(repo, "CustomerKey", "createCustomerKey", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerKey .
func saveCustomerKey(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerKey", "saveCustomerKey", e, now)
		ormErrorLog(repo, "CustomerKey", "saveCustomerKey", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findCustomerExtensionMetadata .
func findCustomerExtensionMetadata(repo GORMRepository, result *po.CustomerExtensionMetadata, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadata", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadata", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findCustomerExtensionMetadataListByPrimarys .
func findCustomerExtensionMetadataListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.CustomerExtensionMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadataListByPrimarys", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadatasByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findCustomerExtensionMetadataByWhere .
func findCustomerExtensionMetadataByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.CustomerExtensionMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadataByWhere", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadataByWhere", e, query, args)
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

// findCustomerExtensionMetadataByMap .
func findCustomerExtensionMetadataByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.CustomerExtensionMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadataByMap", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadataByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findCustomerExtensionMetadataList .
func findCustomerExtensionMetadataList(repo GORMRepository, query po.CustomerExtensionMetadata, builders ...Builder) (results []po.CustomerExtensionMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadataList", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadatas", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findCustomerExtensionMetadataListByWhere .
func findCustomerExtensionMetadataListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.CustomerExtensionMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadataListByWhere", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadatasByWhere", e, query, args)
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

// findCustomerExtensionMetadataListByMap .
func findCustomerExtensionMetadataListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.CustomerExtensionMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "findCustomerExtensionMetadataListByMap", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "findCustomerExtensionMetadatasByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createCustomerExtensionMetadata .
func createCustomerExtensionMetadata(repo GORMRepository, object *po.CustomerExtensionMetadata) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "createCustomerExtensionMetadata", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "createCustomerExtensionMetadata", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveCustomerExtensionMetadata .
func saveCustomerExtensionMetadata(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("CustomerExtensionMetadata", "saveCustomerExtensionMetadata", e, now)
		ormErrorLog(repo, "CustomerExtensionMetadata", "saveCustomerExtensionMetadata", e, object)
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

// findBehaviourFeatureMetadata .
func findBehaviourFeatureMetadata(repo GORMRepository, result *po.BehaviourFeatureMetadata, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadata", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadata", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findBehaviourFeatureMetadataListByPrimarys .
func findBehaviourFeatureMetadataListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.BehaviourFeatureMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadataListByPrimarys", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadatasByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findBehaviourFeatureMetadataByWhere .
func findBehaviourFeatureMetadataByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.BehaviourFeatureMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadataByWhere", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadataByWhere", e, query, args)
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

// findBehaviourFeatureMetadataByMap .
func findBehaviourFeatureMetadataByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.BehaviourFeatureMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadataByMap", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadataByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findBehaviourFeatureMetadataList .
func findBehaviourFeatureMetadataList(repo GORMRepository, query po.BehaviourFeatureMetadata, builders ...Builder) (results []po.BehaviourFeatureMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadataList", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadatas", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findBehaviourFeatureMetadataListByWhere .
func findBehaviourFeatureMetadataListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.BehaviourFeatureMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadataListByWhere", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadatasByWhere", e, query, args)
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

// findBehaviourFeatureMetadataListByMap .
func findBehaviourFeatureMetadataListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.BehaviourFeatureMetadata, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "findBehaviourFeatureMetadataListByMap", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "findBehaviourFeatureMetadatasByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createBehaviourFeatureMetadata .
func createBehaviourFeatureMetadata(repo GORMRepository, object *po.BehaviourFeatureMetadata) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "createBehaviourFeatureMetadata", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "createBehaviourFeatureMetadata", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveBehaviourFeatureMetadata .
func saveBehaviourFeatureMetadata(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeatureMetadata", "saveBehaviourFeatureMetadata", e, now)
		ormErrorLog(repo, "BehaviourFeatureMetadata", "saveBehaviourFeatureMetadata", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findBehaviourFeature .
func findBehaviourFeature(repo GORMRepository, result *po.BehaviourFeature, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeature", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeature", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findBehaviourFeatureListByPrimarys .
func findBehaviourFeatureListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.BehaviourFeature, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeatureListByPrimarys", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeaturesByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findBehaviourFeatureByWhere .
func findBehaviourFeatureByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.BehaviourFeature, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeatureByWhere", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeatureByWhere", e, query, args)
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

// findBehaviourFeatureByMap .
func findBehaviourFeatureByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.BehaviourFeature, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeatureByMap", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeatureByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findBehaviourFeatureList .
func findBehaviourFeatureList(repo GORMRepository, query po.BehaviourFeature, builders ...Builder) (results []po.BehaviourFeature, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeatureList", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeatures", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findBehaviourFeatureListByWhere .
func findBehaviourFeatureListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.BehaviourFeature, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeatureListByWhere", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeaturesByWhere", e, query, args)
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

// findBehaviourFeatureListByMap .
func findBehaviourFeatureListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.BehaviourFeature, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "findBehaviourFeatureListByMap", e, now)
		ormErrorLog(repo, "BehaviourFeature", "findBehaviourFeaturesByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createBehaviourFeature .
func createBehaviourFeature(repo GORMRepository, object *po.BehaviourFeature) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "createBehaviourFeature", e, now)
		ormErrorLog(repo, "BehaviourFeature", "createBehaviourFeature", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveBehaviourFeature .
func saveBehaviourFeature(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("BehaviourFeature", "saveBehaviourFeature", e, now)
		ormErrorLog(repo, "BehaviourFeature", "saveBehaviourFeature", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findBehaviour .
func findBehaviour(repo GORMRepository, result *po.Behaviour, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviour", e, now)
		ormErrorLog(repo, "Behaviour", "findBehaviour", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findBehaviourListByPrimarys .
func findBehaviourListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.Behaviour, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviourListByPrimarys", e, now)
		ormErrorLog(repo, "Behaviour", "findBehavioursByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findBehaviourByWhere .
func findBehaviourByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.Behaviour, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviourByWhere", e, now)
		ormErrorLog(repo, "Behaviour", "findBehaviourByWhere", e, query, args)
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

// findBehaviourByMap .
func findBehaviourByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.Behaviour, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviourByMap", e, now)
		ormErrorLog(repo, "Behaviour", "findBehaviourByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findBehaviourList .
func findBehaviourList(repo GORMRepository, query po.Behaviour, builders ...Builder) (results []po.Behaviour, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviourList", e, now)
		ormErrorLog(repo, "Behaviour", "findBehaviours", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findBehaviourListByWhere .
func findBehaviourListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.Behaviour, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviourListByWhere", e, now)
		ormErrorLog(repo, "Behaviour", "findBehavioursByWhere", e, query, args)
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

// findBehaviourListByMap .
func findBehaviourListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.Behaviour, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "findBehaviourListByMap", e, now)
		ormErrorLog(repo, "Behaviour", "findBehavioursByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createBehaviour .
func createBehaviour(repo GORMRepository, object *po.Behaviour) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "createBehaviour", e, now)
		ormErrorLog(repo, "Behaviour", "createBehaviour", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveBehaviour .
func saveBehaviour(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Behaviour", "saveBehaviour", e, now)
		ormErrorLog(repo, "Behaviour", "saveBehaviour", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findAnalysisReport .
func findAnalysisReport(repo GORMRepository, result *po.AnalysisReport, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReport", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReport", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findAnalysisReportListByPrimarys .
func findAnalysisReportListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.AnalysisReport, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReportListByPrimarys", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReportsByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findAnalysisReportByWhere .
func findAnalysisReportByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.AnalysisReport, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReportByWhere", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReportByWhere", e, query, args)
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

// findAnalysisReportByMap .
func findAnalysisReportByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.AnalysisReport, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReportByMap", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReportByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findAnalysisReportList .
func findAnalysisReportList(repo GORMRepository, query po.AnalysisReport, builders ...Builder) (results []po.AnalysisReport, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReportList", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReports", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findAnalysisReportListByWhere .
func findAnalysisReportListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.AnalysisReport, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReportListByWhere", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReportsByWhere", e, query, args)
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

// findAnalysisReportListByMap .
func findAnalysisReportListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.AnalysisReport, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "findAnalysisReportListByMap", e, now)
		ormErrorLog(repo, "AnalysisReport", "findAnalysisReportsByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createAnalysisReport .
func createAnalysisReport(repo GORMRepository, object *po.AnalysisReport) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "createAnalysisReport", e, now)
		ormErrorLog(repo, "AnalysisReport", "createAnalysisReport", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveAnalysisReport .
func saveAnalysisReport(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("AnalysisReport", "saveAnalysisReport", e, now)
		ormErrorLog(repo, "AnalysisReport", "saveAnalysisReport", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}

// findAnalysis .
func findAnalysis(repo GORMRepository, result *po.Analysis, builders ...Builder) (e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysis", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysis", e, result)
	}()
	db := repo.db()
	if len(builders) == 0 {
		e = db.Where(result).Last(result).Error
		return
	}
	e = builders[0].Execute(db.Limit(1), result)
	return
}

// findAnalysisListByPrimarys .
func findAnalysisListByPrimarys(repo GORMRepository, primarys ...interface{}) (results []po.Analysis, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysisListByPrimarys", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysissByPrimarys", e, primarys)
	}()

	e = repo.db().Find(&results, primarys).Error
	return
}

// findAnalysisByWhere .
func findAnalysisByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (result po.Analysis, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysisByWhere", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysisByWhere", e, query, args)
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

// findAnalysisByMap .
func findAnalysisByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (result po.Analysis, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysisByMap", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysisByMap", e, query)
	}()

	db := repo.db().Where(query)
	if len(builders) == 0 {
		e = db.Last(&result).Error
		return
	}

	e = builders[0].Execute(db.Limit(1), &result)
	return
}

// findAnalysisList .
func findAnalysisList(repo GORMRepository, query po.Analysis, builders ...Builder) (results []po.Analysis, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysisList", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysiss", e, query)
	}()
	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// findAnalysisListByWhere .
func findAnalysisListByWhere(repo GORMRepository, query string, args []interface{}, builders ...Builder) (results []po.Analysis, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysisListByWhere", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysissByWhere", e, query, args)
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

// findAnalysisListByMap .
func findAnalysisListByMap(repo GORMRepository, query map[string]interface{}, builders ...Builder) (results []po.Analysis, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "findAnalysisListByMap", e, now)
		ormErrorLog(repo, "Analysis", "findAnalysissByMap", e, query)
	}()

	db := repo.db().Where(query)

	if len(builders) == 0 {
		e = db.Find(&results).Error
		return
	}
	e = builders[0].Execute(db, &results)
	return
}

// createAnalysis .
func createAnalysis(repo GORMRepository, object *po.Analysis) (rowsAffected int64, e error) {
	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "createAnalysis", e, now)
		ormErrorLog(repo, "Analysis", "createAnalysis", e, *object)
	}()

	db := repo.db().Create(object)
	rowsAffected = db.RowsAffected
	e = db.Error
	return
}

// saveAnalysis .
func saveAnalysis(repo GORMRepository, object saveObject) (rowsAffected int64, e error) {
	if len(object.Location()) == 0 {
		return 0, errors.New("location cannot be empty")
	}

	now := time.Now()
	defer func() {
		freedom.Prometheus().OrmWithLabelValues("Analysis", "saveAnalysis", e, now)
		ormErrorLog(repo, "Analysis", "saveAnalysis", e, object)
	}()

	db := repo.db().Table(object.TableName()).Where(object.Location()).Updates(object.GetChanges())
	e = db.Error
	rowsAffected = db.RowsAffected
	return
}
