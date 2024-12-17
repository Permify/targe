package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func (op *Api) FindUser(ctx context.Context, username string) (*iam.GetUserOutput, error) {
	return op.client.GetUser(ctx, &iam.GetUserInput{
		UserName: aws.String(username),
	})
}

func (op *Api) FindPolicy(ctx context.Context, arn string) (*iam.GetPolicyOutput, error) {
	return op.client.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	})
}

func (op *Api) FindGroup(ctx context.Context, groupname string) (*iam.GetGroupOutput, error) {
	return op.client.GetGroup(ctx, &iam.GetGroupInput{
		GroupName: aws.String(groupname),
	})
}

func (op *Api) FindRole(ctx context.Context, rolename string) (*iam.GetRoleOutput, error) {
	return op.client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(rolename),
	})
}
