package v6_config

type AppConfig struct {
	KafkaConf `ini:"kafka"`
	EtcdConf  `ini:"etcd"`
}

type KafkaConf struct {
	Addr    string `ini:"address"`
	Topic   string `ini:"topic"`
	MaxSize int    `ini:"max_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"collect_log_key"`
}

// ---  unused ---
type TaillogConf struct {
	Filename string `ini:"filename"`
}
