# igz_rates

Get the rates per user from harvest

Install aws cli
https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

GOOS=linux go build -o my-lambda-binary main.go

zip function.zip my-lambda-binary

aws lambda update-function-code --function-name arn:aws:lambda:eu-west-1:374208052150:function:igz_rates --zip-file fileb://function.zip 

aws lambda invoke     --function-name igz_rates      --cli-binary-format raw-in-base64-out     --payload '{"id": "tpaschalis", "val": 100, "flag": true}'     response.json