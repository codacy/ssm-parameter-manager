package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"gopkg.in/yaml.v3"
)

// class xpto () {
// 	// inint svc aqui
// }
// Check how to remove passing ssm every time

const unencryptedParameterFile = "../ssm-test-values.sops.yaml"

func main() {
	environment := flag.String("e", "default", "Name of the AWS Environment")
	//flag.Bool("v", false, "Verbose") Warning this will print your secrets to stdout in PLAIN TEXT
	flag.Parse()

	session := newAWSSession(*environment)
	svc := ssm.New(session)

	data := parseConfigurationFile(unencryptedParameterFile)

	processParameters(svc, data, *environment)
}

func parseConfigurationFile(filePath string) map[string]map[string]string {
	yfile, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]map[string]string)

	if err := yaml.Unmarshal(yfile, &data); err != nil {
		log.Fatal(err)
	}

	return data
}

func processParameters(svc ssmiface.SSMAPI, parameters map[string]map[string]string, environment string) {
	for k, v := range parameters[environment] {
		// print info on what it is parsing, number, etc
		resultsPut, err := putParameter(svc, k, v, "String", true)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(*resultsPut)

		resultsTag, err := tagParameter(svc, k, "Parameter", createTags())

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(*resultsTag) // Write keys to terminal (hide value in case of secrets)
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

func newAWSSession(profile string) *session.Session {
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
