package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/yaml.v3"
)

func main() {
	environment := flag.String("e", "", "Target AWS Environment")
	verbose := flag.Bool("v", false, "Prints information about the parameters being processed. This will print secrets to stdout in plain text!")
	plainFiles := flag.String("plainFiles", "", "Path to plain yaml files, separated by comma")
	encryptedFiles := flag.String("encryptedFiles", "", "Path to encrypted yaml files, separated by comma")
	flag.Parse()

	if *environment == "" {
		log.Fatal("No AWS environment was specified")
	}

	session := newAWSSession(*environment)
	svc := ssm.New(session)

	for _, s := range strings.Split(*plainFiles, ",") {
		if s != "" {
			plainData := parseConfigurationFile(s, false)
			processParameters(svc, plainData, *environment, *verbose)
		}
	}

	for _, s := range strings.Split(*encryptedFiles, ",") {
		if s != "" {
			encryptedData := parseConfigurationFile(s, true)
			processParameters(svc, encryptedData, *environment, *verbose)
		}
	}
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
