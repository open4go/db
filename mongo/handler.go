package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

var DBPool *DataBasePool

type DataBasePool struct {
	mu      sync.Mutex
	clients map[string]*mongo.Client
	handler map[string]*mongo.Database
}

func NewDataBasePool(ctx context.Context) *DataBasePool {
	// 创建一个新的数据库连接池
	return &DataBasePool{
		clients: make(map[string]*mongo.Client),
		handler: make(map[string]*mongo.Database),
	}
}

func (p *DataBasePool) GetClient(ctx context.Context, host string, name string) (*mongo.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := host + name
	if client, ok := p.clients[key]; ok {
		return client, nil
	}

	clientOptions := options.Client().ApplyURI(host)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			return nil, err
		} // disconnect if ping fails
		return nil, err
	}

	p.clients[key] = client
	p.handler[name] = client.Database(name)

	return client, nil
}

// CloseAll 关闭所有 MongoDB 客户端连接
func (p *DataBasePool) CloseAll(ctx context.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, client := range p.clients {
		if client != nil {
			err := client.Disconnect(ctx)
			if err != nil {
				return
			}
		}
	}
}

func (p *DataBasePool) GetHandler(name string) (*mongo.Database, error) {
	if handler, ok := p.handler[name]; ok {
		return handler, nil
	}
	return nil, errors.New("handler not found")
}
