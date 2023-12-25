# AWS Lambda Opentelemetry Example

This repository is an example of collecting telemetry from AWS Lambda and exporting to AWS X-Ray.

# Prerequisites

- awscli
- Python 3.11+
- Nodejs 16+
- AWS Account and Locally configured AWS credential

# Installation

Deploying CDK provisions below infrastructure on your AWS account

## Install Lambda Dependencies

Install dependencies

```bash
$ npm i -g cdk@2.116.1
```

## Provision Infrastructure

Install cdk in global context and run `cdk init` if you did not initailize cdk yet.

```bash
$ cd infra
$ npm i
$ cdk init
$ cdk bootstrap
```

copy [infra/config/dev.toml](./infra/config/dev.toml) to `.toml`.

```bash
$ cp config/dev.toml .toml
```

Deploy infrastructure using CDK on AWS

```bash
$ cdk deploy "*" --require-approval never
```

# Usage

Invoke lambda function few times via AWS Console or CLI

and visit AWS X-Ray Dashboard.
