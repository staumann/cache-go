package cache

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	cfg := Config{
		Enabled: true,
		TTL:     "2s",
		Logging: struct {
			Enabled bool
		}{Enabled: true},
	}
	Init(cfg)
	m.Run()
}

func TestGetFromCache(t *testing.T) {
	content, err := GetFromCache("testId_testLocale", func() (bytes []byte, e error) {
		return []byte("Testcontent"), nil
	})
	assert.Equal(t, "Testcontent", string(content))
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)
	c, _ := cacheMap.Load("testId_testLocale")

	assert.True(t, c.(entry).Expires.After(time.Now()))
	assert.Equal(t, "Testcontent", string(c.(entry).Data))
}

func TestGetFromCacheWithError(t *testing.T) {
	cacheMap = new(sync.Map)
	c, e := GetFromCache("testId2", func() (bytes []byte, e error) {
		return nil, errors.New("dam an error")
	})
	assert.Nil(t, c)
	assert.Equal(t, "dam an error", e.Error())
	v, _ := cacheMap.Load("testId2")
	assert.Nil(t, v)
}

func TestGetFromCacheFilled(t *testing.T) {
	cacheMap = new(sync.Map)

	cacheMap.Store("existingId_testingLocale", entry{
		Data:    []byte("existing content"),
		Expires: time.Now().Add(5 * time.Hour),
	})

	c, err := GetFromCache("existingId_testingLocale", func() (bytes []byte, e error) {
		return nil, nil
	})

	assert.Nil(t, err)
	assert.Equal(t, "existing content", string(c))
}

func TestCleanCache(t *testing.T) {
	cacheMap = new(sync.Map)
	cacheMap.Store("test", entry{
		Expires: time.Now(),
	})
	time.Sleep(1 * time.Second)
	keepCacheClean()
	time.Sleep(1 * time.Second)
	v, _ := cacheMap.Load("test")
	assert.Nil(t, v)
}

func TestGetTTL(t *testing.T) {
	ttl = nil
	config.TTL = "15 kps"

	d := getTTL()

	assert.Equal(t, float64(5), d.Minutes())
}

func BenchmarkGetFromCache(b *testing.B) {
	config.TTL = "15s"
	ttl = nil
	for n := 0; n < b.N; n++ {
		GetFromCache(fmt.Sprintf("entry_%d", n), func() (bytes []byte, e error) {
			return []byte(fmt.Sprintf("entry _ %d", n)), nil
		})
	}
}
