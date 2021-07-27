package geo

import (
	"sync"
	"time"

	"github.com/8treenet/cdp-service/utils"
	"github.com/8treenet/freedom"
	"github.com/8treenet/freedom/infra/requests"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	freedom.Prepare(func(initiator freedom.Initiator) {
		obj := &GEO{}
		obj.ipcom.maxCount = 15
		obj.ipcom.diffSec = 65
		obj.ipcom.firstTime = time.Now()
		obj.client = requests.NewHTTPClient(30*time.Second, 10*time.Second)

		initiator.BindInfra(true, obj)
		initiator.InjectController(func(ctx freedom.Context) (com *GEO) {
			initiator.FetchInfra(ctx, &com)
			return
		})
	})
}

type dispatcher interface {
	Idle() bool
	Do(requests.Client, []string, map[string]*GEOInfo) bool
}

// GEO
type GEO struct {
	freedom.Infra
	client           requests.Client
	backupDispatcher []dispatcher
	ipcom            struct {
		lock      sync.Mutex
		count     int       //使用次数
		maxCount  int       //最大使用次数
		diffSec   int64     //时间限制
		firstTime time.Time //最近的使用时间
	}
}

type GEOInfo struct {
	Addr    string `gorm:"column:addr" json:"query"`
	Country string `gorm:"column:country" json:"country"`
	Region  string `gorm:"column:region" json:"regionName"`
	City    string `gorm:"column:city" json:"city"`
}

// Start .
func (geo *GEO) visitConfig() {
}

// Booting .
func (geo *GEO) Booting(bootManager freedom.BootManager) {
	geo.visitConfig()
}

// ParseIP
func (geo *GEO) ParseIP(addr string) (*GEOInfo, error) {
	m, err := geo.ParseBatchIP([]string{addr})
	if err != nil {
		return nil, err
	}
	info, _ := m[addr]
	return info, nil
}

// ParseIP
func (geo *GEO) ParseBatchIP(addr []string) (map[string]*GEOInfo, error) {
	mapData := geo.fetchDB(addr)
	missAddr := []string{}
	for _, v := range addr {
		if _, ok := mapData[v]; ok {
			continue
		}
		missAddr = append(missAddr, v)
	}

	if len(missAddr) == 0 {
		return mapData, nil
	}

	ipMapData, err := geo.fetchIPCom(missAddr)
	if err != nil {
		return nil, err
	}
	for key, v := range ipMapData {
		mapData[key] = v
	}
	geo.insertDB(ipMapData)
	return mapData, nil
}

func (geo *GEO) fetchDB(addr []string) (result map[string]*GEOInfo) {
	result = make(map[string]*GEOInfo)
	var db *gorm.DB
	if err := geo.FetchOnlyDB(&db); err != nil {
		freedom.Logger().Error("geo fetchDB error:", err)
		return
	}

	var list []*GEOInfo
	if err := db.Table("cdp_ip_addr").Where("addr in (?)", addr).Find(&list).Error; err != nil {
		freedom.Logger().Error("geo find error:", err)
		return
	}

	for i := 0; i < len(list); i++ {
		result[list[i].Addr] = list[i]
	}
	return
}

func (geo *GEO) insertDB(mapData map[string]*GEOInfo) {
	var db *gorm.DB
	if err := geo.FetchOnlyDB(&db); err != nil {
		freedom.Logger().Error("geo insertDB error:", err)
		return
	}
	var list []*GEOInfo
	for _, v := range mapData {
		list = append(list, v)
	}
	if err := db.Clauses(clause.Insert{Modifier: "IGNORE"}).Table("cdp_ip_addr").CreateInBatches(list, 500).Error; err != nil {
		freedom.Logger().Error("geo CreateInBatches error:", err)
	}
	return
}

func (geo *GEO) fetchIPCom(addr []string) (result map[string]*GEOInfo, err error) {
	geo.ipcom.lock.Lock()
	defer geo.ipcom.lock.Unlock()
	result = make(map[string]*GEOInfo)

	var list [][]string
	err = utils.SliceUp(addr, &list, 100)
	if err != nil {
		freedom.Logger().Error("geo fetchIPCom SliceUp error:", err)
		return
	}

	check := func() bool {
		now := time.Now()
		//执行次数等于最大 并且 当前时间-61 小于首次时间
		if geo.ipcom.count >= geo.ipcom.maxCount && now.Unix()-geo.ipcom.diffSec < geo.ipcom.firstTime.Unix() {
			freedom.Logger().Info("等待ipcom限制", geo.ipcom.count, geo.ipcom.maxCount, now.Unix()-geo.ipcom.diffSec, geo.ipcom.firstTime.Unix())
			return false
		}
		return true
	}

	reset := func() {
		geo.ipcom.count++
		now := time.Now()
		if now.Unix()-geo.ipcom.diffSec < geo.ipcom.firstTime.Unix() {
			return
		}
		geo.ipcom.count = 1
		geo.ipcom.firstTime = now
	}
	for i := 0; i < len(list); i++ {
		if !check() {
			if geo.backup(list[i], result) {
				continue //备用成功
			}
			time.Sleep(1 * time.Second)
			i--
			continue
		}
		reset()

		resp := []*GEOInfo{}
		if e := geo.NewHTTPRequest("http://ip-api.com/batch?lang=zh-CN&fields=city,country,query,regionName").Post().SetClient(geo.client).SetJSONBody(list[i]).ToJSON(&resp).Error; e != nil {
			return nil, e
		}
		for _, v := range resp {
			result[v.Addr] = v
		}
	}
	return
}

func (geo *GEO) backup(addr []string, result map[string]*GEOInfo) bool {
	if len(geo.backupDispatcher) == 0 {
		return false
	}

	for _, v := range geo.backupDispatcher {
		if !v.Idle() {
			continue
		}
		return v.Do(geo.client, addr, result)
	}
	return false
}
