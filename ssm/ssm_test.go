package ssm

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type mockSSMClient struct {
	ssmiface.SSMAPI
}

var parameters = map[string]string{
	"/codacy/test/foo":   "bar",
	"/codacy/test/hello": "world",
}

func (m *mockSSMClient) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	tier := "Standard"
	version := int64(1)
	out := &ssm.PutParameterOutput{
		Tier:    &tier,
		Version: &version,
	}

	return out, nil
}

func (m *mockSSMClient) AddTagsToResource(input *ssm.AddTagsToResourceInput) (*ssm.AddTagsToResourceOutput, error) {
	out := &ssm.AddTagsToResourceOutput{}
	return out, nil
}

func (m *mockSSMClient) GetParametersByPath(input *ssm.GetParametersByPathInput) (*ssm.GetParametersByPathOutput, error) {
	var keyA = "/codacy/test/foo"
	var keyB = "/codacy/test/sneaky"
	var valA = "old bar"
	var valB = "will be deleted"
	out := &ssm.GetParametersByPathOutput{
		NextToken: nil,
		Parameters: []*ssm.Parameter{
			{
				ARN:              new(string),
				DataType:         new(string),
				LastModifiedDate: &time.Time{},
				Name:             &keyA,
				Selector:         new(string),
				SourceResult:     new(string),
				Type:             new(string),
				Value:            &valA,
				Version:          new(int64),
			},
			{
				ARN:              new(string),
				DataType:         new(string),
				LastModifiedDate: &time.Time{},
				Name:             &keyB,
				Selector:         new(string),
				SourceResult:     new(string),
				Type:             new(string),
				Value:            &valB,
				Version:          new(int64),
			},
		},
	}

	return out, nil
}

func (m *mockSSMClient) DeleteParameters(input *ssm.DeleteParametersInput) (*ssm.DeleteParametersOutput, error) {
	var s = "deleted test"
	out := &ssm.DeleteParametersOutput{
		DeletedParameters: []*string{
			&s,
		},
		InvalidParameters: nil,
	}

	return out, nil
}

func TestProcessParametersWithoutErrors(t *testing.T) {

	mockSvc := &mockSSMClient{}

	ProcessParameters(mockSvc, parameters, true)
}

func TestDeleteParametersWihoutErrors(t *testing.T) {
	mockSvc := &mockSSMClient{}
	CleanParameters(mockSvc, "/codacy/test/", true, parameters, nil)
}
