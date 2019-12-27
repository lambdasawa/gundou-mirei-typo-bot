.PHONY: deploy destroy

deploy:
	test -n "${CONSUMER_KEY}" || exit 1
	test -n "${CONSUMER_SECRET}" || exit 1
	test -n "${ACCESS_TOKEN}" || exit 1
	test -n "${ACCESS_SECRET}" || exit 1

	GOOS=linux GOARCH=amd64 go build -trimpath -o main main.go
	mkdir -p infra/lambda/
	mv main infra/lambda/

	cd infra &&\
		npm run build &&\
		npm run cdk deploy &&\
		npm run post-deploy

destroy:
	cd infra && cdk destroy
