package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func FindUser(ctx context.Context, cfg aws.Config, username string) (*iam.GetUserOutput, error) {
	client := iam.NewFromConfig(cfg)
	return client.GetUser(ctx, &iam.GetUserInput{
		UserName: aws.String(username),
	})
}

func FindPolicy(ctx context.Context, cfg aws.Config, arn string) (*iam.GetPolicyOutput, error) {
	client := iam.NewFromConfig(cfg)
	return client.GetPolicy(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	})
}

func FindGroup(ctx context.Context, cfg aws.Config, groupname string) (*iam.GetGroupOutput, error) {
	client := iam.NewFromConfig(cfg)
	return client.GetGroup(ctx, &iam.GetGroupInput{
		GroupName: aws.String(groupname),
	})
}

func FindRole(ctx context.Context, cfg aws.Config, rolename string) (*iam.GetRoleOutput, error) {
	client := iam.NewFromConfig(cfg)
	return client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(rolename),
	})
}
