package core

type Config struct {
	CanalAddr      string
	CanalPort      int
	CanalUser      string
	CanalPassword  string
	CanalDest      string
	CanalBatchSize int
	ESAddrs        []string
	SyncTables     map[string]string // key为MySQL表名，value为ES索引名
	MySQLDatabase  string
}

var DefaultConfig = Config{}

func LoadConfig() {
	cfg := Config{
		CanalAddr:      "canal-server",
		CanalPort:      11111,
		CanalUser:      "canal",
		CanalPassword:  "canal",
		CanalDest:      "example",
		CanalBatchSize: 100,
		ESAddrs:        []string{"http://es.xyq777.com"},
		MySQLDatabase:  "tag",
	}
	DefaultConfig = cfg

}
