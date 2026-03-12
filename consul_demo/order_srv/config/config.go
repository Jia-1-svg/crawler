package config

import (
	"encoding/json"
	"os"
	"strconv"
)

// DBConfig 数据库配置
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// ConsulConfig Consul 配置
type ConsulConfig struct {
	Address string `json:"address"`
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	ServiceName string `json:"service_name"`
	ServiceID   string `json:"service_id"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
}

// Config 总配置
type Config struct {
	Database DBConfig      `json:"database"`
	Consul   ConsulConfig  `json:"consul"`
	Service  ServiceConfig `json:"service"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	cfg := &Config{
		Database: DBConfig{
			Host:     "115.190.88.2",
			Port:     3306,
			User:     "root",
			Password: "4ay1nkal3u8ed77y",
			DBName:   "zg2lx1",
		},
		Consul: ConsulConfig{
			Address: "115.190.88.2:8500",
		},
		Service: ServiceConfig{
			ServiceName: "order-service",
			ServiceID:   "order-service-1",
			Host:        "127.0.0.1",
			Port:        50051,
		},
	}

	// 从环境变量覆盖服务配置
	if serviceID := os.Getenv("SERVICE_ID"); serviceID != "" {
		cfg.Service.ServiceID = serviceID
	}
	if servicePort := os.Getenv("SERVICE_PORT"); servicePort != "" {
		if port, err := strconv.Atoi(servicePort); err == nil {
			cfg.Service.Port = port
		}
	}

	return cfg
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, nil
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
