package ssm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

func ProcessParameters(svc ssmiface.SSMAPI, parameters map[string]string, verbose bool) (int, error) {
	for k, v := range parameters {
		if verbose {
			fmt.Printf("**PUTTED** \"%s\" - \"%s\"\n", k, v)
		}

		_, err := putParameter(svc, k, v, "String", true)

		if err != nil {
			return -1, err
		}

		_, err = tagParameter(svc, k, "Parameter", createTags())

		if err != nil {
			return -1, err
		}
	}

	return len(parameters), nil
}

func CleanParameters(svc ssmiface.SSMAPI, path string, verbose bool, plainParameters map[string]string, encryptedParameters map[string]string) (int, error) {
	var allParams = make(map[string]string)
	var result *ssm.GetParametersByPathOutput
	var nextToken *string
	var err error

	for result == nil || result.NextToken != nil {
		result, err = getParametersByPrefix(svc, nextToken, 10, path, true, false)

		if err != nil {
			return -1, err
		}

		for _, v := range result.Parameters {
			allParams[*v.Name] = *v.Value
		}

		nextToken = result.NextToken
	}

	// Remove parameters contained in the config files to avoid deleting them
	for k := range plainParameters {
		delete(allParams, k)
	}

	for k := range encryptedParameters {
		delete(allParams, k)
	}

	if len(allParams) == 0 {
		fmt.Printf("No parameters to delete.\n")
		return 0, nil
	}

	fmt.Printf("Found %d parameters not contained in the ssm configuration files. Deleting...\n", len(allParams))
	var paramsToDelete []*string

	for k, v := range allParams {
		if !strings.HasSuffix(path, "/") || !strings.HasPrefix(k, path) {
			return -1, errors.New("prefix doesn't end with \"/\" or would delete a parameter that doesn't start with the specified prefix")
		}

		if verbose {
			fmt.Printf("**DELETING**  \"%s\" - \"%s\" \n", k, v)
		}

		var s = k
		paramsToDelete = append(paramsToDelete, &s)
	}

	results, err := deleteParameters(svc, paramsToDelete)

	if err != nil {
		return -1, err
	}

	if len(paramsToDelete) != len(results.DeletedParameters) {
		return -1, fmt.Errorf("expected to delete %d parameters but deleted %d instead", len(paramsToDelete), len(results.DeletedParameters))
	}

	return len(results.DeletedParameters), nil
}

func createTags() []*ssm.Tag {
	tagKey := "ssm-managed"
	tagValue := "true"

	t := ssm.Tag{
		Key:   &tagKey,
		Value: &tagValue,
	}

	var tags []*ssm.Tag

	tags = append(tags, &t)

	return tags
}

func putParameter(svc ssmiface.SSMAPI, name string, value string, paramType string, overwrite bool) (*ssm.PutParameterOutput, error) {
	results, err := svc.PutParameter(&ssm.PutParameterInput{
		Name:      &name,
		Value:     &value,
		Type:      &paramType,
		Overwrite: &overwrite,
	})

	return results, err
}

func tagParameter(svc ssmiface.SSMAPI, resourceId string, resourceType string, tags []*ssm.Tag) (*ssm.AddTagsToResourceOutput, error) {
	results, err := svc.AddTagsToResource(&ssm.AddTagsToResourceInput{
		ResourceId:   &resourceId,
		ResourceType: &resourceType,
		Tags:         tags,
	})

	return results, err
}

func getParametersByPrefix(svc ssmiface.SSMAPI, nextToken *string, maxResults int64, path string, recursive bool, decrypt bool) (*ssm.GetParametersByPathOutput, error) {
	results, err := svc.GetParametersByPath(&ssm.GetParametersByPathInput{
		MaxResults:       &maxResults,
		NextToken:        nextToken,
		ParameterFilters: nil,
		Path:             &path,
		Recursive:        &recursive,
		WithDecryption:   &decrypt,
	})

	return results, err
}

func deleteParameters(svc ssmiface.SSMAPI, names []*string) (*ssm.DeleteParametersOutput, error) {
	results, err := svc.DeleteParameters(&ssm.DeleteParametersInput{
		Names: names,
	})

	return results, err
}
