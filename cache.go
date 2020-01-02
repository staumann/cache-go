package cache

import (
	"log"
	"sync"
	"time"
)

var (
	cacheMap *sync.Map
	ttl      *time.Duration
	config   Config
	initDone bool
)

func init() {
	cacheMap = new(sync.Map)
	ttl = nil
	initDone = false
}

//ShutDown use this method to shutdown the cache. To restart it you need to call the init method with an configuration struct
func ShutDown() {
	initDone = false
	config = Config{
		Enabled: false,
		Logging: struct {
			Enabled bool
		}{
			false,
		},
	}
	cacheMap = new(sync.Map)
}

//Init is the method to initialize the caching mechanism
func Init(cfg Config) {
	if !initDone {
		initDone = true
		config = cfg
		if config.Enabled {
			log.Print("Caching is enabled")
			go func() {
				for {
					time.Sleep(getTTL())
					keepCacheClean()
				}
			}()
		} else {
			log.Print("Caching is disabled")
		}
	} else {
		log.Printf("caution the caching mechanism was already initialized")
	}
}

//GetFromCache is the method to retrieve a cached entry.
//entryID the id of the cacheEntry
//fetcher the function to fill the cache if the entry is not present
//if the cache is not enabled the method returns the result of fetcher
func GetFromCache(entryID string, fetcher func() ([]byte, error)) ([]byte, error) {
	if config.Enabled {
		var returnValue []byte
		var err error
		cacheEntry, ok := cacheMap.Load(entryID)
		if ok {
			returnValue = cacheEntry.(entry).Data
		} else {
			returnValue, err = fillCache(entryID, fetcher)
		}

		return returnValue, err
	}
	return fetcher()
}

func fillCache(id string, fetcher func() ([]byte, error)) ([]byte, error) {
	data, err := fetcher()
	if err == nil {
		cacheMap.Store(id,
			entry{
				Data:    data,
				Expires: time.Now().Add(getTTL()),
			})
	} else {
		log.Printf("error filling cache: %s", err.Error())
	}
	return data, err
}

func keepCacheClean() {
	t := time.Now()
	counter := 0
	cacheMap.Range(func(key, value interface{}) bool {
		if value.(entry).Expires.Before(t) {
			cacheMap.Delete(key)
		}
		counter = counter + 1
		return true
	})

	if config.Logging.Enabled {
		log.Printf("current cache size: %d", counter)
	}
}

func getTTL() time.Duration {
	if ttl == nil {
		d, e := time.ParseDuration(config.TTL)
		if e != nil {
			log.Print("error parsing cache ttl duration. Making fallback to 5 minutes")
			d = 5 * time.Minute
		}
		ttl = &d
	}

	return *ttl
}
