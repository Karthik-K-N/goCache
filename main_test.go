package main

import (
	"testing"
	"time"

	"k8s.io/client-go/tools/cache"

	. "github.com/onsi/gomega"
)

func TestAddToCache(t *testing.T) {

	cacheStore := cache.NewTTLStore(cacheKeyFunc, cacheTTL)
	testKey := keyValue{
		key:   "testKey",
		value: "testValue",
	}
	err := addToCache(cacheStore, testKey)
	if err != nil {
		t.Fatalf("expecting error to be nil but got err %v", err)
	}
}

func TestFetchFromCache(t *testing.T) {

	defaultCacheStoreFunc := func() cache.Store {
		return cache.NewTTLStore(cacheKeyFunc, cacheTTL)
	}
	testKey := keyValue{
		key:   "testKey",
		value: "testValue",
	}

	testCases := []struct {
		name          string
		expectedError error
		expectedValue string
		keyName       string
		sleep         bool
		cacheStore    func() cache.Store
	}{
		{
			name:          "exists in cache",
			expectedValue: "testValue",
			keyName:       "testKey",
			cacheStore:    defaultCacheStoreFunc,
		},
		{
			name:          "not exists in cache",
			expectedValue: "",
			keyName:       "notTestKey",
			cacheStore:    defaultCacheStoreFunc,
		},
		{
			name:          "key in cache expired",
			expectedValue: "",
			keyName:       "testKey",
			sleep:         true,
			cacheStore: func() cache.Store {
				// set the cache timeout to 1 Millisecond so that by the time we fetch the key will be expired
				return cache.NewTTLStore(cacheKeyFunc, time.Millisecond)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gs := NewWithT(t)
			cacheStore := tc.cacheStore()
			err := addToCache(cacheStore, testKey)
			if err != nil {
				t.Fatalf("failed to add key to cache error %v", err)
			}
			if tc.sleep {
				time.Sleep(time.Second)
			}
			value, err := fetchFromCache(cacheStore, tc.keyName)
			if err != nil {
				if tc.expectedError != nil {
					gs.Expect(err).To(HaveOccurred())
					gs.Expect(err.Error()).To(Equal(tc.expectedError.Error()))
				} else {
					gs.Expect(err).ToNot(HaveOccurred())
				}
			}
			gs.Expect(value).To(Equal(tc.expectedValue))
		})
	}
}
