package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/resourcemanager"
	"github.com/stretchr/testify/mock"
)

type MockResourceManager struct {
	mock.Mock
}

func (m *MockResourceManager) ListStacks(ctx context.Context, request resourcemanager.ListStacksRequest) (resourcemanager.ListStacksResponse, error) {

	args := m.Called(ctx, request)
	return args.Get(0).(resourcemanager.ListStacksResponse), args.Error(1)
}
