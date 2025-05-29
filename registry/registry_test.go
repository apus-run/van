package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	// Test the NoopRegistry implementation
	reg := noopRegistry{}
	ctx := context.Background()
	assert.Nil(t, reg.Deregister(ctx, &ServiceInstance{}))
	assert.Nil(t, reg.Register(ctx, &ServiceInstance{}))
}
