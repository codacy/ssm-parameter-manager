package ssm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

func validParameterType(parameterType string) bool {
	for _, t := range ssm.ParameterType_Values() {
		if parameterType == t {
			return true
		}
	}

	return false
}

func parseParameter(key string, parameter interface{}) (*string, *string, error) {
	var parameterType, parameterValue string
	switch p := parameter.(type) {
	case string:
		parameterType = "String"
		parameterValue = p
	case map[string]interface{}:
		var ok bool
		if parameterType, ok = p["type"].(string); !ok {
			return nil, nil, fmt.Errorf("key [%s] doesnt have a defined type", key)
		}
		if !validParameterType(parameterType) {
			return nil, nil, fmt.Errorf("invalid parameter type [%s] for key [%s]", parameterType, key)
		}

		if parameterValue, ok = p["value"].(string); !ok {
			return nil, nil, fmt.Errorf("key [%s] doesnt have a defined value", key)
		}

		if parameterValue == "" || (parameterType == ssm.ParameterTypeStringList && !strings.Contains(parameterValue, ",")) {
			return nil, nil, fmt.Errorf("invalid value [%s] for key [%s], it needs to be a valid list", parameterValue, key)
		}

	default:
		return nil, nil, errors.New("unknown parameter definition")
	}

	return &parameterType, &parameterValue, nil
}

// ProcessParameters takes a map of parameters and pushes them to the parameter store of the configured AWS environemnt
func ProcessParameters(svc ssmiface.SSMAPI, parameters map[string]interface{}, verbose bool) error {
	for k, v := range parameters {
		var parameterType, parameterValue, err = parseParameter(k, v)

		if err != nil {
			return err
		}

		if verbose {
			fmt.Printf("**PUSHED** \"%s\" - \"%s\"\n", k, *parameterValue)
		}

		_, err = putParameter(svc, k, *parameterValue, *parameterType, true)

		if err != nil {
			return err
		}

		_, err = tagParameter(svc, k, "Parameter", createTags())

		if err != nil {
			return err
		}
	}

	return nil
}

// CleanParameters deletes SSM parameters for a given path, if they're not present in the specified maps
func CleanParameters(svc ssmiface.SSMAPI, path string, verbose bool, plainParameters map[string]interface{}, encryptedParameters map[string]interface{}) error {
	var allParams = make(map[string]string)
	var result *ssm.GetParametersByPathOutput
	var nextToken *string
	var err error

	for result == nil || result.NextToken != nil {
		result, err = getParametersByPrefix(svc, nextToken, 10, path, true, false)

		if err != nil {
			return err
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
		return nil
	}

	fmt.Printf("Found %d parameters not contained in the ssm configuration files. Deleting...\n", len(allParams))
	var paramsToDelete []*string

	for k, v := range allParams {
		if !strings.HasSuffix(path, "/") || !strings.HasPrefix(k, path) {
			return errors.New("prefix doesn't end with \"/\" or would delete a parameter that doesn't start with the specified prefix")
		}

		if verbose {
			fmt.Printf("**DELETING**  \"%s\" - \"%s\" \n", k, v)
		}

		var s = k
		paramsToDelete = append(paramsToDelete, &s)
	}

	results, err := deleteParameters(svc, paramsToDelete)

	if err != nil {
		return err
	}

	if len(paramsToDelete) != len(results.DeletedParameters) {
		return fmt.Errorf("expected to delete %d parameters but deleted %d instead", len(paramsToDelete), len(results.DeletedParameters))
	}

	return nil
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

func tagParameter(svc ssmiface.SSMAPI, resourceID string, resourceType string, tags []*ssm.Tag) (*ssm.AddTagsToResourceOutput, error) {
	results, err := svc.AddTagsToResource(&ssm.AddTagsToResourceInput{
		ResourceId:   &resourceID,
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
