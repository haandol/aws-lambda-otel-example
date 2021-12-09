# AWS Lambda Opentelemetry Example

This repository is an example of collecting telemetry from AWS Lambda and exporting to LightStep

# Prerequisites

- awscli
- Python 3.8
- Nodejs 12
- AWS Account and Locally configured AWS credential

# Installation

Deploying CDK provisions below infrastructure on your AWS account

## Install Lambda Dependencies

Install dependencies

```bash
$ cd infra
$ pip install -r lib/functions/requirements.txt -t lib/functions
```

## Setup LightStep

Login LightStep

visit settings

copy AccessToken to clipboard

paste it on [**collector.yml**](lib/functions/collector.yml)

## Provision Infrastructure

Install project dependencies

```bash
$ npm i
```

Install cdk in global context and run `cdk init` if you did not initailize cdk yet.

```bash
$ npm i -g cdk@1.134.0
$ cdk init
$ cdk bootstrap
```

Deploy infrastructure using CDK on AWS

```bash
$ cdk deploy "*" --require-approval never
```

# Usage

Invoke lambda function few times

and visit LightStep Service Discovery Dashboard.