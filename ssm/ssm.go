package ssm

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

func ProcessParameters(svc ssmiface.SSMAPI, parameters map[string]map[string]string, environment string, verbose bool) {
	fmt.Printf("Processing %d parameters for %s environment\n", len(parameters[environment]), environment)
	for k, v := range parameters[environment] {
		if verbose {
			fmt.Printf("Putting and tagging parameter with key \"%s\" and value \"%s\"\n", k, v)
		}

		_, err := putParameter(svc, k, v, "String", true)

		if err != nil {
			log.Fatal(err)
		}

		_, err = tagParameter(svc, k, "Parameter", createTags())

		if err != nil {
			log.Fatal(err)
		}
	}
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

func NewAWSSession(profile string) *session.Session {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))

	return session
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
