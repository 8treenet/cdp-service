//package conf generated by 'freedom new-project github.com/8treenet/cdp-service'
package conf

import (
	"os"
	"runtime"
	"sync"

	"github.com/8treenet/freedom"
)

func init() {
	entryPoint()
}

// Get .
func Get() *Configuration {
	once.Do(func() {
		cfg = &Configuration{
			DB:     newDBConf(),
			App:    newAppConf(),
			Redis:  newRedisConf(),
			System: newSystemConf(),
		}
	})
	return cfg
}

var once sync.Once
var cfg *Configuration

// Configuration .
type Configuration struct {
	DB     *DBConf
	App    *freedom.Configuration
	Redis  *RedisConf
	System *SystemConf
}

// SystemConf .
type SystemConf struct {
	JobEnteringHouseSleep    int `toml:"job_entering_house_sleep"`
	JobEnteringHouseMaxCount int `toml:"job_entering_house_max_count"`
	JobTruncateHour          int `toml:"job_truncate_hour"`
	JobAnalysisHour          int `toml:"job_analysis_hour"`
	JobAnalysisRefreshHour   int `toml:"job_analysis_refresh_hour"`
	JobPersonaHour           int `toml:"job_persona_hour"`
	JobPersonaRefreshHour    int `toml:"job_persona_refresh_hour"`
}

// DBConf .
type DBConf struct {
	Addr            string `toml:"addr"`
	MaxOpenConns    int    `toml:"max_open_conns"`
	MaxIdleConns    int    `toml:"max_idle_conns"`
	ConnMaxLifeTime int    `toml:"conn_max_life_time"`
}

// RedisConf .
type RedisConf struct {
	Addr               string `toml:"addr"`
	Password           string `toml:"password"`
	DB                 int    `toml:"db"`
	MaxRetries         int    `toml:"max_retries"`
	PoolSize           int    `toml:"pool_size"`
	ReadTimeout        int    `toml:"read_timeout"`
	WriteTimeout       int    `toml:"write_timeout"`
	IdleTimeout        int    `toml:"idle_timeout"`
	IdleCheckFrequency int    `toml:"idle_check_frequency"`
	MaxConnAge         int    `toml:"max_conn_age"`
	PoolTimeout        int    `toml:"pool_timeout"`
}

func newSystemConf() *SystemConf {
	var result *SystemConf
	var appData struct {
		SystemConf SystemConf `toml:"system"`
	}
	if err := freedom.Configure(&appData, "app.toml"); err != nil {
		panic(err)
	}
	result = &appData.SystemConf

	if result.JobEnteringHouseSleep == 0 {
		result.JobEnteringHouseSleep = 8
	}
	if result.JobTruncateHour == 0 || result.JobTruncateHour > 23 || result.JobTruncateHour < 0 {
		result.JobTruncateHour = 4
	}
	if result.JobEnteringHouseMaxCount == 0 {
		result.JobEnteringHouseMaxCount = 1000
	}
	return result
}

func newAppConf() *freedom.Configuration {
	result := freedom.DefaultConfiguration()
	result.Other["listen_addr"] = ":8000"
	result.Other["service_name"] = "default"
	freedom.Configure(&result, "app.toml")
	return &result
}

func newDBConf() *DBConf {
	result := &DBConf{}
	if err := freedom.Configure(result, "db.toml"); err != nil {
		panic(err)
	}
	return result
}

func newRedisConf() *RedisConf {
	result := &RedisConf{
		MaxRetries:         0,
		PoolSize:           10 * runtime.NumCPU(),
		ReadTimeout:        3,
		WriteTimeout:       3,
		IdleTimeout:        300,
		IdleCheckFrequency: 60,
	}
	if err := freedom.Configure(result, "redis.toml"); err != nil {
		panic(err)
	}
	return result
}

func entryPoint() {
	// [./github.com/8treenet/cdp-service -c ./server/conf]
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] != "-c" {
			continue
		}
		if i+1 >= len(os.Args) {
			break
		}
		os.Setenv(freedom.ProfileENV, os.Args[i+1])
	}
}
