package registry

import (
	"errors"
	"sync"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

// --- Core cache structure ---
type ServiceInfo struct {
	Host       string
	Port       string
	LastUpdate time.Time
}

var (
	cache sync.Map // map[string]ServiceInfo
	ttl   = 60 * time.Second
)

// --- Public: get from cache or Master ---
func GetService(module string) (ServiceInfo, error) {
	if v, ok := cache.Load(module); ok {
		info := v.(ServiceInfo)
		if time.Since(info.LastUpdate) < ttl {
			return info, nil
		}
	}

	// Cache miss — try Master
	host := viper.GetString("microservices.masterservice.host")
	port := viper.GetString("microservices.masterservice.port")

	req := &pb.Request{
		Module: sf.StringPtr("masterservice"),
		IR: &pb.InternalRequest{
			Param:  sf.StringPtr("getmicroservicebypath"),
			Method: sf.StringPtr("GET"),
			Args:   map[string]*anypb.Any{},
		},
	}

	ans := sf.GRPCConnect(host, port, req)
	if ans["httpcode"] != nil {
		return ServiceInfo{}, errors.New("masterservice unavailable")
	}

	info := ServiceInfo{
		Host:       ans["host"].(string),
		Port:       ans["port"].(string),
		LastUpdate: time.Now(),
	}

	cache.Store(module, info)
	return info, nil
}

// --- Background refresh loop ---
func StartRefresher() {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			RefreshCache()
		}
	}()
}

// --- Refresh all entries ---
func RefreshCache() {
	cache.Range(func(key, val any) bool {
		mod := key.(string)
		info := val.(ServiceInfo)

		if time.Since(info.LastUpdate) > ttl {
			// Lazy revalidation
			newInfo, err := GetService(mod)
			if err == nil {
				cache.Store(mod, newInfo)
			}
		}
		return true
	})
}
