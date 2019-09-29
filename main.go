package main

import (
	"encoding/base64"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/lambdasawa/gundou-mirei-typo-bot/tw"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	switch {
	case isAWSLambdaEnv():
		lambda.Start(GetLambdaHandler())
	default:
		start()
	}
}

func isAWSLambdaEnv() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}

func GetLambdaHandler() func() error {
	kssService := kms.New(session.New(), aws.NewConfig().WithRegion("ap-northeast-1"))

	conf := tw.TwitterConfig{}
	var err error

	conf.ConsumerKey, err = decodeKMS(kssService, "CONSUMER_KEY")
	if err != nil {
		log.Fatal(err)
	}

	conf.ConsumerSecret, err = decodeKMS(kssService, "CONSUMER_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	conf.AccessToken, err = decodeKMS(kssService, "ACCESS_TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	conf.AccessSecret, err = decodeKMS(kssService, "ACCESS_SECRET")
	if err != nil {
		log.Fatal(err)
	}

	return func() error {
		return tw.Favorite(conf)
	}
}

func decodeKMS(svc *kms.KMS, envName string) (string, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(os.Getenv(envName))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode KMS data as Base64")
	}

	var in = &kms.DecryptInput{
		CiphertextBlob: dataBytes,
	}
	out, err := svc.Decrypt(in)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt KMS value")
	}

	return string(out.Plaintext), nil
}

func start() {
	if err := tw.Favorite(tw.TwitterConfig{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
		AccessToken:    os.Getenv("ACCESS_TOKEN"),
		AccessSecret:   os.Getenv("ACCESS_SECRET"),
	}); err != nil {
		log.Fatal(err)
	}
}
