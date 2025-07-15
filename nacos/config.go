package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"time"
)

type Option func(*config)
type config struct {
	Address      string
	Params       vo.ConfigParam
	Port         uint64
	Namespace    string
	Temp         string
	serverConfig []constant.ServerConfig
	clientConfig *constant.ClientConfig
}

func WithServerConfig(address string, port uint64) Option {
	return func(c *config) {
		c.serverConfig = append(c.serverConfig, constant.ServerConfig{
			IpAddr: address,
			Port:   port,
		})
	}
}

func WithNamespace(namespace string) Option {
	return func(c *config) {
		c.Namespace = namespace
	}
}

func WithTemp(temp string) Option {
	return func(c *config) {
		if temp == "" {
			temp = "tmp/log/"
		}
		c.Temp = temp
	}
}

func WithParams(dataId, group string) Option {
	return func(c *config) {
		if group == "" {
			group = "DEFAULT_GROUP"
		}
		c.Params.Group = group
		c.Params.DataId = dataId
	}
}

func WithClientConfig() Option {
	return func(c *config) {
		c.clientConfig = &constant.ClientConfig{
			NamespaceId:         c.Namespace,
			TimeoutMs:           uint64(30 * time.Millisecond),
			NotLoadCacheAtStart: true,
			LogLevel:            "warn",
			LogDir:              c.Temp,
		}
	}
}

func (c *config) CreateClient() (config_client.IConfigClient, error) {
	return clients.NewConfigClient(
		vo.NacosClientParam{
			ServerConfigs: c.serverConfig,
			ClientConfig:  c.clientConfig,
		},
	)

}

func (c *config) GetConfig() (string, error) {
	client, err := c.CreateClient()
	if err != nil {
		return "", err
	}
	return client.GetConfig(c.Params)
}
