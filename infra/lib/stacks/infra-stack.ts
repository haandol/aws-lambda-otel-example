import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as lambdaGo from '@aws-cdk/aws-lambda-go-alpha';

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const fn = new lambdaGo.GoFunction(this, 'LambdaFunction', {
      entry: path.resolve(__dirname, '..', '..', 'functions', 'cmd', 'api'),
      bundling: {
        goBuildFlags: ['-ldflags "-s -w"'],
      },
      architecture: lambda.Architecture.ARM_64,
      runtime: lambda.Runtime.PROVIDED_AL2,
      timeout: cdk.Duration.seconds(10),
      tracing: lambda.Tracing.ACTIVE,
      layers: [
        // https://aws-otel.github.io/docs/getting-started/lambda/lambda-go
        lambda.LayerVersion.fromLayerVersionArn(
          this,
          `OtelLayer`,
          `arn:aws:lambda:${this.region}:901920570463:layer:aws-otel-collector-arm64-ver-0-90-1:1`
        ),
      ],
    });
    fn.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ['xray:*'],
        resources: ['*'],
      })
    );
  }
}
