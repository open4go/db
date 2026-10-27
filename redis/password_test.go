package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/spf13/viper"
)

func TestGetClientWithoutPasswordStillWorks(t *testing.T) {
	mr := miniredis.RunT(t)
	NewDataBasePool()
	c, err := DBPool.GetClient(context.Background(), mr.Addr(), "plain", 0)
	if err != nil || c == nil {
		t.Fatalf("legacy GetClient must keep working without auth: %v", err)
	}
	DBPool.CloseAll()
}

func TestGetClientWithPassword(t *testing.T) {
	mr := miniredis.RunT(t)
	mr.RequireAuth("s3cret")
	NewDataBasePool()
	if _, err := DBPool.GetClient(context.Background(), mr.Addr(), "noauth", 0); err == nil {
		t.Fatal("no password against an auth-enabled server must fail")
	}
	if _, err := DBPool.GetClientWithPassword(context.Background(), mr.Addr(), "bad", 0, "wrong"); err == nil {
		t.Fatal("wrong password must fail")
	}
	c, err := DBPool.GetClientWithPassword(context.Background(), mr.Addr(), "ok", 0, "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Set(context.Background(), "k", "v", 0).Err(); err != nil {
		t.Fatal(err)
	}
	h, err := DBPool.GetHandler("ok")
	if err != nil || h != c {
		t.Fatalf("handler lookup: %v", err)
	}
	DBPool.CloseAll()
}

func TestInitReadsPasswordFromConfig(t *testing.T) {
	mr := miniredis.RunT(t)
	mr.RequireAuth("s3cret")
	viper.Set("db.redis", []map[string]interface{}{
		{"host": mr.Addr(), "name": "cache", "db": 0, "password": "s3cret"},
	})
	t.Cleanup(func() { viper.Set("db.redis", nil) })
	Init(context.Background())
	h, err := DBPool.GetHandler("cache")
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	DBPool.CloseAll()
}
