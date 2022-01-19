package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"codacy/ssm-parameter-manager/sops"
	"codacy/ssm-parameter-manager/ssm"

	"github.com/aws/aws-sdk-go/aws/session"
	awsSsm "github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/yaml.v3"
)

func main() {
	verbose := flag.Bool("v", false, "Prints information about the parameters being processed. This will print secrets to stdout in plain text!")
	plainFile := flag.String("plainFile", "", "Path to plain yaml files, separated by comma")
	encryptedFile := flag.String("encryptedFile", "", "Path to encrypted yaml files, separated by comma")
	parameterPrefix := flag.String("parameterPrefix", "", "Prefix for the parameters to be checked and deleted if they are not contained in the config files. If empty, will not delete any parameters.")
	flag.Parse()

	environment := os.Getenv("AWS_PROFILE")

	if environment == "" {
		log.Fatal("AWS_PROFILE is not set.")
	}

	svc := awsSsm.New(newAWSSession(environment))

	plainData := parseConfigurationFile(*plainFile, false)
	encryptedData := parseConfigurationFile(*encryptedFile, true)

	fmt.Printf("SSM Parameter Manager working environment: \"%s\"\n", environment)

	if *parameterPrefix == "" {
		fmt.Printf("Parameter prefix is not set - proceeding without deleting parameters.\n")
	} else {
		if !strings.HasSuffix(*parameterPrefix, "/") {
			*parameterPrefix += "/"
		}

		fmt.Printf("Checking parameters with prefix \"%s\"\n", *parameterPrefix)
		_, err := ssm.CleanParameters(svc, *parameterPrefix, true, plainData, encryptedData)

		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("Processing %d plain text parameters\n", len(plainData))
	_, err := ssm.ProcessParameters(svc, plainData, *verbose)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Processing %d encrypted parameters\n", len(encryptedData))
	_, err = ssm.ProcessParameters(svc, encryptedData, *verbose)

	if err != nil {
		log.Fatalln(err)
	}
}

func parseConfigurationFile(filePath string, encrypted bool) map[string]string {
	var yfile []byte
	var err error

	data := make(map[string]string)

	if filePath == "" {
		return data
	}

	if encrypted {
		yfile, err = sops.Decrypt(filePath)
	} else {
		yfile, err = ioutil.ReadFile(filePath)
	}

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Fatal(string(exiterr.Stderr))
		}

		log.Fatal(err)
	}

	if err := yaml.Unmarshal(yfile, &data); err != nil {
		log.Fatal(err)
	}

	return data
}

func newAWSSession(profile string) *session.Session {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))

	return session
}
