package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var defExpiration time.Duration = time.Minute * 10

func redisJSONGet(client redis.Cmdable, key string, value interface{}) error {
	bytes, err := client.Get(key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

func redisJSONSet(client redis.Cmdable, key string, value interface{}, expirations ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	expiration := defExpiration
	if len(expirations) > 0 {
		expiration = expirations[0]
	}
	return client.Set(key, data, expiration).Err()
}

// Limit .
type Limit struct {
	length int
	fields []string
	orders []string
}

// NewDescLimit .
func NewDescLimit(column string, columns ...string) *Limit {
	return newDefaultLimit("desc", column, columns...)
}

// NewAscLimit .
func NewAscLimit(column string, columns ...string) *Limit {
	return newDefaultLimit("asc", column, columns...)
}

// NewDescOrder .
func newDefaultLimit(sort, field string, args ...string) *Limit {
	fields := []string{field}
	fields = append(fields, args...)
	orders := []string{}
	for index := 0; index < len(fields); index++ {
		orders = append(orders, sort)
	}
	return &Limit{
		fields: fields,
		orders: orders,
	}
}

// Order .
func (p *Limit) Order() interface{} {
	if len(p.fields) == 0 {
		return nil
	}
	args := []string{}
	for index := 0; index < len(p.fields); index++ {
		args = append(args, fmt.Sprintf("`%s` %s", p.fields[index], p.orders[index]))
	}

	return strings.Join(args, ",")
}

// SetLength .
func (p *Limit) SetLength(length int) *Limit {
	p.length = length
	return p
}

// Execute .
func (p *Limit) Execute(db *gorm.DB, object interface{}) error {
	orderValue := p.Order()
	if orderValue != nil {
		db = db.Order(orderValue)
	}
	db = db.Limit(p.length)

	resultDB := db.Find(object)
	if resultDB.Error != nil {
		return resultDB.Error
	}
	return nil
}
