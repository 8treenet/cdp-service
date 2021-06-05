package repository

import (
	"context"
	"time"

	"github.com/8treenet/crm-service/domain/entity"
	"github.com/8treenet/crm-service/domain/po"
	"github.com/8treenet/crm-service/infra"
	"github.com/8treenet/freedom"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

const (
	noneIndex   = 0
	btreeIndex  = 1
	uniqueIndex = 2
	textIndex   = 3
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindRepository(func() *CustomerRepository {
			return &CustomerRepository{customerCollection: "cdp_customer", customerTemplateCacheKey: "cdp_customer_template"}
		})
	})
}

// CustomerRepository .
type CustomerRepository struct {
	freedom.Repository
	Common                   *infra.CommonRequest
	Mongo                    *infra.Mongo
	customerCollection       string
	customerTemplateCacheKey string
}

// AddTempleteField .
func (repo *CustomerRepository) AddTempleteField(name, kind, dict string, index int) error {
	defer func() {
		if e := repo.Redis().Del(repo.customerTemplateCacheKey).Err(); e != nil {
			repo.Worker().Logger().Error(e)
		}
	}()
	pobject := &po.CustomerTemplate{
		Name:    name,
		Kind:    kind,
		Index:   index,
		Dict:    dict,
		Created: time.Now(),
		Updated: time.Now(),
	}
	_, err := createCustomerTemplate(repo, pobject)
	if err != nil {
		return err
	}

	opt := options.Index().SetBackground(true).SetName(name)
	var sort interface{} = -1
	switch index {
	case noneIndex:
		return nil
	case uniqueIndex:
		opt.SetUnique(true)
	case textIndex:
		sort = "text"
	}
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{name, sort},
		},
		Options: opt,
	}

	collection := repo.Mongo.GetCollection(repo.customerCollection)
	_, err = collection.Indexes().CreateOne(context.TODO(), indexModel)
	return err
}

// NewCustomer .
func (repo *CustomerRepository) NewCustomer(source map[string]interface{}) (*entity.Customer, error) {
	result := &entity.Customer{}
	result.Source = source
	repo.InjectBaseEntity(result)

	templetes, err := repo.getTempletes()
	if err != nil {
		return nil, err
	}

	result.Templetes = templetes
	result.Source["_id"] = primitive.NewObjectID().Hex()
	result.Source["_created"] = time.Now()

	if e := result.Verify(); e != nil {
		return nil, e
	}
	collection := repo.Mongo.GetCollection(repo.customerCollection)
	_, err = collection.InsertOne(context.TODO(), result.Source)
	return result, err
}

// NewCustomers .
func (repo *CustomerRepository) NewCustomers(sources []map[string]interface{}) ([]*entity.Customer, error) {
	result := []*entity.Customer{}

	templetes, err := repo.getTempletes()
	if err != nil {
		return nil, err
	}

	created := time.Now()
	for _, source := range sources {
		source["_created"] = created
		source["_id"] = primitive.NewObjectID().Hex()
		entity := &entity.Customer{
			Source:    source,
			Templetes: templetes,
		}
		if e := entity.Verify(); e != nil {
			return nil, e
		}
		result = append(result, entity)
	}
	repo.InjectBaseEntitys(result)

	collection := repo.Mongo.GetCollection(repo.customerCollection)
	_, err = collection.InsertMany(context.TODO(), repo.Mongo.ToDocuments(sources))
	return result, err
}

// SaveCustomer .
func (repo *CustomerRepository) SaveCustomer(customer *entity.Customer) error {
	collection := repo.Mongo.GetCollection(repo.customerCollection)
	updateOpt := options.Update().SetUpsert(true)
	if e := customer.Verify(); e != nil {
		return e
	}

	value := bson.M{
		"$set": customer.GetChanges(),
	}
	_, err := collection.UpdateOne(context.TODO(), customer.Location(), value, updateOpt)
	return err
}

// GetCustomer .
func (repo *CustomerRepository) GetCustomer(id string) (result *entity.Customer, e error) {
	templetes, err := repo.getTempletes()
	if err != nil {
		return nil, err
	}

	collection := repo.Mongo.GetCollection(repo.customerCollection)
	singleResult := collection.FindOne(context.TODO(), map[string]interface{}{"_id": id})
	if e = singleResult.Err(); e != nil {
		return nil, e
	}

	result = &entity.Customer{}
	result.Source = map[string]interface{}{}
	result.Templetes = templetes
	repo.InjectBaseEntity(result)
	e = singleResult.Decode(&result.Source)
	return
}

func (repo *CustomerRepository) getTempletes() (result []*po.CustomerTemplate, err error) {
	err = redisJSONGet(repo.Redis(), repo.customerTemplateCacheKey, &result)
	if err == nil || err != redis.Nil {
		return
	}

	list, err := findCustomerTemplateList(repo, po.CustomerTemplate{})
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		result = append(result, &list[i])
	}
	err = redisJSONSet(repo.Redis(), repo.customerTemplateCacheKey, result)
	return
}

// db .
func (repo *CustomerRepository) db() *gorm.DB {
	var db *gorm.DB
	if err := repo.FetchDB(&db); err != nil {
		panic(err)
	}
	return db
}
