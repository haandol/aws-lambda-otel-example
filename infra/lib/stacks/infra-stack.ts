import * as path from 'path';
import * as cdk from '@aws-cdk/core';
import * as iam from '@aws-cdk/aws-iam';
import * as lambda from '@aws-cdk/aws-lambda';
import * as lambdaPython from '@aws-cdk/aws-lambda-python';
import { Namespace } from '../constants/config';

export class InfraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const fn = new lambdaPython.PythonFunction(this, 'LambdaFunction', {
      functionName: `${Namespace}OpenTelemetryTest`,
      entry: path.resolve(__dirname, '..', 'functions'),
      runtime: lambda.Runtime.PYTHON_3_8,
      index: 'index.py',
      handler: 'handler',
      tracing: lambda.Tracing.ACTIVE,
      layers: [
        // ver 1-7-1:1 is only for ap-northeast-2, use 1-5-0:3 for other regions
        lambda.LayerVersion.fromLayerVersionArn(this, `OtelLayer`, `arn:aws:lambda:${cdk.Stack.of(this).region}:901920570463:layer:aws-otel-python38-ver-1-7-1:1`)
      ],
      environment: {
        AWS_LAMBDA_EXEC_WRAPPER: '/opt/otel-instrument',
        OPENTELEMETRY_COLLECTOR_CONFIG_FILE: '/var/task/collector.yml',
      },
      timeout: cdk.Duration.seconds(10),
    })
    fn.addToRolePolicy(new iam.PolicyStatement({
      actions: ['xray:*'],
      resources: ['*'],
    }));
  }

}
