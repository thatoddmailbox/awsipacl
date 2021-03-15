# awsipacl
A web app that lets you manage the IPs of an EC2 Security Group. It runs as an AWS Lambda function behind an Amazon API Gateway, with usage low enough to effectively be free. (less than a cent per month)

## Setup
## Configuration
Go to Security Groups in VPC and copy the security group ID.

## Deployment
_Make sure you have Go 1.16 or newer installed._

First, install the [AWS CLI version 2](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html) and [configure it](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html).

Then, from the Lambda console, create a new function. Select "Author from scratch", give it a name, and select "Go 1.x" as the runtime.

Open up the newly-created function, and, in the "Code" tab, edit the "Runtime settings". Set the "Handler" to "awsipacl". Also, in the "Configuration" tab, edit the memory to 128 MB.

In the "Configuration" tab, select "Permissions", and then click on the role under "Role name". This will open a new tab in the IAM console. From there, click "Attach policies", then "Create policy" (this will open a new tab). Select the JSON tab and paste the following in (MAKE SURE TO REPLACE `SECURITY-GROUP-ID-HERE` with your security group ID, including the `sg-` prefix):

```json
{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": [
				"ec2:AuthorizeSecurityGroupEgress",
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:RevokeSecurityGroupEgress",
				"ec2:RevokeSecurityGroupIngress"
			],
			"Effect": "Allow",
			"Resource": "arn:aws:ec2:*:*:security-group/SECURITY-GROUP-ID-HERE"
		},
        {
            "Action": [
                "ec2:DescribeSecurityGroups"
            ],
            "Effect": "Allow",
            "Resource": "*"
        }
	]
}
```

Click "Next: Tags" and don't enter any tags. Give the policy a name (ideally matching the name of your Lambda function) and click "Create policy". Close the tab you just used and go back to the "Attach Permissions" screen. From here, search for the policy you just made and attach it.

Now, go to the API Gateway console, and build a new HTTP API. Add your Lambda function as an integration. Give the API the same name as your function. Set the resource path to "/{proxy+}". Leave the rest of the settings as their defaults.

In a terminal, go to this folder and run `GOOS=linux go build .`, followed by `zip awsipacl.zip awsipacl`. Then run `aws lambda update-function-code --function-name FUNCTION-NAME-HERE --zip-file fileb://./awsipacl.zip`, but replace FUNCTION-NAME-HERE with the name of your Lambda function.