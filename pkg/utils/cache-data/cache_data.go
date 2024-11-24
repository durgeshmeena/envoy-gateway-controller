package cache

import (
	"log"
	"sync"
	"time"
)

// data fields to store data
type PolicyData struct {

}

// in-memory cache to store PolicyData
type Cache struct {
	mu      sync.RWMutex
	data 	map[string]PolicyData
}

// retrieves data from the cache
func (c *Cache) Get(key string) (PolicyData, bool) {
 	c.mu.RLock()
    defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

// store data in the cache
func (c *Cache) Set(key string, val PolicyData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = val
}


// SyncData periodically fetches data from the database and updates the cache
func (c *Cache) SyncData() {
	for {
		//  fetch data from the database
		data, err := fetchDataFromDatabase()
		if err != nil {
			log.Println("Error fetching data: ", err)
		    time.Sleep(1 * time.Minute)
			continue
		}

		// update the cache
		c.mu.Lock()
		c.data = data
		c.mu.Unlock()
		
		// sync every 10 minutes
		time.Sleep(10 * time.Minute)
	}
}


// fetchDataFromDatabase fetches data from the database
func fetchDataFromDatabase() (map[string]PolicyData, error) {
    // Implement database fetching logic here
    return map[string]PolicyData{}, nil
}