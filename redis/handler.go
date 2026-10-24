package redis

import (
	"context"
	"errors"
	"github.com/open4go/log"
	v9 "github.com/redis/go-redis/v9"
	"strings"
	"sync"
)

var DBPool *DataBasePool

type DataBasePool struct {
	mu      sync.Mutex
	clients map[string]*v9.Client
}

func NewDataBasePool() {
	DBPool = &DataBasePool{
		clients: make(map[string]*v9.Client),
	}
}

// GetClient
// host: localhost:2379
// name: auth, middle, cache, etc
// db: 0, 2, 3, 4
func (p *DataBasePool) GetClient(ctx context.Context,
	host string, name string, db int) (*v9.Client, error) {
	return p.GetClientWithPassword(ctx, host, name, db, "")
}

// GetClientWithPassword is GetClient with an optional AUTH password.
func (p *DataBasePool) GetClientWithPassword(ctx context.Context,
	host string, name string, db int, password string) (*v9.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果没有添加斜线，则为其加上
	// 便于后续步骤添加数据库名称
	if !strings.HasSuffix(host, "/") {
		host = strings.Trim(host, "/")
	}

	if client, ok := p.clients[name]; ok {
		return client, nil
	}
	uri := host
	client := v9.NewClient(&v9.Options{
		Addr:       uri,
		Password:   password,
		DB:         db,
		PoolSize:   100,
		MaxRetries: 3,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Log(ctx).WithField("uri", uri).WithField("name", name).Error(err)
		_ = client.Close()
		return nil, err
	}
	log.Log(ctx).WithField("uri", uri).WithField("name", name).Info("Redis server is reachable")
	p.clients[name] = client
	return client, nil
}

// CloseAll 关闭所有 MongoDB 客户端连接
func (p *DataBasePool) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, client := range p.clients {
		if client != nil {
			_ = client.Close()
		}
	}
}

func (p *DataBasePool) GetHandler(name string) (*v9.Client, error) {
	if handler, ok := p.clients[name]; ok {
		return handler, nil
	}
	return nil, errors.New("no found any handler")
}
