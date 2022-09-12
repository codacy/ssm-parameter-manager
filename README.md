# SSM Parameter Manager

A generic tool to achieve SSM parameter configuration "as code" and avoid manual input. The goal is to have a way to process **yaml** files containing project specific configuration, which will automatically be created in the respective AWS account/environment. 

Additionally, this tool can also be used to delete parameters which are not contained within the configuration files.

Each project using this tool needs to wire it in the respective circleci configuration and also needs to have the parameter files commited in its own repo.

It is capable of handling sensitive values, such as passwords, has it makes use of [Mozilla SOPS](https://github.com/mozilla/sops) combined with [AWS KMS](https://aws.amazon.com/kms/). Note that the native integration with SOPS and Go is not very well supported right now, so it makes use of a SOPS instalation in the OS for the decryption.

## Requirements

(**only necessary when working with encrypted values**)

*   Configured AWS CLI
*   Mozilla SOPS 
*   Auxiliary encryption script

Installing Mozilla SOPS: 
````sh
brew install sops
````
Auxiliary encryption script:
````sh
https://bitbucket.org/qamine/codacy-scripts/src/master/misc/sops-encryption.sh
````

## Usage

This tool will take up to two lists of yaml parameter files, one with plain text configurations and one with values encrypted by SOPS.

First, encrypted files (if any) will be decrypted using a KMS key called **sops-ssm-parameter-manager**, which exists in each AWS environment.

Each file will then be parsed, and individual entries will be created in the SSM Parameter Store of the configured AWS account.

Currently, this tool only works with the **AWS_PROFILE** variable, which will be used to point to the right AWS account.

If `-parameterPrefix` flag is specified, existing parameters with that prefix will be deleted from AWS.

Parameters created with this tool will be tagged with

```sh
ssm-managed=true
```

## Command

```sh
ssm-parameter-manager [flags]
```

## Flags

```
  -plainFile              Path to the plain yaml file
  -encryptedFile          Path to encrypted yaml file
  -parameterPrefix        Prefix for the parameters to be checked and deleted if they are not contained in the config files. If empty, will not delete any parameters.
  -v                      Prints parameters being processed, including secrets in plain text
```

### Warning

As of release 0.0.7, this tool has the capability to delete parameters of a given prefix if they are not defined in the configuration files, enabled by passing a flag when running and **disabled** by default. This means that if you delete a parameter from the configuration files it will be deleted from the AWS environment. This is the desired behaviour to encourage configuration as code, but this can be dangerous if some critical parameters are defined in some other way other than the config files, so caution is advised.

## Working with files with sensitive parameters

In order to safely commit configuration files that contain sensitive data, it is necessary to first encrypt them. There is an [auxiliary script](https://bitbucket.org/qamine/codacy-scripts/src/master/misc/sops-encryption.sh) to ease this process.

### Encrypting files

```sh
sops-encryption.sh "path/to/file/to/encrypt.yaml" 
```
This will encrypt the file for the current environment specified in the **AWS_PROFILE** variable

### Decrypting files

```sh
sops -d -i "path/to/file/to/decrypt.yaml" 
```

Make sure you're in **AWS_PROFILE** that corresponds to the file being decrytpted. Files encrypted by SOPS have enough metadata information to decrypt without further configuration.

### Adding information to encrypted files

Simply decrypt the file, add the information and encrypt the file again, using the steps described above.

## Examples

### Plain text configuration file example

See parameter key convention names defined on [the handbook](https://handbook.dev.codacy.org/engineering/guidelines/application-parameters.html#ssm-parameter-conventions).

>All values must be quoted, either with `'` or `"`.

```yaml
ssm/parameter/path/a: 'a'
ssm/parameter/path/b: 'b'
ssm/parameter/path/c: '3'
ssm/parameter/path/d: 'true'
ssm/parameter/path/e:
  type: StringList
  value: "a,string,list"
```

### Encrypted configuration file example

```yaml
ssm/parameter/path/super-secret-key: ENC[AES256_GCM,data:vQK6Gg+OzUK7QQ==,iv:w6bdRet/EVwvXwDwrDaxisO/IY1sP3fN/GkvPN+euzA=,tag:2qDAm80zvxDh8UbVlQWiXA==,type:str]
ssm/parameter/path/password: ENC[AES256_GCM,data:nZrFqSlh+sUU6fk=,iv:V/QGow5xbuoHeACDgmmz3P7x/ptsh8yfC/yB//hEvPU=,tag:4b7CVFGEVL3T5IAyZ2GiOw==,type:str]
sops:
    kms:
        - arn: arn:aws:kms:eu-west-1:364192610488:key/56dbc6b0-2c0c-4370-8ddc-c081224b5998
          created_at: "2021-11-17T15:33:30Z"
          enc: AQICAHhaT9VxfjV6Mz7/D51Xv9TEykcYSMnG46Hcc8rBp8MT6gFPA3VFkp8noTBK9TpRnfBMAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMA1QHMJykps7DVsSzAgEQgDsnK2KzEbh6C35fo221FI5WtnwIOLeLVhqyFwU5N5/73+ynWP3Fjvm/xkRH2Y+nYNzXK+mYxUHCljLluQ==
          aws_profile: ""
    gcp_kms: []
    azure_kv: []
    hc_vault: []
    age: []
    lastmodified: "2021-11-17T15:33:31Z"
    mac: ENC[AES256_GCM,data:uI4ygLB5Jj2awSnYYQGyNC9uG+QbdulmBYfO7YR+UFigXmBC2XXN9vo9tMQ5l32RYg447h9qsL+f7A8WZqfWbeIoe3T/lW4r7uEsvdsk+rX23ONczThrELYHF5YBE0wQcQDSNu5hxR2e30f755OU11ohcx159dFxyKUc1WyYUIM=,iv:2GsiXeZt3iLdUEAO2bdVZdRzqZVdnW2hdBcc62RT3Iw=,tag:oz0xFgSUZ5mC1a9vTpTpNw==,type:str]
    pgp: []
    unencrypted_suffix: _unencrypted
    version: 3.7.1

```

#### Example command
```sh
ssm-parameter-manager -plainFiles "/path/to/plain/file.yaml" -encryptedFiles "/path/to/encryptedfie.yaml"
```
