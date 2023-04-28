package oci

import (
	"context"
	"testing"

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

func (m *MockResourceManager) DeleteStack(ctx context.Context, req resourcemanager.DeleteStackRequest) (resourcemanager.DeleteStackResponse, error) {

	args := m.Called(ctx, req)
	return args.Get(0).(resourcemanager.DeleteStackResponse), args.Error(1)
}

func (m *MockResourceManager) CreateStack(ctx context.Context, req resourcemanager.CreateStackRequest) (resourcemanager.CreateStackResponse, error) {

	args := m.Called(ctx, req)
	return args.Get(0).(resourcemanager.CreateStackResponse), args.Error(1)
}

func TestFindStack(t *testing.T) {
	client := new(MockResourceManager)
	stack := Stack{
		Client:        client,
		Name:          "test",
		CompartmentID: "test",
	}
	// Test stack not found
	client.On("ListStacks", mock.Anything, mock.Anything).Return(resourcemanager.ListStacksResponse{}, nil)
	_, err := stack.findStack()
	if err != nil {
		t.Errorf("Error should be nil, got %v", err)
	}

	// Test stack found
	client.On("ListStacks", mock.Anything, mock.Anything).Return(resourcemanager.ListStacksResponse{
		Items: []resourcemanager.StackSummary{
			{
				Id: &[]string{"test"}[0],
			},
		},
	}, nil)
	_, err = stack.findStack()
	if err != nil {
		t.Errorf("Error should be nil, got %v", err)
	}

}
