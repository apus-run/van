package registry

import (
	"context"
	"fmt"
	"sort"
)

// Registry for discovery
type Registry interface {
	// Register service to registry
	Register(ctx context.Context, ins *ServiceInstance) error
	// Deregister service from registry
	Deregister(ctx context.Context, ins *ServiceInstance) error
	// GetService return the service in memory according to the service name.
	GetService(ctx context.Context, name string) ([]*ServiceInstance, error)
	// Watch creates a watcher according to the service name.
	Watch(ctx context.Context, name string) (Watcher, error)
}

// Watcher is service watcher.
type Watcher interface {
	// Next returns services in the following two cases:
	// 1.the first time to watch and the service instance list is not empty.
	// 2.any service instance changes found.
	// if the above two conditions are not met, it will block until context deadline exceeded or canceled
	Next() ([]*ServiceInstance, error)
	// Stop close the watcher.
	Stop() error
}

// ServiceInstance is an instance of a service in a discovery system.
type ServiceInstance struct {
	// ID is the unique instance ID as registered.
	ID string `json:"id"`
	// Name is the service name as registered.
	Name string `json:"name"`
	// Version is the version of the compiled.
	Version string `json:"version"`
	// Metadata is the kv pair metadata associated with the service instance.
	Metadata map[string]string `json:"metadata"`
	// Endpoints are endpoint addresses of the service instance.
	// schema:
	//   http://127.0.0.1:8000?isSecure=false
	//   grpc://127.0.0.1:9000?isSecure=false
	Endpoints []string `json:"endpoints"`
}

func (i *ServiceInstance) String() string {
	return fmt.Sprintf("%s-%s", i.Name, i.ID)
}

// Equal returns whether i and o are equivalent.
func (i *ServiceInstance) Equal(o interface{}) bool {
	if i == nil && o == nil {
		return true
	}

	if i == nil || o == nil {
		return false
	}

	t, ok := o.(*ServiceInstance)
	if !ok {
		return false
	}

	if len(i.Endpoints) != len(t.Endpoints) {
		return false
	}

	sort.Strings(i.Endpoints)
	sort.Strings(t.Endpoints)
	for j := 0; j < len(i.Endpoints); j++ {
		if i.Endpoints[j] != t.Endpoints[j] {
			return false
		}
	}

	if len(i.Metadata) != len(t.Metadata) {
		return false
	}

	for k, v := range i.Metadata {
		if v != t.Metadata[k] {
			return false
		}
	}

	return i.ID == t.ID && i.Name == t.Name && i.Version == t.Version
}

var NoopRegistry Registry = &noopRegistry{}

// NoopRegistry is a no-op registry implementation.
type noopRegistry struct{}

func (n *noopRegistry) Register(ctx context.Context, ins *ServiceInstance) error {
	return nil
}
func (n *noopRegistry) Deregister(ctx context.Context, ins *ServiceInstance) error {
	return nil
}
func (n *noopRegistry) GetService(ctx context.Context, name string) ([]*ServiceInstance, error) {
	return nil, nil
}
func (n *noopRegistry) Watch(ctx context.Context, name string) (Watcher, error) {
	return nil, nil
}
