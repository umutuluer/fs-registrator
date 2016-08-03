package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type KvBackend interface {
	Read(key string, recursive bool) (*string, error)
	Write(key string, value string, ttl string) error
}

// Credit: http://matthewbrown.io/2016/01/23/factory-pattern-in-golang/

func init() {
	RegisterKvBackend("etcd", NewKvBackendEtcd)
	// Add new backends here as they become available.
}

type KvBackendFactory func(conf map[string]string) (KvBackend, error)

var kvBackendFactories = make(map[string]KvBackendFactory)

func RegisterKvBackend(name string, factory KvBackendFactory) {
	if factory == nil {
		log.Fatalf("K/V backend factory '%s' does not exist.", name)
	}
	_, registered := kvBackendFactories[name]
	if registered {
		log.Printf("K/V backend factory '%s' already registered. Ignoring.", name)
	}
	kvBackendFactories[name] = factory
}

func CreateKvBackend(conf map[string]string) (KvBackend, error) {
	if _, ok := conf["backend"]; ok == false {
		return nil, errors.New("'backend' key does not exist in conf.")
	}

	kvBackendFactory, ok2 := kvBackendFactories[conf["backend"]]

	if ok2 == false {
		// Factory has not been registered.
		// Make a list of all available datastore factories for logging.
		availableKvBackends := make([]string, len(kvBackendFactories))
		for k, _ := range kvBackendFactories {
			availableKvBackends = append(availableKvBackends, k)
		}
		fmt.Printf("backends: %+v\n", availableKvBackends)
		return nil, errors.New(fmt.Sprintf("Invalid K/V Backend Name. Must be one of: %s", strings.Join(availableKvBackends, ", ")))
	}

	// Run the factory with the configuration.
	return kvBackendFactory(conf)
}
