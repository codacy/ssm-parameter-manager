package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v3"
)

const unencryptedParameterFile = "/Users/loliveira/workspace/onboarding-project/ssm-test-values.sops.yaml"
const encryptedParameterFile = "/Users/loliveira/workspace/onboarding-project/ssm-test-values-encrypted.sops.yaml"

func main() {
	environment := flag.String("e", "default", "Name of the AWS Environment")
	//flag.Bool("v", false, "Verbose") Warning this will print your secrets to stdout in PLAIN TEXT
	flag.Parse()

	session := newAWSSession(*environment)
	svc := ssm.New(session)

	plainData := parseConfigurationFile(unencryptedParameterFile, false)
	encryptedData := parseConfigurationFile(encryptedParameterFile, true)

	processParameters(svc, plainData, *environment)
	processParameters(svc, encryptedData, *environment)
}

func parseConfigurationFile(filePath string, encrypted bool) map[string]map[string]string {
	var yfile []byte
	var err error

	if encrypted {
		yfile, err = decryptWithSops(filePath)
	} else {
		yfile, err = ioutil.ReadFile(filePath)
	}

	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]map[string]string)

	if err := yaml.Unmarshal(yfile, &data); err != nil {
		log.Fatal(err)
	}

	return data
}
