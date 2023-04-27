package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

const (
	tf_version string = "1.2.x"
)

type stackInterface interface {
	ListStacks(ctx context.Context, request resourcemanager.ListStacksRequest) (response resourcemanager.ListStacksResponse, err error)
	DeleteStack(ctx context.Context, req resourcemanager.DeleteStackRequest) (resourcemanager.DeleteStackResponse, error)
	CreateStack(ctx context.Context, req resourcemanager.CreateStackRequest) (resourcemanager.CreateStackResponse, error)
}

type Stack struct {
	Client           stackInterface
	Name             string
	CompartmentID    string
	Zipfile          string
	TerraformVersion string
}

// Get stack info by its name
func (s *Stack) Get() (resourcemanager.StackSummary, error) {
	req := resourcemanager.ListStacksRequest{
		CompartmentId: &s.CompartmentID,
		DisplayName:   &s.Name,
	}
	resp, err := s.Client.ListStacks(context.Background(), req)
	if err != nil {
		return resourcemanager.StackSummary{}, err
	}
	if len(resp.Items) == 0 {
		return resourcemanager.StackSummary{}, nil
	}
	return resp.Items[0], nil
}

// DeleteStack deletes a stack
func (s *Stack) Delete(stackID string) error {

	getStack, err := s.Get()
	if err != nil {
		return err
	
	req := resourcemanager.DeleteStackRequest{
		StackId: &stackID,
	}
	_, err := s.Client.DeleteStack(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}

// CreateStack creates a resource manager stack
func (s *Stack) Create() (resourcemanager.CreateStackResponse, error) {

	req := resourcemanager.CreateStackRequest{
		CreateStackDetails: resourcemanager.CreateStackDetails{
			CompartmentId: &s.CompartmentID,
			DisplayName:   &s.Name,
			Description:   common.String("okectl generated stack for cluster"),
			ConfigSource: resourcemanager.CreateZipUploadConfigSourceDetails{
				ZipFileBase64Encoded: common.String(s.Zipfile),
			},
			TerraformVersion: common.String(tf_version),
		},
	}

	resp, err := s.Client.CreateStack(context.Background(), req)
	if err != nil {
		return resourcemanager.CreateStackResponse{}, err
	}
	return resp, nil
}