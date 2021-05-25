package conf

// LogTransfer 全局配置
type LogTransferCfg struct {
	ESCfg    `ini:"es"`
	KafkaCfg `ini:"kafka"`
}

type ESCfg struct {
	Address string `ini:"address"`
}

type KafkaCfg struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
