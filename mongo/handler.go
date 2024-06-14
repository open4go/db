package mongo

import (
	"context"
	"errors"
	"github.com/open4go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
)

var DBPool *DataBasePool

type DataBasePool struct {
	mu      sync.Mutex
	clients map[string]*mongo.Client
	Handler map[string]*mongo.Database
	ctx     context.Context
}

func NewDataBasePool() {
	DBPool = &DataBasePool{
		clients: make(map[string]*mongo.Client),
		Handler: make(map[string]*mongo.Database),
	}
	return
}

func (p *DataBasePool) GetClient(ctx context.Context, host string, name string) (*mongo.Client, error) {
	p.ctx = ctx
	p.mu.Lock()
	defer p.mu.Unlock()

	// 如果没有添加斜线，则为其加上
	// 便于后续步骤添加数据库名称
	if !strings.HasSuffix(host, "/") {
		host = host + "/"
	}

	client, ok := p.clients[host+name]
	if ok {
		// client has been init, no need to connect again
		return client, nil
	}
	uri := host + name
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Log(p.ctx).Fatal(err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Log(p.ctx).WithField("uri", uri).
			Fatal("Failed to ping MongoDB server: %v", err)
		// Handle error
	} else {
		log.Log(p.ctx).WithField("uri", uri).
			Info("MongoDB server is reachable")
		// MongoDB server is reachable, proceed with your logic
		p.clients[host+name] = client
		handler := client.Database(name)
		if handler != nil {
			p.Handler[name] = handler
		} else {
			log.Log(p.ctx).Fatal("handler is nil")
		}
	}
	return client, nil
}

// CloseAll 关闭所有 MongoDB 客户端连接
func (p *DataBasePool) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, client := range p.clients {
		if client != nil {
			_ = client.Disconnect(context.Background())
		}
	}
}

func (p *DataBasePool) GetHandler(name string) (*mongo.Database, error) {
	if handler, ok := p.Handler[name]; ok {
		return handler, nil
	}
	return nil, errors.New("no found any handler")
}
