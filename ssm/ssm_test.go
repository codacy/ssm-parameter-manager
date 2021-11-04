package ssm

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

type mockSSMClient struct {
	ssmiface.SSMAPI
}

var parameters = map[string]string{
	"foo":   "bar",
	"hello": "world",
}

var data = map[string]map[string]string{
	"test": parameters,
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

func TestProcessParametersWithoutErrors(t *testing.T) {

	mockSvc := &mockSSMClient{}

	ProcessParameters(mockSvc, data, "test", true)
}
