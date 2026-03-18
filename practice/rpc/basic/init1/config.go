package init1

import (
	"github.com/Jia-1-svg/crawler/practice/rpc/basic/config"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigFile("../../app.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	var nacos config.Nacos
	err := viper.UnmarshalKey("nacos", &nacos)
	if err != nil {
		panic(err)
	}
	// Nacos服务器地址
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacos.Addr,
			Port:   uint64(nacos.Port),
		},
	}
	// 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacos.Namespace, // 如果不需要命名空间，可以留空
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	configs, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacos.DataId,
		Group:  nacos.Group,
	})
	if err != nil {
		panic(err)
	}
	err = viper.ReadConfig(strings.NewReader(configs))
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config.Config)
	if err != nil {
		panic(err)
	}
	//fmt.Println(config.Config)
}
