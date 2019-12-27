import "source-map-support/register";
import aws = require("aws-sdk");

const findCfnOutputValue = async (
  cloudFormation: aws.CloudFormation,
  key: string
): Promise<string> => {
  const stacksRes = await cloudFormation
    .describeStacks({
      StackName: "GundouMireiTypoBotInfraStack"
    })
    .promise();
  const stacks = stacksRes.Stacks || [];
  const outputs = stacks[0].Outputs || [];
  const output = outputs.find(o => o.OutputKey === key);
  return output?.OutputValue || "";
};

const encryptKMS = async (
  kms: aws.KMS,
  keyId: string,
  plaintext: string
): Promise<string> => {
  const encryptRes = await kms
    .encrypt({
      KeyId: keyId,
      Plaintext: plaintext
    })
    .promise();
  const cipherTextBlob = encryptRes.CiphertextBlob;
  if (cipherTextBlob instanceof Buffer) {
    return cipherTextBlob.toString("base64");
  }
  if (cipherTextBlob instanceof String) {
    return cipherTextBlob.toString();
  }
  throw new Error("unknown cipher text type");
};

const main = async (): Promise<void> => {
  aws.config.getCredentials((err: aws.AWSError) => {
    if (err) throw new Error(err.stack);
  });
  aws.config.region = "ap-northeast-1";
  aws.config.logger = console;

  const cloudFormation = new aws.CloudFormation();
  const keyId = await findCfnOutputValue(cloudFormation, "KeyId");
  const functionName = await findCfnOutputValue(cloudFormation, "FunctionName");

  const kms = new aws.KMS();
  const keys = [
    "CONSUMER_KEY",
    "CONSUMER_SECRET",
    "ACCESS_TOKEN",
    "ACCESS_SECRET"
  ];
  const envVars: { [_: string]: string } = {};
  for (let i = 0; i < keys.length; i++) {
    const key = keys[i];
    console.log(key);
    envVars[key] = await encryptKMS(kms, keyId, process.env[key] || "");
    console.log(envVars);
  }

  const lambda = new aws.Lambda();
  await lambda
    .updateFunctionConfiguration({
      FunctionName: functionName,
      Environment: {
        Variables: envVars
      }
    })
    .promise();
};

main();
