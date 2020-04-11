.PHONY: deploy undeploy

deploy:
	GOOS=linux GOARCH=amd64 go build -trimpath -o bin/main main.go
	sls deploy

undeploy:
	sls destroy
