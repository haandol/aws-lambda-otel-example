import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as lambdaPython from '@aws-cdk/aws-lambda-python-alpha';

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const fn = new lambdaPython.PythonFunction(this, 'LambdaFunction', {
      entry: path.resolve(__dirname, '..', 'functions'),
      runtime: lambda.Runtime.PYTHON_3_12,
      architecture: lambda.Architecture.ARM_64,
      tracing: lambda.Tracing.ACTIVE,
      layers: [
        // https://aws-otel.github.io/docs/getting-started/lambda/lambda-python
        lambda.LayerVersion.fromLayerVersionArn(
          this,
          `OtelLayer`,
          `arn:aws:lambda:${
            cdk.Stack.of(this).region
          }:901920570463:layer:aws-otel-python-arm64-ver-1-20-0:3`
        ),
      ],
      environment: {
        AWS_LAMBDA_EXEC_WRAPPER: '/opt/otel-instrument',
        OPENTELEMETRY_COLLECTOR_CONFIG_FILE: '/var/task/collector.yml',
      },
    });
    fn.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['xray:*'],
        resources: ['*'],
      })
    );
  }
}
