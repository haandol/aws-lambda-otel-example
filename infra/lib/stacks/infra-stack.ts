import * as path from 'path';
import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as lambdaGo from '@aws-cdk/aws-lambda-go-alpha';

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const bucket = new s3.Bucket(this, 'Bucket');
    const fn = new lambdaGo.GoFunction(this, 'LambdaFunction', {
      entry: path.resolve(__dirname, '..', '..', 'functions', 'cmd', 'api'),
      bundling: {
        goBuildFlags: ['-ldflags "-s -w"'],
      },
      architecture: lambda.Architecture.ARM_64,
      timeout: cdk.Duration.seconds(8),
      environment: {
        AWS_LAMBDA_EXEC_WRAPPER: '/opt/otel-instrument',
        OPENTELEMETRY_COLLECTOR_CONFIG_FILE: `s3://${bucket.bucketName}.s3.${this.region}.amazonaws.com/collector.yml`,
      },
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
        actions: ['xray:*', 'cloudwatch:PutMetricData', 's3:GetObject'],
        resources: ['*'],
      })
    );
    bucket.grantRead(fn);
  }
}
