package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lambdasawa/gundou-mirei-typo-bot/tw"
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
	conf := tw.TwitterConfig{}
	conf.ConsumerKey = os.Getenv("CONSUMER_KEY")
	conf.ConsumerSecret = os.Getenv("CONSUMER_SECRET")
	conf.AccessToken = os.Getenv("ACCESS_TOKEN")
	conf.AccessSecret = os.Getenv("ACCESS_SECRET")

	return func() error {
		return tw.Favorite(conf)
	}
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
