.PHONY: run deploy

run:
	go run main.go

deploy:
	rm -f main handler.zip
	go build -o main main.go
	zip handler.zip main
	aws lambda update-function-code --function-name ${FUNCTION_NAME} --zip-file fileb://${PWD}/handler.zip
