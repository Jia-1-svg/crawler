package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// ServiceInfo 服务信息结构体
type ServiceInfo struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Healthy bool
}

// DiscoverService 发现指定名称的健康服务
// serviceName: 服务名称
// 返回健康的服务实例列表
func (c *Client) DiscoverService(serviceName string) ([]ServiceInfo, error) {
	// 使用 Health().Service() 获取健康的服务实例
	services, _, err := c.apiClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	var result []ServiceInfo
	for _, entry := range services {
		// 检查服务是否通过健康检查
		isHealthy := true
		for _, check := range entry.Checks {
			if check.Status != api.HealthPassing {
				isHealthy = false
				break
			}
		}

		service := ServiceInfo{
			ID:      entry.Service.ID,
			Name:    entry.Service.Service,
			Address: entry.Service.Address,
			Port:    entry.Service.Port,
			Tags:    entry.Service.Tags,
			Healthy: isHealthy,
		}
		result = append(result, service)
	}

	return result, nil
}

// DiscoverServiceWithTag 根据标签发现服务
// serviceName: 服务名称
// tag: 服务标签
func (c *Client) DiscoverServiceWithTag(serviceName, tag string) ([]ServiceInfo, error) {
	services, _, err := c.apiClient.Health().Service(serviceName, tag, true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s with tag %s: %w", serviceName, tag, err)
	}

	var result []ServiceInfo
	for _, entry := range services {
		isHealthy := true
		for _, check := range entry.Checks {
			if check.Status != api.HealthPassing {
				isHealthy = false
				break
			}
		}

		service := ServiceInfo{
			ID:      entry.Service.ID,
			Name:    entry.Service.Service,
			Address: entry.Service.Address,
			Port:    entry.Service.Port,
			Tags:    entry.Service.Tags,
			Healthy: isHealthy,
		}
		result = append(result, service)
	}

	return result, nil
}

// GetAllServices 获取所有注册的服务
func (c *Client) GetAllServices() (map[string][]string, error) {
	services, err := c.apiClient.Agent().Services()
	if err != nil {
		return nil, fmt.Errorf("failed to get all services: %w", err)
	}

	result := make(map[string][]string)
	for id, service := range services {
		result[id] = service.Tags
	}

	return result, nil
}

// GetService 获取单个服务的详细信息
// serviceID: 服务唯一标识
func (c *Client) GetService(serviceID string) (*api.AgentService, error) {
	service, _, err := c.apiClient.Agent().Service(serviceID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get service %s: %w", serviceID, err)
	}

	return service, nil
}
