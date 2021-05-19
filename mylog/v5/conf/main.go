package v5_config

type AppConfig struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"tailog"`
}

type KafkaConf struct {
	Addr  string `ini:"address"`
	Topic string `ini:"topic"`
}

type TaillogConf struct {
	Filename string `ini:"filename"`
}
