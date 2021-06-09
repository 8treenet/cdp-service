package infra

/*
	需要在main中 隐式初始化组件
  	import _"github.com/8treenet/cdp-service/infra"
*/

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/8treenet/freedom"
	driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		initiator.BindInfra(true, &Mongo{})
	})
}

// Mongo .
type Mongo struct {
	freedom.Infra
	db            *driver.Client
	opts          *options.ClientOptions
	databaseName  string
	database      *driver.Database
	collectionMap map[string]*driver.Collection
	mutex         sync.Mutex
}

// Start .
func (mongo *Mongo) visitConfig() {
	var mongoConf struct {
		Database string `toml:"database"`
		URI      string `toml:"uri"`
	}
	//default
	mongoConf.Database = "default"
	mongoConf.URI = "mongodb://root:123123@localhost:27017"

	freedom.Configure(&mongoConf, "mongo.toml")
	mongo.opts = options.Client().ApplyURI(mongoConf.URI)
	mongo.databaseName = mongoConf.Database
}

// Booting .
func (mongo *Mongo) Booting(bootManager freedom.BootManager) {
	mongo.visitConfig()

	ctx, fun := context.WithTimeout(context.Background(), time.Second*2)
	defer fun()
	client, err := driver.Connect(ctx, mongo.opts)
	if err != nil {
		freedom.Logger().Fatal(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		freedom.Logger().Fatal(err)
	}

	bootManager.RegisterShutdown(func() {
		client.Disconnect(ctx)
	})

	mongo.db = client
	mongo.database = mongo.db.Database(mongo.databaseName)
	mongo.collectionMap = make(map[string]*driver.Collection)
	freedom.Logger().Infof("Connect mongo success, database:%s.", mongo.databaseName)
}

func (mongo *Mongo) GetCollection(data interface{}) *driver.Collection {
	type getCollection interface {
		Collection() string
	}

	collectionName := fmt.Sprint(data)
	impl, ok := data.(getCollection)
	if ok {
		collectionName = impl.Collection()
	}

	mongo.mutex.Lock()
	defer mongo.mutex.Unlock()

	if result, ok := mongo.collectionMap[collectionName]; ok {
		return result
	}

	result := mongo.database.Collection(collectionName)
	mongo.collectionMap[collectionName] = result
	return result
}

func (mongo *Mongo) ToDocuments(data interface{}) (documents []interface{}) {
	value := reflect.ValueOf(data)
	if value.Type().Kind() != reflect.Slice && value.Type().Kind() != reflect.Array {
		freedom.Logger().Error("ToDocuments:传入类型data错误,必须是一个数组.")
		return
	}

	for i := 0; i < value.Len(); i++ {
		documents = append(documents, value.Index(i).Interface())
	}
	return
}
