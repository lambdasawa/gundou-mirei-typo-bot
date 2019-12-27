import cdk = require("@aws-cdk/core");
import iam = require("@aws-cdk/aws-iam");
import kms = require("@aws-cdk/aws-kms");
import lambda = require("@aws-cdk/aws-lambda");
import events = require("@aws-cdk/aws-events");
import eventsTargets = require("@aws-cdk/aws-events-targets");
import { Duration } from "@aws-cdk/core";

export class GundouMireiTypoBotInfraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const account = new iam.AccountPrincipal(this.account);

    const key = new kms.Key(this, "Key", {
      enableKeyRotation: true
    });
    const keyAlias = key.addAlias("alias/gundou-mirei-typo-bot");

    keyAlias.grantEncryptDecrypt(account);

    const fn = new lambda.Function(this, "MyFunction", {
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
      code: lambda.Code.fromAsset("lambda"),
      timeout: Duration.minutes(3)
    });

    keyAlias.grantEncryptDecrypt(fn);

    const rule = new events.Rule(this, "Rule", {
      schedule: events.Schedule.expression(`rate(5 minutes)`)
    });

    rule.addTarget(new eventsTargets.LambdaFunction(fn));

    new cdk.CfnOutput(this, "FunctionName", {
      value: fn.functionName
    });
    new cdk.CfnOutput(this, "KeyId", {
      value: keyAlias.keyId
    });
  }
}
