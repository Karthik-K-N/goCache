package main

import (
	"k8s.io/klog/v2"
	"time"

	"k8s.io/client-go/tools/cache"
)

// cacheTTL is the duration of time to hold the key in cache
const cacheTTL = 20 * time.Second

// keyValue contains the fields to store the key and value
type keyValue struct {
	key   string
	value string
}

func main() {
	cacheStore := cache.NewTTLStore(cacheKeyFunc, cacheTTL)

	testKey := "myKey"
	key := keyValue{
		key:   testKey,
		value: "myValue",
	}

	klog.Infof("adding the key %v to cache", key)

	err := addToCache(cacheStore, key)
	if err != nil {
		klog.Fatalf("failed to add the key %v to cache error %v", key, err)
	}

	klog.Infof("fetching the value for key: %s from cache", testKey)

	value, err := fetchFromCache(cacheStore, "myKey")
	if err != nil {
		klog.Fatalf("failed to fetch value for key %S from cache error %v", testKey, err)
	}
	if value == "" {
		klog.Fatalf("the value for key %s is empty", testKey)
	}

	klog.Infof("successfully fetched the value for key %s from cache value: %s", testKey, value)

	klog.Infof("deleting the key %s from cache", testKey)
	err = deleteFromCache(cacheStore, key)
	if err != nil {
		klog.Fatalf("failed to delete key %s from cache error %v", testKey, err)
	}
}

func addToCache(cacheStore cache.Store, object keyValue) error {
	err := cacheStore.Add(object)
	if err != nil {
		klog.Errorf("failed to add key value to cache error", err)
		return err
	}
	return nil
}

func fetchFromCache(cacheStore cache.Store, key string) (string, error) {
	obj, exists, err := cacheStore.GetByKey(key)
	if err != nil {
		klog.Errorf("failed to add key value to cache error", err)
		return "", err
	}
	if !exists {
		klog.Errorf("object does not exist in the cache")
		return "", nil
	}
	return obj.(keyValue).value, nil
}

func deleteFromCache(cacheStore cache.Store, object keyValue) error {
	return cacheStore.Delete(object)
}

// cacheKeyFunc defines the key function required in TTLStore.
func cacheKeyFunc(obj interface{}) (string, error) {
	return obj.(keyValue).key, nil
}
