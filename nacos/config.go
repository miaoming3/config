package nacos

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"time"
)

type NacosInterface interface {
	GetConfig() (string, error)
}
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

func NewNacosConfig(opts ...Option) (NacosInterface, error) {
	cfg := &config{}

	for _, opt := range opts {
		opt(cfg)
	}
	if len(cfg.serverConfig) == 0 {
		return nil, errors.New("服务器配置不能为空")
	}
	return cfg, nil
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

func (c *config) Listen(handler func(content string, err error)) (func(), error) {
	client, err := c.CreateClient()
	if err != nil {
		return nil, err
	}
	err = client.ListenConfig(vo.ConfigParam{
		DataId: c.Params.DataId,
		Group:  c.Params.Group,
		OnChange: func(namespace, group, dataId, data string) {
			handler(data, nil)
		},
	})
	if err != nil {
		return nil, err
	}
	return func() {
		_ = client.CancelListenConfig(vo.ConfigParam{
			DataId: c.Params.DataId,
			Group:  c.Params.Group,
		})
	}, nil
}
