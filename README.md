# awsipacl
## Setup
First, install the [AWS CLI version 2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) and [configure it](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html).

Then, from the Lambda console, create a new function. Select "Author from scratch", give it a name, and select "Go 1.x" as the runtime.

Open up the newly-created function, and, in the "Code" tab, edit the "Runtime settings". Set the "Handler" to "backend". Also, in the "Configuration" tab, edit the memory to 128 MB.

In a terminal, go to the `backend` folder and run `GOOS=linux go build .`, followed by `zip backend.zip backend`. Then run `aws lambda update-function-code --function-name FUNCTION-NAME-HERE --zip-file fileb://./backend.zip`, but replace FUNCTION-NAME-HERE with the name of your Lambda function.