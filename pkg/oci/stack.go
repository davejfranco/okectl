package oci

import (
	"context"
	"errors"
	"strings"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/resourcemanager"
)

const (
	tf_version string = "1.2.x"
)

type Stack struct {
	id               string
	client           RMClient
	Name             string
	CompartmentID    string
	TerraformVersion string
}

// NewStack creates a new stack
func NewStack(client RMClient) *Stack {
	return &Stack{
		client:           client,
		TerraformVersion: tf_version,
	}
}

// Get stack info by its name
func GetStack(name, compartmenID string, client RMClient) (*Stack, error) {
	req := resourcemanager.ListStacksRequest{
		CompartmentId: &compartmenID,
		DisplayName:   &name,
	}
	resp, err := client.ListStacks(context.Background(), req)
	if err != nil {
		return &Stack{}, err
	}
	if len(resp.Items) == 0 {
		return &Stack{}, errors.New("stack not found")
	}

	stack := NewStack(client)
	stack.id = *resp.Items[0].Id
	stack.Name = *resp.Items[0].DisplayName
	stack.CompartmentID = *resp.Items[0].CompartmentId
	stack.TerraformVersion = *resp.Items[0].TerraformVersion
	return stack, nil
}

// DeleteStack deletes a stack
func (s *Stack) Delete(stackID string) error {

	req := resourcemanager.DeleteStackRequest{
		StackId: &stackID,
	}

	_, err := s.client.DeleteStack(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

// CreateStack creates a resource manager stack
func (s *Stack) Create(zipfile string) (resourcemanager.CreateStackResponse, error) {

	req := resourcemanager.CreateStackRequest{
		CreateStackDetails: resourcemanager.CreateStackDetails{
			CompartmentId: &s.CompartmentID,
			DisplayName:   &s.Name,
			Description:   common.String("okectl generated stack for cluster"),
			ConfigSource: resourcemanager.CreateZipUploadConfigSourceDetails{
				ZipFileBase64Encoded: common.String(zipfile),
			},
			TerraformVersion: common.String(tf_version),
		},
	}

	resp, err := s.client.CreateStack(context.Background(), req)
	if err != nil {
		return resourcemanager.CreateStackResponse{}, err
	}
	//update the stack id
	s.id = *resp.Id

	return resp, nil
}

func (s *Stack) Job(action string) (string, error) {

	switch action {
	case "plan", "apply", "destroy":

		req := resourcemanager.CreateJobRequest{
			CreateJobDetails: resourcemanager.CreateJobDetails{
				StackId:   &s.id,
				Operation: resourcemanager.JobOperationEnum(strings.ToUpper(action)),
				ApplyJobPlanResolution: &resourcemanager.ApplyJobPlanResolution{
					IsAutoApproved: common.Bool(true), //for now this will be auto approved
				},
			},
		}

		resp, err := s.client.CreateJob(context.Background(), req)
		if err != nil {
			return "", err
		}
		return *resp.Id, nil //it returns the job id

	default:
		return "", errors.New("invalid action")
	}
}
