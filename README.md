# awsipacl
## Setup
_Make sure you have Go 1.16 or newer installed._

First, install the [AWS CLI version 2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) and [configure it](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html).

Then, from the Lambda console, create a new function. Select "Author from scratch", give it a name, and select "Go 1.x" as the runtime.

Open up the newly-created function, and, in the "Code" tab, edit the "Runtime settings". Set the "Handler" to "awsipacl". Also, in the "Configuration" tab, edit the memory to 128 MB.

Now, go to the API Gateway console, and build a new HTTP API. Add your Lambda function as an integration. Give the API the same name as your function. Set the resource path to "/{proxy+}". Leave the rest of the settings as their defaults.

At the top of the Lambda console page, select the "Add trigger" button and select "API Gateway" from the list. You'll want to create a new API, and select the "HTTP API" type. The security can be set to "Open".

Under the "Configuration" tab, select "Triggers", and then open the API Gateway that you just created. From the API Gateway console, select "Routes", and then the "ANY" route. Edit that route, setting its path to "/{proxy+}".

In a terminal, go to this folder and run `GOOS=linux go build .`, followed by `zip awsipacl.zip awsipacl`. Then run `aws lambda update-function-code --function-name FUNCTION-NAME-HERE --zip-file fileb://./awsipacl.zip`, but replace FUNCTION-NAME-HERE with the name of your Lambda function.