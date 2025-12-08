package watchservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// ServiceInstance 表示一个服务实例
type ServiceInstance struct {
	ServiceName string            `json:"serviceName"`
	InstanceID  string           `json:"instanceId"`
	Host        string           `json:"host"`
	Port        int             `json:"port"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	LastHeartbeat time.Time      `json:"lastHeartbeat"`
}

// ServiceRegistry 服务注册中心
type ServiceRegistry struct {
	redisClient *redis.Client
	ttl         time.Duration
}

// NewServiceRegistry 创建新的服务注册中心
func NewServiceRegistry(redisClient *redis.Client, ttl time.Duration) *ServiceRegistry {
	return &ServiceRegistry{
		redisClient: redisClient,
		ttl:         ttl,
	}
}

// Register 注册服务实例
func (sr *ServiceRegistry) Register(ctx context.Context, instance *ServiceInstance) error {
	instance.LastHeartbeat = time.Now()
	data, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("marshal instance failed: %w", err)
	}

	key := fmt.Sprintf("service:%s:%s", instance.ServiceName, instance.InstanceID)
	err = sr.redisClient.Set(ctx, key, string(data), sr.ttl).Err()
	if err != nil {
		return fmt.Errorf("register service failed: %w", err)
	}
	return nil
}

// Deregister 注销服务实例
func (sr *ServiceRegistry) Deregister(ctx context.Context, serviceName, instanceID string) error {
	key := fmt.Sprintf("service:%s:%s", serviceName, instanceID)
	err := sr.redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("deregister service failed: %w", err)
	}
	return nil
}

// GetService 获取指定服务的所有实例
func (sr *ServiceRegistry) GetService(ctx context.Context, serviceName string) ([]*ServiceInstance, error) {
	pattern := fmt.Sprintf("service:%s:*", serviceName)
	keys, err := sr.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("get service keys failed: %w", err)
	}

	instances := make([]*ServiceInstance, 0, len(keys))
	for _, key := range keys {
		data, err := sr.redisClient.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, fmt.Errorf("get service instance failed: %w", err)
		}

		var instance ServiceInstance
		if err := json.Unmarshal([]byte(data), &instance); err != nil {
			return nil, fmt.Errorf("unmarshal instance failed: %w", err)
		}
		instances = append(instances, &instance)
	}

	return instances, nil
}

// Heartbeat 服务心跳更新
func (sr *ServiceRegistry) Heartbeat(ctx context.Context, serviceName, instanceID string) error {
	key := fmt.Sprintf("service:%s:%s", serviceName, instanceID)
	
	// 获取现有实例数据
	data, err := sr.redisClient.Get(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("get instance for heartbeat failed: %w", err)
	}

	var instance ServiceInstance
	if err := json.Unmarshal([]byte(data), &instance); err != nil {
		return fmt.Errorf("unmarshal instance for heartbeat failed: %w", err)
	}

	// 更新心跳时间
	instance.LastHeartbeat = time.Now()
	updatedData, err := json.Marshal(instance)
	if err != nil {
		return fmt.Errorf("marshal updated instance failed: %w", err)
	}

	// 更新 Redis 并刷新 TTL
	err = sr.redisClient.Set(ctx, key, string(updatedData), sr.ttl).Err()
	if err != nil {
		return fmt.Errorf("update heartbeat failed: %w", err)
	}

	return nil
}