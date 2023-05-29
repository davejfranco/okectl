package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/resourcemanager"
)

// ResourceManager struct that implements the ResourceManager interface
type ResourceManager struct {
	Client resourcemanager.ResourceManagerClient
}

// ResourceManager interface
type RMClient interface {
	ListStacks(ctx context.Context, request resourcemanager.ListStacksRequest) (response resourcemanager.ListStacksResponse, err error)
	DeleteStack(ctx context.Context, req resourcemanager.DeleteStackRequest) (resourcemanager.DeleteStackResponse, error)
	CreateStack(ctx context.Context, req resourcemanager.CreateStackRequest) (resourcemanager.CreateStackResponse, error)
	CreateJob(ctx context.Context, req resourcemanager.CreateJobRequest) (resourcemanager.CreateJobResponse, error)
}

func (rm ResourceManager) ListStacks(ctx context.Context, request resourcemanager.ListStacksRequest) (resourcemanager.ListStacksResponse, error) {
	return rm.Client.ListStacks(ctx, request)
}

func (rm ResourceManager) CreateStack(ctx context.Context, req resourcemanager.CreateStackRequest) (resourcemanager.CreateStackResponse, error) {
	return rm.Client.CreateStack(ctx, req)
}

func (rm ResourceManager) DeleteStack(ctx context.Context, req resourcemanager.DeleteStackRequest) (resourcemanager.DeleteStackResponse, error) {
	return rm.Client.DeleteStack(ctx, req)
}

func (rm ResourceManager) CreateJob(ctx context.Context, req resourcemanager.CreateJobRequest) (resourcemanager.CreateJobResponse, error) {
	return rm.Client.CreateJob(ctx, req)
}
