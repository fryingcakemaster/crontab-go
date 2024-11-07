package worker

type ConfigData struct {
	Title   string      `toml:"title"`
	Version string      `toml:"version"`
	Etcd    EtcdConfig  `toml:"etcd"`
	MongoDB MongoConfig `toml:"mongodb"`
	Log     LogConfig   `toml:"log"`
}

type EtcdConfig struct {
	EndPoints   []string `toml:"endpoints"`
	DialTimeout int      `toml:"timeout"`
}

type MongoConfig struct {
	Uri            string `toml:"uri"`
	ConnectTimeout int    `toml:"timeout"`
}

type LogConfig struct {
	Console  bool   `toml:"console"`  // 是否打印到控制台
	Level    int    `toml:"level"`    // 日志等级
	FileName string `toml:"filename"` // 保存日志的文件名
	Daily    bool   `toml:"daily"`    //是否按天rotate
	MaxDays  int    `toml:"max_days"` // 保存日志天数
}
